package ts

import (
	"fmt"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/mickeyreiss/firemodel"
	"github.com/mickeyreiss/firemodel/version"
	"github.com/pkg/errors"
)

func init() {
	firemodel.RegisterModeler("ts", &Modeler{})
}

type Modeler struct{}

func (m *Modeler) Model(schema *firemodel.Schema, sourceCoder firemodel.SourceCoder) error {
	f, err := sourceCoder.NewFile("firemodel.d.ts")
	if err != nil {
		return errors.Wrapf(err, "firemodel/ts: create typescript file")
	}
	defer f.Close()

	if err := tpl.Execute(f, schema); err != nil {
		return errors.Wrapf(err, "firemodel/ts: generating typescript")
	}
	return nil
}

var (
	tpl = template.Must(template.
		New("file").
		Funcs(map[string]interface{}{
			"firemodelVersion": func() string { return version.Version },
			"toTypescriptType": toTypescriptType,
			"ToScreamingSnake": strcase.ToScreamingSnake,
			"ToLowerCamel":     strcase.ToLowerCamel,
			"ToCamel":          strcase.ToCamel,
			"getModelOption":   getModelOption,
			"getSchemaOption":  getSchemaOption,
			"interfaceName":    interfaceName,
		}).
		Parse(file),
	)
	_ = template.Must(tpl.New("model").Parse(model))
	_ = template.Must(tpl.New("enum").Parse(enum))
)

func interfaceName(sym string) string {
	return fmt.Sprintf("I%s", sym)
}

func toTypescriptType(firetype firemodel.SchemaFieldType, extras *firemodel.SchemaFieldExtras) string {
	switch firetype {
	case firemodel.Boolean:
		return "boolean"
	case firemodel.Integer, firemodel.Double:
		return "number"
	case firemodel.Timestamp:
		return "firestore.Timestamp"
	case firemodel.String:
		if extras != nil && extras.EnumType != "" {
			return extras.EnumType
		} else if extras != nil && extras.URL {
			return "URL"
		} else {
			return "string"
		}
	case firemodel.Bytes:
		return "firestore.Blob"
	case firemodel.Reference:
		if extras != nil && extras.ReferenceTo != "" {
			return fmt.Sprintf("DocumentReference<%s>", interfaceName(extras.ReferenceTo))
		} else {
			return "firestore.DocumentReference"
		}
	case firemodel.GeoPoint:
		return "firestore.GeoPoint"
	case firemodel.Array:
		if extras != nil && extras.ArrayOfModel != "" {
			return fmt.Sprintf("%s[]", interfaceName(extras.ArrayOfModel))
		} else if extras != nil && extras.ArrayOfEnum != "" {
			return fmt.Sprintf("%s[]", extras.ArrayOfEnum)
		} else if extras != nil && extras.ArrayOfPrimitive != "" {
			return fmt.Sprintf("%s[]", toTypescriptType(extras.ArrayOfPrimitive, nil))
		} else {
			return "any[]"
		}
	case firemodel.Map:
		if extras != nil && extras.File {
			return "IFile"
		} else if extras != nil && extras.MapToModel != "" {
			return interfaceName(extras.MapToModel)
		} else if extras != nil && extras.MapToEnum != "" {
			return extras.MapToEnum
		} else if extras != nil && extras.MapToPrimitive != "" {
			return fmt.Sprintf("{ [key: string]: %s; }", toTypescriptType(extras.MapToPrimitive, nil))
		} else {
			return `{ [key: string]: any; }`
		}
	default:
		err := errors.Errorf("firemodel/ts: unknown type %s", firetype)
		panic(err)
	}
}

func getSchemaOption(namespace string, key string, defaultValue string, options firemodel.SchemaOptions) string {
	ns, ok := options[namespace]
	if !ok {
		return defaultValue
	}
	opt, ok := ns[key]
	if !ok {
		return defaultValue
	}
	return opt
}

func getModelOption(namespace string, key string, required bool, options firemodel.SchemaModelOptions) string {
	ns, ok := options[namespace]
	if !ok {
		if required {
			err := errors.Errorf("option %s.%s is required but not set", namespace, key)
			panic(err)
		} else {
			return ""
		}
	}
	opt, ok := ns[key]
	if !ok {
		if required {
			err := errors.Errorf("option %s.%s is required but not set", namespace, key)
			panic(err)
		} else {
			return ""
		}
	}
	return opt
}

