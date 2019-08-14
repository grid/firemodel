package ios

import (
	"fmt"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/visor-tax/firemodel"
	"github.com/visor-tax/firemodel/version"
)

func init() {
	firemodel.RegisterModeler("ios", &Modeler{})
}

type Modeler struct{}

func (m *Modeler) Model(schema *firemodel.Schema, sourceCoder firemodel.SourceCoder) error {
	f, err := sourceCoder.NewFile("Firemodel.swift")
	if err != nil {
		return errors.Wrapf(err, "firemodel/ios: create swift file")
	}
	defer f.Close()

	if err := tpl.Execute(f, schema); err != nil {
		return errors.Wrapf(err, "firemodel/ios: generating swift")
	}
	return nil
}

var (
	tpl = template.Must(template.
		New("file").
		Funcs(map[string]interface{}{
			"firemodelVersion":             func() string { return version.Version },
			"toSwiftType":                  toSwiftType,
			"toScreamingSnake":             strcase.ToScreamingSnake,
			"toCamel":                      strcase.ToCamel,
			"toLowerCamel":                 strcase.ToLowerCamel,
			"toSwiftFieldName":             toSwiftFieldName,
			"filterFieldsEnumsOnly":        filterFieldsEnumsOnly,
			"asEnum":                       asEnum,
			"asArray":                      asArray,
			"hasAnyAssociatedValues":       hasAnyAssociatedValues,
			"filterFieldsNonEnumsOnly":     filterFieldsNonEnumsOnly,
			"filterFieldsStructsOnly":      filterFieldsStructsOnly,
			"filterFieldsStructArraysOnly": filterFieldsStructArraysOnly,
			"filterFieldsEnumArraysOnly":   filterFieldsEnumArraysOnly,
			"isDecodingType":               isDecodingType,
			//"firestoreModelName":           firestoreModelName,
		}).
		Parse(file),
	)
	_ = template.Must(tpl.New("base").Parse(base))
	_ = template.Must(tpl.New("refs").Parse(refs))
	_ = template.Must(tpl.New("client").Parse(client))
	_ = template.Must(tpl.New("decoder").Parse(decoder))
	_ = template.Must(tpl.New("coding").Parse(coding))
	_ = template.Must(tpl.New("enum").Parse(enum))
)

func isDecodingType(codingType string, in *firemodel.SchemaField) bool {
	switch codingType {
	case "decodeIfPresent":
		switch in.Type.(type) {
		case *firemodel.String,
			*firemodel.Boolean,
			*firemodel.Integer,
			*firemodel.Double,
			*firemodel.GeoPoint,
			*firemodel.Timestamp,
			*firemodel.Bytes,
			*firemodel.Reference,
			*firemodel.Map,
			*firemodel.URL,
			*firemodel.Struct,
			*firemodel.File:
			return true
		}
	case "decodeEnum":
		switch in.Type.(type) {
		case *firemodel.Enum:
			return true
		}
	case "decodeArray":
		switch in.Type.(type) {
		case *firemodel.Array:
			return true
		}
	}
	return false
}

func asEnum(in firemodel.SchemaFieldType) *firemodel.SchemaEnum {
	return in.(*firemodel.Enum).T
}

func asArray(in firemodel.SchemaFieldType) *firemodel.Array {
	return in.(*firemodel.Array)
}

func hasAnyAssociatedValues(field firemodel.SchemaFieldType) bool {
	enum := asEnum(field)
	for _, value := range enum.Values {
		if value.AssociatedValue != nil {
			return true
		}
	}
	return false
}

func filterFieldsEnumsOnly(in []*firemodel.SchemaField) []*firemodel.SchemaField {
	var out []*firemodel.SchemaField
	for _, i := range in {
		if _, ok := i.Type.(*firemodel.Enum); !ok {
			continue
		}
		out = append(out, i)
	}
	return out
}

func filterFieldsNonEnumsOnly(in []*firemodel.SchemaField) []*firemodel.SchemaField {
	var out []*firemodel.SchemaField
	for _, i := range in {
		if _, ok := i.Type.(*firemodel.Enum); ok {
			continue
		}
		out = append(out, i)
	}
	return out
}

