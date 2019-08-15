package ios

import (
	"fmt"
	"github.com/jinzhu/inflection"
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

func (m *Modeler) Model(schema *firemodel.Schema, options firemodel.GenOptions, sourceCoder firemodel.SourceCoder) error {
	f, err := sourceCoder.NewFile("Firemodel.swift")
	if err != nil {
		return errors.Wrapf(err, "firemodel/ios: create swift file")
	}
	defer f.Close()

	if err := newTpl(schema).Execute(f, schema); err != nil {
		return errors.Wrapf(err, "firemodel/ios: generating swift")
	}
	return nil
}

func newTpl(schema *firemodel.Schema) *template.Template {
	tpl := template.Must(template.
		New("file").
		Funcs(map[string]interface{}{
			"firemodelVersion":            func() string { return version.Version },
			"toSwiftType":                 toSwiftType,
			"toScreamingSnake":            strcase.ToScreamingSnake,
			"toCamel":                     strcase.ToCamel,
			"toLowerCamel":                strcase.ToLowerCamel,
			"toSwiftFieldName":            toSwiftFieldName,
			"asEnum":                      asEnum,
			"asArray":                     asArray,
			"asMap":                       asMap,
			"hasAnyAssociatedValues":      hasAnyAssociatedValues,
			"isDecodingType":              isDecodingType,
			"pluralize":                   inflection.Plural,
			"DirectSubcollectionsOfModel": schema.DirectSubcollectionsOfModel,
			"RootModels":                  schema.RootModels,
			"ParentModel":                 schema.ParentModel,
			"lastTemplatePart":            lastTemplatePart,
		}).
		Parse(file),
	)
	_ = template.Must(tpl.New("base").Parse(base))
	_ = template.Must(tpl.New("refs").Parse(refs))
	_ = template.Must(tpl.New("client").Parse(client))
	_ = template.Must(tpl.New("decoder").Parse(decoder))
	_ = template.Must(tpl.New("coding").Parse(coding))
	_ = template.Must(tpl.New("refcoding").Parse(refcoding))
	_ = template.Must(tpl.New("enum").Parse(enum))
	_ = template.Must(tpl.New("protocol").Parse(protocol))
	return tpl
}

func lastTemplatePart(xs []firemodel.SchemaModelPathTemplatePart) firemodel.SchemaModelPathTemplatePart {
	return xs[len(xs)-1]
}

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
			*firemodel.URL,
			*firemodel.Struct:
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
	case "decodeMap":
		switch in.Type.(type) {
		case *firemodel.Map:
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

func asMap(in firemodel.SchemaFieldType) *firemodel.Map {
	return in.(*firemodel.Map)
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

// MARK: - Interfaces
{{range .Interfaces -}}
{{template "protocol" .}}
{{- end}}

// MARK: - References
{{range .Models -}}
{{template "refs" .}}
{{- end}}

// MARK: - Coding 
{{range .Models -}}
{{template "coding" .}}
{{- end}}
{{range .Models -}}
{{template "refcoding" .}}
{{- end}}
{{range .Structs -}}
{{template "coding" .}}
{{- end}}

{{template "client" .}}

// MARK: - Decoder
{{template "decoder" .}}

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

{{- $type := . -}}
{{- range .Implements }}
extension {{ $type.Name | toCamel }}: {{.Name | toCamel }} {
}
{{- end }}
`

	coding = `
extension {{.Name | toCamel}}: Decodable {
    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
		{{- range $field := .Fields}}
		{{- if . | isDecodingType "decodeIfPresent" }}
        self.{{.Name | toSwiftFieldName}} = try container.decodeIfPresent({{.Type | toSwiftType false}}.self, forKey: .{{.Name | toSwiftFieldName}})
		{{- else if . | isDecodingType "decodeEnum" }}
		{{- if .Type | hasAnyAssociatedValues }}
        let {{ .Name | toSwiftFieldName }}Container = try container.nestedContainer(keyedBy: {{.Type | toSwiftType false}}Type.self, forKey: .{{.Name | toSwiftFieldName}})
		{{- end }}
        let {{ .Name | toSwiftFieldName }}Value = try container.decodeIfPresent(String.self, forKey: .{{.Name | toSwiftFieldName}})
        switch {{ .Name | toSwiftFieldName }}Value {
        {{- $enum := .Type | asEnum}}
		{{- range $enum.Values }}
		case "{{ .Name | toScreamingSnake }}":
		{{- if .AssociatedValue }}
			self.{{ $field.Name | toSwiftFieldName }} = .{{.Name | toSwiftFieldName}}(try {{ $field.Name | toSwiftFieldName }}Container.decode({{ .AssociatedValue | toSwiftType false }}.self, forKey: {{ $field.Type | toSwiftType false }}Type.{{ .Name | toSwiftFieldName }}))
		{{- else }}
			self.{{ $field.Name | toSwiftFieldName }} = .{{.Name | toSwiftFieldName}}
  		{{- end }}
  		{{- end }}
		default:
			self.{{ $field.Name | toSwiftFieldName }} = .invalid({{ .Name | toSwiftFieldName }}Value)
		}
		{{- else if . | isDecodingType "decodeArray" }}
        self.{{ .Name | toSwiftFieldName }} =  try container.decode({{ .Type | asArray | toSwiftType true}}.self, forKey: .{{ .Name | toSwiftFieldName }})
		{{- else if . | isDecodingType "decodeMap" }}
        self.{{ .Name | toSwiftFieldName }} =  try container.decode({{ .Type | asMap | toSwiftType true}}.self, forKey: .{{ .Name | toSwiftFieldName }})
		{{else -}}
        self.{{.Name | toSwiftFieldName}} = /* TODO: {{ .Type }} decoding */
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

	refcoding = `
extension {{.Name | toCamel}}Ref: Decodable {
    init(from decoder: Decoder) throws {
        guard let client = decoder.userInfo[firestoreClientDecodingKey] as? FiremodelClient else {
            assertionFailure("firemodel client is missing in user info")
            throw DocumentSnapshotDecodingError.firestoreClientMissing
        }
        let container = try decoder.singleValueContainer()
        self.client = client
        self.ref  = client.rawDocumentReference(try container.decode(String.self))
    }
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
// MARK: - Client

class FiremodelClient {
    private let firestore: FirebaseFirestore.Firestore

	// Deprecated. Please use DI initializer instead.
    static func dev() -> FiremodelClient {
        let firestore = FirebaseFirestore.Firestore.firestore()
        return FiremodelClient(firestore: firestore)
    }

    init(firestore: FirebaseFirestore.Firestore) {
        self.firestore = firestore
    }

    // MARK: - Root Collections

	{{- range RootModels }}

    func {{ .Name | toSwiftFieldName | pluralize }}() -> {{ .Name }}CollectionRef {
        return {{ .Name }}CollectionRef(ref: firestore.collection("{{ (lastTemplatePart .FirestorePath.CollectionParts).CollectionName  }}"), client: self)
    }

    func {{ .Name | toSwiftFieldName }}(id: String) -> {{ .Name }}Ref {
        return {{ .Name }}Ref(ref: firestore.collection("{{ (lastTemplatePart .FirestorePath.CollectionParts).CollectionName  }}").document(id), client: self)
    }
	{{- end }}
}

// MARK: - Subscription Helpers

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

enum FiremodelCollectionEvent<T> {
    case snapshot(_: [T], diff: (additions: [FiremodelChange<T>], modifications: [FiremodelChange<T>], removals: [FiremodelChange<T>]), metadata: SnapshotMetadata)
    case error(Error)
}

struct FiremodelChange<T> {
    let document: T
    let oldIndex: UInt
    let newIndex: UInt
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


// MARK: - Decoding Helpers

extension FiremodelClient {


    fileprivate func decode<T>(_ type: T.Type, from snapshot: FirebaseFirestore.DocumentSnapshot) throws -> T where T: Decodable {
        let decoder = DocumentSnapshotDecoder(documentSnapshot: snapshot,
                                              codingPath: [],
                                              userInfo: [firestoreClientDecodingKey: self])

        return try type.init(from: decoder)
    }

    fileprivate func rawDocumentReference(_ path: String) -> DocumentReference {
        return self.firestore.document(path)
    }
}

enum FiremodelError: Error {
    case typeError
    case internalError
}

fileprivate let firestoreClientDecodingKey = CodingUserInfoKey(rawValue: "firestore")!
`

	refs = `
struct {{ .Name | toCamel }}CollectionRef {
	fileprivate let ref: CollectionReference
	fileprivate let client: FiremodelClient

	{{- with ParentModel . }}

	// MARK: - Parent Ref

    func parent() -> {{ .Name }}Ref {
        return {{ .Name }}Ref(ref: ref.parent!, client: client)
    }
	{{- end }}


	// MARK: - Child Ref

    func {{ .Name | toSwiftFieldName }}(id: String) -> {{ .Name }}Ref {
        return {{ .Name }}Ref(ref: ref.document(id), client: client)
    }
}

extension {{ .Name | toCamel }}CollectionRef: FiremodelCollectionSubscriber {
    func subscribe(withQuery applyQuery: ((Query) -> Query)? = nil,
                   receiver publish: @escaping (FiremodelCollectionEvent<{{.Name }}>) -> Void) -> FiremodelUnsubscriber {

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

struct {{ .Name | toCamel }}Ref {
	fileprivate let ref: DocumentReference
	fileprivate let client: FiremodelClient

	// MARK: - Parent Ref

    func parent() -> {{ .Name }}CollectionRef {
        return {{ .Name }}CollectionRef(ref: ref.parent, client: client)
    }

	{{- with $directSubcollections := . | DirectSubcollectionsOfModel }}

	// MARK: -  Subcollection Refs

	{{- range $directSubcollections }}

    func {{ .Name | toSwiftFieldName | pluralize }}() -> {{ .Name }}CollectionRef {
        return {{ .Name }}CollectionRef(ref: ref.collection("{{ (lastTemplatePart .FirestorePath.CollectionParts ).CollectionName  }}"), client: client)
    }

	{{- end }}
	{{- end }}
}

extension {{ .Name | toCamel }}Ref: FiremodelDocumentSubscriber {

    func subscribe(receiver publish: @escaping (FiremodelDocumentEvent<{{.Name}}>) -> Void) -> FiremodelUnsubscriber {
        let registration = ref
            .addSnapshotListener { (snap: DocumentSnapshot?, error: Error?) in
                if let error = error {
                    publish(.error(error))
                    return
                }
                guard let snap = snap else {
                    assertionFailure("Error was nil but Snapshot was also nil. This is unexpected behavior from addSnapshotListener!")
                    publish(.error(FiremodelError.internalError))
                    return
                }
                
                do {
                    let model = try self.client.decode({{.Name}}.self, from: snap)
                    publish(.snapshot(model, metadata: snap.metadata))
                } catch {
                    publish(.error(error))
                    return
                }
        }
        
        return FiremodelUnsubscriber(listenerRegistration: registration)
    }
}

`

	enum = `
{{- if .Comment}}
// {{.Comment}}
{{- end}}
enum {{.Name | toCamel }} {
	// An unknown enum value with its raw string.
	case invalid(String?)
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

	protocol = `
{{- if .Comment}}
// {{.Comment}}
{{- end}}
protocol {{.Name | toCamel }} {
    {{- range .Fields}}
    {{- if .Comment}}
    // {{.Comment}}
    {{- end}}
    var {{.Name | toSwiftFieldName}}: {{ .Type | toSwiftType true }} { get }
    {{- end}}
}
`
)