const (
	file = `// DO NOT EDIT - Code generated by firemodel {{firemodelVersion}}.

import { firestore } from 'firebase';

export interface DocumentSnapshot<DataType = firestore.DocumentData>
  extends firestore.DocumentSnapshot {
  data(options?: firestore.SnapshotOptions): DataType | undefined;
}
export interface QueryDocumentSnapshot<
  DataType = firestore.DocumentData
> extends firestore.QueryDocumentSnapshot {
  data(options?: firestore.SnapshotOptions): DataType | undefined;
}
export interface QuerySnapshot<DataType = firestore.DocumentData>
  extends firestore.QuerySnapshot {
  readonly docs: QueryDocumentSnapshot<DataType>[];
}
export interface DocumentSnapshotExpanded<
  DataType = firestore.DocumentData
> {
  exists: firestore.DocumentSnapshot['exists'];
  ref: firestore.DocumentSnapshot['ref'];
  id: firestore.DocumentSnapshot['id'];
  metadata: firestore.DocumentSnapshot['metadata'];
  data: DataType;
}
export interface QuerySnapshotExpanded<
  DataType = firestore.DocumentData
> {
  metadata: {
    hasPendingWrites: firestore.QuerySnapshot['metadata']['hasPendingWrites'];
    fromCache: firestore.QuerySnapshot['metadata']['fromCache'];
  };
  size: firestore.QuerySnapshot['size'];
  empty: firestore.QuerySnapshot['empty'];
  docs: {
    [docId: string]: DocumentSnapshotExpanded<DataType>;
  };
}
export interface DocumentReference<DataType> extends firestore.DocumentReference {
  data(options?: firestore.SnapshotOptions): DataType | undefined;
  get(options?: firestore.GetOptions): Promise<DocumentSnapshot<DataType>>;
}
export interface CollectionReference<
  DataType = firestore.DocumentData
> extends firestore.CollectionReference {
  get(options?: firestore.GetOptions): Promise<QuerySnapshot<DataType>>;
}
export interface Collection<DataType = firestore.DocumentData> {
  [id: string]: DocumentSnapshotExpanded<DataType>;
}

// tslint:disable-next-line:no-namespace
export namespace {{.Options | getSchemaOption "ts" "namespace" "firemodel"}} {
  type URL = string;

  export interface IFile {
    url: URL;
    mimeType: string;
    name: string;
	}

  {{- range .Enums -}}
  {{- template "enum" .}}
  {{- end}}
  {{- range .Models -}}
  {{- template "model" .}}
  {{- end}}
}
`
	model = `
  {{- if .Comment}}

  /** {{.Comment}} */
  {{- else}}

  /** TODO: Add documentation to {{.Name}}. */
  {{- end}}
  export interface {{.Name | interfaceName | ToCamel}} {
    {{- range .Collections}}
    {{- if .Comment}}
    /** {{.Comment}} */
    {{- else }}
    /** TODO: Add documentation to {{.Name}}. */
    {{- end}}
    {{.Name | ToLowerCamel}}: CollectionReference<{{.Type | interfaceName | ToCamel}}>;
    {{- end}}

    {{- range .Fields}}
    {{- if .Comment}}
    /** {{.Comment}} */
    {{- else }}
    /** TODO: Add documentation to {{.Name}}. */
    {{- end}}
    {{.Name | ToLowerCamel -}}?: {{toTypescriptType .Type .Extras}};
    {{- end}}
    {{- if .Options | getModelOption "firestore" "autotimestamp" false}}

    /** Record creation timestamp. */
    createdAt?: firestore.Timestamp;
    /** Record update timestamp. */
    updatedAt?: firestore.Timestamp;
    {{- end}}
  }`

	enum = `
  {{- if .Comment}}

  /** {{.Comment}} */
  {{- else}}

  /** TODO: Add documentation to {{.Name}}. */
  {{- end}}
  export enum {{.Name | ToCamel}} {
    {{- range .Values}}
    {{- if .Comment}}
    /** {{.Comment}} */
    {{- else}}
    /** TODO: Add documentation to {{.Name}}. */
    {{- end}}
    {{.Name}} = "{{.Name | ToScreamingSnake}}",
    {{- end}}
  }`
)