func filterFieldsStructsOnly(in []*firemodel.SchemaField) []*firemodel.SchemaField {
	var out []*firemodel.SchemaField
	for _, i := range in {
		if _, ok := i.Type.(*firemodel.Struct); !ok {
			continue
		}
		out = append(out, i)
	}
	return out
}

func filterFieldsStructArraysOnly(in []*firemodel.SchemaField) []*firemodel.SchemaField {
	var out []*firemodel.SchemaField
	for _, i := range in {
		t, ok := i.Type.(*firemodel.Array)
		if !ok {
			continue
		}
		if _, ok := t.T.(*firemodel.Struct); !ok {
			continue
		}
		out = append(out, i)
	}
	return out
}

func filterFieldsEnumArraysOnly(in []*firemodel.SchemaField) []*firemodel.SchemaField {
	var out []*firemodel.SchemaField
	for _, i := range in {
		t, ok := i.Type.(*firemodel.Array)
		if !ok {
			continue
		}
		if _, ok := t.T.(*firemodel.Enum); !ok {
			continue
		}
		out = append(out, i)
	}
	return out
}

func toSwiftType(root bool, firetype firemodel.SchemaFieldType) string {
	switch firetype := firetype.(type) {
	case *firemodel.Boolean:
		if root {
			return "Bool?"
		} else {
			return "Bool"
		}
	case *firemodel.Integer:
		if root {
			return "Int?"
		} else {
			return "Int"
		}
	case *firemodel.Double:
		if root {
			return "Float?"
		} else {
			return "Float"
		}
	case *firemodel.Timestamp:
		if root {
			return "Date?"
		} else {
			return "Date"
		}
	case *firemodel.URL:
		if root {
			return "URL?"
		} else {
			return "URL"
		}
	case *firemodel.String:
		if root {
			return "String?"
		} else {
			return "String"
		}
	case *firemodel.Bytes:
		if root {
			return "Data?"
		} else {
			return "Data"
		}
	case *firemodel.Reference:
		// TODO: Generic on reference should be required by compiler.
		if firetype.T == nil {
			panic("Reference type is required.")
		}
		if root {
			return fmt.Sprintf("%s?", refType(firetype.T))
		} else {
			return refType(firetype.T)
		}
	case *firemodel.GeoPoint:
		if root {
			return "GeoPoint?"
		} else {
			return "GeoPoint"
		}
	case *firemodel.Array:
		if firetype.T == nil {
			panic("Untyped arrays are not supported")
		}

		if root {
			return fmt.Sprintf("[%s]", toSwiftType(false, firetype.T))
		} else {
			return fmt.Sprintf("[%s]?", toSwiftType(false, firetype.T))
		}
	case *firemodel.File:
		if root {
			return "Pring.File?"
		} else {
			return "Pring.File"
		}
	case *firemodel.Struct:
		if root {
			return fmt.Sprintf("%s?", firetype.T.Name)
		} else {
			return firetype.T.Name
		}
	case *firemodel.Enum:
		if root {
			return fmt.Sprintf("%s?", strcase.ToCamel(firetype.T.Name))
		} else {
			return strcase.ToCamel(firetype.T.Name)
		}
	case *firemodel.Map:
		if firetype.T == nil {
			panic("Untyped maps are not supported")
		}
		return fmt.Sprintf("[String: %s]", toSwiftType(false, firetype.T))
	default:
		err := errors.Errorf("firemodel/ios: unknown type %s", firetype)
		panic(err)
	}
}

func refType(schemaModel *firemodel.SchemaModel) string {
	return fmt.Sprintf("%sRef", strcase.ToCamel(schemaModel.Name))
}

func firestoreModelName(model firemodel.SchemaModel) string {
	modelName, err := model.Options.GetFirestoreModelName()
	if err != nil {
		panic(err)
	}

	return modelName
}

func toSwiftFieldName(in string) string {
	lower := strcase.ToLowerCamel(in)
	switch lower {
	case "where":
		return fmt.Sprintf("`%s`", lower)
	}
	return lower
}

const (
	file = `// DO NOT EDIT - Code generated by firemodel {{firemodelVersion}}.

import Foundation
import FirebaseFirestore

// MARK: - Models
{{range .Models -}}
{{- template "base" .}}
{{- end}}

// MARK: - Structs
{{range .Structs -}}
{{template "base" .}}
{{- end}}

// MARK: - Enums
{{range .Enums -}}
{{template "enum" .}}
{{- end}}

// MARK: - References
{{range .Models -}}
{{template "refs" .}}
{{- end}}

// MARK: - Coding 
{{range .Models -}}
{{template "coding" .}}
{{- end}}
{{range .Structs -}}
{{template "coding" .}}
{{- end}}
{{/* range .Enums -}}
{{template "coding" .}}
{{- end */}}

{{template "client" .}}

// MARK: - Decoder
{{template "decoder" .}}

// MARK: - Protocols 

// MARK: - Standard Types

struct GeoPoint {
  let latitude: Float
  let longitude: Float
}

`

	base = `
{{- if .Comment}}
// {{.Comment}}
{{- end}}
struct {{.Name | toCamel}} {
    {{- range .Fields}}
    {{- if .Comment}}
    // {{.Comment}}
    {{- end}}
    let {{.Name | toSwiftFieldName -}}: {{.Type | toSwiftType true}}
    {{- end}}
}
`

	coding = `
extension {{.Name | toCamel}}: Decodable {
    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
		{{- range $field := .Fields}}
		{{- if . | isDecodingType "decodeIfPresent" }}
        self.{{.Name | toSwiftFieldName}} = try container.decodeIfPresent({{.Type | toSwiftType false}}.self, forKey: .{{.Name | toSwiftFieldName}})
		{{- else if . | isDecodingType "decodeEnum" }}
        let {{.Name | toSwiftFieldName}}Type = try container.decodeIfPresent(String.self, forKey: .{{.Name | toSwiftFieldName}})
        let {{ .Name | toSwiftFieldName }} = try container.nestedContainer(keyedBy: {{.Type | toSwiftType false}}Type.self, forKey: .{{.Name | toSwiftFieldName}})
		{{- if .Type | hasAnyAssociatedValues }}
		switch {{.Name | toSwiftFieldName}}Type {
        {{- $enum := .Type | asEnum}}
		{{- range $enum.Values }}
		{{- if .AssociatedValue }}
		case {{ $field.Type | toSwiftType false }}.{{ .Name | toSwiftFieldName }}.rawValue:
		self.{{.Name}} = try container.decodeIfPresent({{ .AssociatedValue | toSwiftType false }}.self, forKey: .{{ .Name | toSwiftFieldName }})
  		{{- end }}
  		{{- end }}
		default:
			break
		}
		{{- end -}}
		{{- else if . | isDecodingType "decodeArray" }}
        self.{{ .Name | toSwiftFieldName }} =  try container.decodeIfPresent({{ .Type | asArray | toSwiftType true}})
		{{else -}}
        self.{{.Name | toSwiftFieldName}} = try container.decode me
		{{end -}}
		{{- end}}
    }

	// Coding keys for {{ .Name }}.
    enum CodingKeys: String, CodingKey {
		{{- range .Fields}}
		case {{ .Name | toSwiftFieldName }} = "{{ .Name }}"
		{{- end}}
    }

	{{- range .Fields }}
	{{- if . | isDecodingType "decodeEnum" }}
    // Coding keys for the {{ .Type | toSwiftType false }} enum’s associated value.
	enum {{ .Type | toSwiftType false -}}Type: String, CodingKey {
	{{- range ( .Type | asEnum ).Values }}
		case {{ .Name | toSwiftFieldName }} = "{{- .Name -}}"
	{{- end }}
	}
	{{- end }}
	{{- end }}
}
`

	decoder = `
struct DocumentSnapshotKey: CodingKey {
    let stringValue: String

    let intValue: Int? = nil

    init?(stringValue: String) {
        self.stringValue = stringValue
    }

    init?(intValue: Int) {
        return nil
    }
}

enum DocumentSnapshotDecodingError: Error {
    case firestoreClientMissing
}

struct DocumentSnapshotDecoder: Decoder {
    let documentSnapshot: DocumentSnapshot
    let codingPath: [CodingKey]
    let userInfo: [CodingUserInfoKey : Any]

    func container<Key>(keyedBy type: Key.Type) throws -> KeyedDecodingContainer<Key> where Key : CodingKey {
        return KeyedDecodingContainer(DocumentSnapshotKeyedDecodingContainerProtocol<Key>(documentSnapshot: documentSnapshot, codingPath: codingPath, userInfo: userInfo))
    }
    func singleValueContainer() throws -> SingleValueDecodingContainer {
        return DocumentSnapshotSingleValueDecodingContainer(documentSnapshot: documentSnapshot, codingPath: codingPath, userInfo: userInfo)
    }
    func unkeyedContainer() throws -> UnkeyedDecodingContainer {
        throw FiremodelError.typeError
    }
}

struct DocumentSnapshotUnkeyedDecodingContainer: UnkeyedDecodingContainer {
    let documentSnapshot: DocumentSnapshot

    let codingPath: [CodingKey]

    let userInfo: [CodingUserInfoKey: Any]

    var count: Int?

    var isAtEnd: Bool

    var currentIndex: Int

    mutating func decodeNil() throws -> Bool {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Bool.Type) throws -> Bool {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: String.Type) throws -> String {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Double.Type) throws -> Double {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Float.Type) throws -> Float {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Int.Type) throws -> Int {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Int8.Type) throws -> Int8 {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Int16.Type) throws -> Int16 {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Int32.Type) throws -> Int32 {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Int64.Type) throws -> Int64 {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: UInt.Type) throws -> UInt {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: UInt8.Type) throws -> UInt8 {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: UInt16.Type) throws -> UInt16 {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: UInt32.Type) throws -> UInt32 {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: UInt64.Type) throws -> UInt64 {
        throw FiremodelError.internalError
    }

    mutating func decode<T>(_ type: T.Type) throws -> T where T : Decodable {
        throw FiremodelError.internalError
    }

    mutating func nestedContainer<NestedKey>(keyedBy type: NestedKey.Type) throws -> KeyedDecodingContainer<NestedKey> where NestedKey : CodingKey {
        throw FiremodelError.internalError
    }

    mutating func nestedUnkeyedContainer() throws -> UnkeyedDecodingContainer {
        throw FiremodelError.internalError
    }

    mutating func superDecoder() throws -> Decoder {
        throw FiremodelError.internalError
    }
}

struct DocumentSnapshotSingleValueDecodingContainer: SingleValueDecodingContainer {
    let documentSnapshot: DocumentSnapshot
    let codingPath: [CodingKey]
    let userInfo: [CodingUserInfoKey: Any]

    func decodeNil() -> Bool {
        return false
    }

    func decode(_ type: Bool.Type) throws -> Bool {
        throw FiremodelError.internalError
    }

    func decode(_ type: String.Type) throws -> String {
        throw FiremodelError.internalError
    }

    func decode(_ type: Double.Type) throws -> Double {
        throw FiremodelError.internalError
    }

    func decode(_ type: Float.Type) throws -> Float {
        throw FiremodelError.internalError
    }

    func decode(_ type: Int.Type) throws -> Int {
        throw FiremodelError.internalError
    }

    func decode(_ type: Int8.Type) throws -> Int8 {
        throw FiremodelError.internalError
    }

    func decode(_ type: Int16.Type) throws -> Int16 {
        throw FiremodelError.internalError
    }

    func decode(_ type: Int32.Type) throws -> Int32 {
        throw FiremodelError.internalError
    }

    func decode(_ type: Int64.Type) throws -> Int64 {
        throw FiremodelError.internalError
    }

    func decode(_ type: UInt.Type) throws -> UInt {
        throw FiremodelError.internalError
    }

    func decode(_ type: UInt8.Type) throws -> UInt8 {
        throw FiremodelError.internalError
    }

    func decode(_ type: UInt16.Type) throws -> UInt16 {
        throw FiremodelError.internalError
    }

    func decode(_ type: UInt32.Type) throws -> UInt32 {
        throw FiremodelError.internalError
    }

    func decode(_ type: UInt64.Type) throws -> UInt64 {
        throw FiremodelError.internalError
    }

    func decode<T>(_ type: T.Type) throws -> T where T : Decodable {
        throw FiremodelError.internalError
    }
}

struct DocumentSnapshotKeyedDecodingContainerProtocol<Key>: KeyedDecodingContainerProtocol where Key : CodingKey {
    let documentSnapshot: DocumentSnapshot
    let codingPath: [CodingKey]
    let userInfo: [CodingUserInfoKey: Any]

    var allKeys: [Key] {
        return Array(self.documentSnapshot.data()?.keys.compactMap { Key(stringValue: $0) } ?? [])
    }

    func contains(_ key: Key) -> Bool {
        return self.documentSnapshot.data()?.keys.contains(key.stringValue) ?? false
    }

    func decodeNil(forKey key: Key) throws -> Bool {
        return documentSnapshot.get(key.stringValue) == nil
    }

    private func primativeValue(forKey key: Key) throws -> Any {
        let fp = FieldPath(codingPath.map { $0.stringValue } + [key.stringValue])
        guard case let .some(value) = documentSnapshot.get(fp) else {
            throw DecodingError.keyNotFound(key, DecodingError.Context(codingPath: codingPath, debugDescription: "key missing"))
        }
        return value
    }

    func decode(_ type: Bool.Type, forKey key: Key) throws -> Bool {
        guard let value = try primativeValue(forKey: key) as? Bool else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: String.Type, forKey key: Key) throws -> String {
        guard let value = try primativeValue(forKey: key) as? String else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: Double.Type, forKey key: Key) throws -> Double {
        guard let value = try primativeValue(forKey: key) as? Double else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: Float.Type, forKey key: Key) throws -> Float {
        guard let value = try primativeValue(forKey: key) as? Float else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: Int.Type, forKey key: Key) throws -> Int {
        guard let value = try primativeValue(forKey: key) as? Int else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: Int8.Type, forKey key: Key) throws -> Int8 {
        guard let value = try primativeValue(forKey: key) as? Int8 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: Int16.Type, forKey key: Key) throws -> Int16 {
        guard let value = try primativeValue(forKey: key) as? Int16 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: Int32.Type, forKey key: Key) throws -> Int32 {
        guard let value = try primativeValue(forKey: key) as? Int32 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: Int64.Type, forKey key: Key) throws -> Int64 {
        guard let value = try primativeValue(forKey: key) as? Int64 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: UInt.Type, forKey key: Key) throws -> UInt {
        guard let value = try primativeValue(forKey: key) as? UInt else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: UInt8.Type, forKey key: Key) throws -> UInt8 {
        guard let value = try primativeValue(forKey: key) as? UInt8 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: UInt16.Type, forKey key: Key) throws -> UInt16 {
        guard let value = try primativeValue(forKey: key) as? UInt16 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: UInt32.Type, forKey key: Key) throws -> UInt32 {
        guard let value = try primativeValue(forKey: key) as? UInt32 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: UInt64.Type, forKey key: Key) throws -> UInt64 {
        guard let value = try primativeValue(forKey: key) as? UInt64 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode<T>(_ type: T.Type, forKey key: Key) throws -> T where T : Decodable {
        let t = DocumentSnapshotDecoder(documentSnapshot: documentSnapshot, codingPath: codingPath + [key], userInfo: self.userInfo)
        return try type.init(from: t)
    }

    func nestedContainer<NestedKey>(keyedBy type: NestedKey.Type, forKey key: Key) throws -> KeyedDecodingContainer<NestedKey> where NestedKey : CodingKey {
        return KeyedDecodingContainer(DocumentSnapshotKeyedDecodingContainerProtocol<NestedKey>(documentSnapshot: documentSnapshot, codingPath: codingPath + [key], userInfo: userInfo))
    }

    func nestedUnkeyedContainer(forKey key: Key) throws -> UnkeyedDecodingContainer {
        guard let value = try primativeValue(forKey: key) as? [Any] else {
            throw DecodingError.typeMismatch([Any].self, DecodingError.Context(codingPath: codingPath, debugDescription: "Unexpected type for key \(key)"))
        }
        return DocumentSnapshotUnkeyedDecodingContainer(documentSnapshot: documentSnapshot, codingPath: codingPath + [key], userInfo: userInfo, count: value.count, isAtEnd: value.isEmpty, currentIndex: 0)
    }

    func superDecoder() throws -> Decoder {
        throw FiremodelError.internalError
    }

    func superDecoder(forKey key: Key) throws -> Decoder {
        throw FiremodelError.internalError
    }
}
`

	client = `
typealias Source = FirebaseFirestore.FirestoreSource

// MARK: - Protocols

protocol FiremodelDocumentSubscriber {
    associatedtype DocumentType
    func subscribe(receiver publish: @escaping (FiremodelDocumentEvent<DocumentType>) -> Void) -> FiremodelUnsubscriber
}

enum FiremodelDocumentEvent<T> {
    case snapshot(_: T, metadata: SnapshotMetadata)
    case error(Error)
}

protocol FiremodelCollectionSubscriber {
    associatedtype DocumentType
    func subscribe(withQuery applyQuery: ((Query) -> Query)?,
                   receiver publish: @escaping (FiremodelCollectionEvent<DocumentType>) -> Void) -> FiremodelUnsubscriber
}

class FiremodelUnsubscriber {
    private var listenerRegistration: ListenerRegistration?
    private var unsubscribeOnDeinit: Bool = true

    fileprivate init(listenerRegistration: ListenerRegistration) {
        self.listenerRegistration = listenerRegistration
    }

    // Shared prevents the automatic unsubscribe behavior for rare cases when the listener registration is never retained.
    func shared() {
        self.unsubscribeOnDeinit = false
    }

    deinit {
        if unsubscribeOnDeinit {
            unsubscribe()
        }
    }

    func unsubscribe() {
        listenerRegistration?.remove()
        listenerRegistration = nil
    }
}

enum FiremodelCollectionEvent<T> {
    case snapshot(_: [T], diff: (additions: [FiremodelChange<T>], modifications: [FiremodelChange<T>], removals: [FiremodelChange<T>]), metadata: SnapshotMetadata)
    case error(Error)
}

struct FiremodelChange<T> {
    let document: T
    let oldIndex: UInt
    let newIndex: UInt
}

// MARK: - Client

class FiremodelClient {
    private let firestore: FirebaseFirestore.Firestore

    static func dev() -> FiremodelClient {
        let firestore = FirebaseFirestore.Firestore.firestore()
        return FiremodelClient(firestore: firestore)
    }

    init(firestore: FirebaseFirestore.Firestore) {
        self.firestore = firestore
    }

    // MARK: - Root Collections

    func users() -> UserCollectionRef {
        return UserCollectionRef(ref: firestore.collection("users"), client: self)
    }

    func user(id: String) -> UserRef {
        return users().user(id: id)
    }

    // MARK: - Decoding

    func decode<T>(_ type: T.Type, from snapshot: FirebaseFirestore.DocumentSnapshot) throws -> T where T: Decodable {
        let decoder = DocumentSnapshotDecoder(documentSnapshot: snapshot,
                                              codingPath: [],
                                              userInfo: [firestoreClientDecodingKey: self])

        return try type.init(from: decoder)
    }

    func rawDocumentReference(_ path: String) -> DocumentReference {
        return self.firestore.document(path)
    }
}

fileprivate let firestoreClientDecodingKey = CodingUserInfoKey(rawValue: "firestore")!
`

	refs = `
struct {{ .Name | toCamel }}Ref: FiremodelDocumentSubscriber {
	fileprivate let ref: DocumentReference
	fileprivate let client: FiremodelClient

	// TODO: Subcollection refs
	// TODO: Parent refs
}

struct {{ .Name | toCamel }}CollectionRef: FiremodelCollectionSubscriber {
	fileprivate let ref: DocumentReference
	fileprivate let client: FiremodelClient

	// TODO: Subdoc refs
}

    func subscribe(withQuery applyQuery: ((Query) -> Query)? = nil,
                   receiver publish: @escaping (FiremodelCollectionEvent<{{.Name | toCamel}}>) -> Void) -> FiremodelUnsubscriber {

        let registration = (applyQuery?(ref) ?? ref)
            .addSnapshotListener { (snap: QuerySnapshot?, error: Error?) in
                if let error = error {
                    publish(.error(error))
                    return
                }
                guard let snap = snap else {
                    assertionFailure("Error was nil but Snapshot was also nil. This is unexpected behavior from addSnapshotListener!")
                    publish(.error(FiremodelError.internalError))
                    return
                }


                var documents = [{{.Name | toCamel}}]()
                var diff = (additions: [FiremodelChange<{{.Name | toCamel}}>](), modifications: [FiremodelChange<{{.Name | toCamel}}>](), removals: [FiremodelChange<{{.Name | toCamel}}>]())
                for change in snap.documentChanges {
                    let model: {{.Name | toCamel}}
                    do {
                        model = try self.client.decode({{.Name | toCamel}}.self, from: change.document)
                    } catch {
                        publish(.error(error))
                        return
                    }

                    documents.append(model)

                    let firemodelChange = FiremodelChange(document: model, oldIndex: change.oldIndex, newIndex: change.newIndex)

                    switch change.type {
                    case DocumentChangeType.added:
                        diff.additions.append(firemodelChange)
                    case DocumentChangeType.modified:
                        diff.modifications.append(firemodelChange)
                    case DocumentChangeType.removed:
                        diff.removals.append(firemodelChange)
                    default:
                        assertionFailure("unexpected firestore DocumentChangeType \(change.type)")
                    }

                }

                publish(.snapshot(documents, diff: diff, metadata: snap.metadata as SnapshotMetadata))
        }

        return FiremodelUnsubscriber(listenerRegistration: registration)
    }
}
`

	// CUSTOM COLLECTIONS

	// CUSTOM CODING
	//{{- if .Fields | requiresCustomEncodeDecode }}
	//
	//override func encode(_ key: String, value: Any?) -> Any? {
	//switch key {
	//{{- range .Fields | filterFieldsEnumsOnly}}
	//case "{{.Name | toLowerCamel}}":
	//return self.{{.Name | toLowerCamel}}?.firestoreValue
	//{{- end}}
	//{{- range .Fields | filterFieldsStructArraysOnly}}
	//case "{{.Name | toLowerCamel}}":
	//return self.{{.Name | toLowerCamel}}?.map { $0.rawValue }
	//{{- end}}
	//{{- range .Fields | filterFieldsStructsOnly}}
	//case "{{.Name | toLowerCamel}}":
	//return self.{{.Name | toLowerCamel}}?.rawValue
	//{{- end}}
	//{{- range .Fields | filterFieldsEnumArraysOnly}}
	//case "{{.Name | toLowerCamel}}":
	//return self.{{.Name | toLowerCamel}}?.map { $0.firestoreValue }
	//{{- end}}
	//default:
	//break
	//}
	//return nil
	//}
	//
	//override func decode(_ key: String, value: Any?) -> Bool {
	//switch key {
	//{{- range .Fields | filterFieldsEnumsOnly}}
	//case "{{.Name | toLowerCamel}}":
	//self.{{.Name | toLowerCamel}} = {{.Type | toSwiftType false }}(firestoreValue: value)
	//{{- end}}
	//{{- range .Fields | filterFieldsStructArraysOnly}}
	//case "{{.Name | toLowerCamel}}":
	//self.{{.Name | toLowerCamel}} = (value as? [[String: Any]])?
	//.enumerated()
	//.map { {{.Type.T | toSwiftType false }}(id: "{{.Name | toLowerCamel}}.\($0.offset)", value: $0.element) }
	//{{- end}}
	//{{- range .Fields | filterFieldsStructsOnly}}
	//case "{{.Name | toLowerCamel}}":
	//if let value = value as? [String: Any] {
	//self.{{.Name | toLowerCamel}} = {{.Type | toSwiftType false}}(id: "\(0)", value: value)
	//return true
	//}
	//{{- end}}
	//{{- range .Fields | filterFieldsEnumArraysOnly}}
	//case "{{.Name | toLowerCamel}}":
	//self.{{.Name | toLowerCamel}} = (value as? [String])?.compactMap { {{.Type.T | toSwiftType false }}(firestoreValue: $0) }
	//return true
	//{{- end}}
	//default:
	//break
	//}
	//return false
	//}
	//{{- end}}

	enum = `
{{- if .Comment}}
// {{.Comment}}
{{- end}}
enum {{.Name | toCamel }} {
    {{- range .Values}}
    {{- if .Comment}}
    // {{.Comment}}
    {{- end}}
    case {{.Name | toSwiftFieldName}}
    {{- if .AssociatedValue -}}
    ( {{- .AssociatedValue.T.Name -}} )
    {{- end -}}
    {{- end}}
}
`

	//	structTpl = `
	//{{- if .Comment}}
	//// {{.Comment}}
	//{{- end}}
	//@objcMembers class {{.Name | toCamel }}: Pring.Object {
	//    {{- range .Fields}}
	//    {{- if .Comment}}
	//    // {{.Comment}}
	//    {{- end}}
	//    var {{.Name | toLowerCamel -}}: {{.Type | toSwiftType true}}
	//    {{- end}}
	//    {{- if .Fields | requiresCustomEncodeDecode }}
	//
	//    override func encode(_ key: String, value: Any?) -> Any? {
	//        switch key {
	//        {{- range .Fields | filterFieldsEnumsOnly}}
	//        case "{{.Name | toLowerCamel}}":
	//            return self.{{.Name | toLowerCamel}}?.firestoreValue
	//        {{- end}}
	//        {{- range .Fields | filterFieldsStructArraysOnly}}
	//        case "{{.Name | toLowerCamel}}":
	//            return self.{{.Name | toLowerCamel}}?.map { $0.rawValue }
	//        {{- end}}
	//        {{- range .Fields | filterFieldsStructsOnly}}
	//        case "{{.Name | toLowerCamel}}":
	//            return self.{{.Name | toLowerCamel}}?.rawValue
	//        {{- end}}
	//        {{- range .Fields | filterFieldsEnumArraysOnly}}
	//        case "{{.Name | toLowerCamel}}":
	//            return self.{{.Name | toLowerCamel}}?.map { $0.firestoreValue }
	//        {{- end}}
	//        default:
	//            break
	//        }
	//        return nil
	//    }
	//
	//    override func decode(_ key: String, value: Any?) -> Bool {
	//        switch key {
	//        {{- range .Fields | filterFieldsEnumsOnly}}
	//        case "{{.Name | toLowerCamel}}":
	//            self.{{.Name | toLowerCamel}} = {{.Type | toSwiftType false }}(firestoreValue: value)
	//        {{- end}}
	//        {{- range .Fields | filterFieldsStructArraysOnly}}
	//        case "{{.Name | toLowerCamel}}":
	//            self.{{.Name | toLowerCamel}} = (value as? [[String: Any]])?
	//                .enumerated()
	//                .map { {{.Type.T | toSwiftType false }}(id: "\($0.offset)", value: $0.element) }
	//        {{- end}}
	//        {{- range .Fields | filterFieldsStructsOnly}}
	//        case "{{.Name | toLowerCamel}}":
	//          if let value = value as? [String: Any] {
	//            self.{{.Name | toLowerCamel}} = {{.Type | toSwiftType false}}(id: "\(0)", value: value)
	//            return true
	//          }
	//          {{- end}}
	//        {{- range .Fields | filterFieldsEnumArraysOnly}}
	//        case "{{.Name | toLowerCamel}}":
	//            self.{{.Name | toLowerCamel}} = (value as? [String])?.compactMap { {{.Type.T | toSwiftType false }}(firestoreValue: $0) }
	//			return true
	//        {{- end}}
	//        default:
	//            break
	//        }
	//        return false
	//    }
	//    {{- end}}
	//}
	//`
)
