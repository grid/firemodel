package firemodel

import (
	"testing"

	"os"
	"path"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-test/deep"
)

func TestParseSchema(t *testing.T) {
	tests := []struct {
		name    string
		want    *Schema
		wantErr bool
	}{
		{
			name: "empty",
			want: &Schema{
				Options: SchemaOptions{},
			},
		},
		{
			name: "empty_model",
			want: &Schema{
				Models: []*SchemaModel{
					{
						Name:    "Empty",
						Options: SchemaOptions{},
					},
				},
				Options: SchemaOptions{},
			},
		},
		{
			name: "simple",
			want: &Schema{
				Models: []*SchemaModel{
					{
						Name: "SimpleModel",
						Fields: []*SchemaField{
							{
								Name:   "foo",
								Type:   String,
								Extras: &SchemaFieldExtras{},
							},
						},
						Options: SchemaOptions{},
					},
				},
				Options: SchemaOptions{},
			},
		},
		{
			name: "full",
			want: &Schema{
				Models: []*SchemaModel{
					{
						Name:    "TestModel",
						Comment: "A Test is a test model.",
						Fields: []*SchemaField{
							{Name: "name",
								Comment: "The name.",
								Type:    String,
								Extras:  &SchemaFieldExtras{},
							},
							{
								Name:    "age",
								Comment: "The age.",
								Type:    Integer,
								Extras:  &SchemaFieldExtras{},
							},
							{
								Name:    "pi",
								Comment: "The number pi.",
								Type:    Double,
								Extras:  &SchemaFieldExtras{},
							},
							{
								Name:    "birthdate",
								Comment: "The birth date.",
								Type:    Timestamp,
								Extras:  &SchemaFieldExtras{},
							},
							{
								Name:    "is_good",
								Comment: "True if it is good.",
								Type:    Boolean,
								Extras:  &SchemaFieldExtras{},
							},
							{
								Name: "data", Type: Bytes,
								Extras: &SchemaFieldExtras{},
							},
							{
								Name:   "friend",
								Type:   Reference,
								Extras: &SchemaFieldExtras{},
							},
							{
								Name:   "location",
								Type:   GeoPoint,
								Extras: &SchemaFieldExtras{},
							},
							{
								Name:   "colors",
								Type:   Array,
								Extras: &SchemaFieldExtras{},
							},
							{
								Name:   "meta",
								Type:   Map,
								Extras: &SchemaFieldExtras{},
							},
							{
								Name:    "a_file",
								Comment: "Fake types...",
								Type:    Map,
								Extras:  &SchemaFieldExtras{File: true},
							},
							{
								Name:   "an_url",
								Type:   String,
								Extras: &SchemaFieldExtras{URL: true},
							},
						},
						Options: SchemaOptions{},
					},
				},
				Options: SchemaOptions{},
			},
		},
		{
			name: "extras",
			want: &Schema{
				Models: []*SchemaModel{
					{
						Name: "TestModel",
						Fields: []*SchemaField{
							{
								Name: "other", Type: Reference, Extras: &SchemaFieldExtras{ReferenceTo: "TestModel"}},
							{
								Name:   "unspecified_other",
								Type:   Reference,
								Extras: &SchemaFieldExtras{},
							},
							{
								Name:   "str_ary",
								Type:   Array,
								Extras: &SchemaFieldExtras{ArrayOfPrimitive: String},
							},
							{
								Name:   "model_ary",
								Type:   Array,
								Extras: &SchemaFieldExtras{ArrayOf: "TestModel"}},
							{
								Name:   "str_map",
								Type:   Map,
								Extras: &SchemaFieldExtras{MapToPrimitive: String},
							},
							{
								Name:   "model_map",
								Type:   Map,
								Extras: &SchemaFieldExtras{MapTo: "TestModel"},
							},
						},
						Options: SchemaOptions{},
					},
				},
				Options: SchemaOptions{},
			},
		},
		{
			name: "url",
			want: &Schema{
				Models: []*SchemaModel{
					{
						Name: "TestModel",
						Fields: []*SchemaField{
							{
								Name:   "url",
								Type:   String,
								Extras: &SchemaFieldExtras{URL: true},
							},
						},
						Options: SchemaOptions{},
					},
				},
				Options: SchemaOptions{},
			},
		},
		{
			name: "enums",
			want: &Schema{
				Enums: []*SchemaEnum{
					{
						Comment: "A cardinal direction.",
						Name:    "Direction",
						Values: []*SchemaEnumValue{
							{
								Comment: "Leftwards.",
								Name:    "left",
							},
							{
								Name: "right",
							},
							{
								Name: "up",
							},
							{
								Name: "down",
							},
						},
					},
				},
				Models: []*SchemaModel{
					{
						Name: "TestModel",
						Fields: []*SchemaField{
							{
								Comment: "The direction.",
								Name:    "dir",
								Type:    String,
								Extras:  &SchemaFieldExtras{EnumType: "Direction"},
							},
						},
						Options: SchemaOptions{},
					},
				},
				Options: SchemaOptions{},
			},
		},
		{
			name:    "error_nonsense",
			wantErr: true,
		},
		{
			name: "relational",
			want: &Schema{
				Models: []*SchemaModel{
					{
						Name: "Operator",
						Fields: []*SchemaField{
							{
								Name:   "operator_name",
								Type:   String,
								Extras: &SchemaFieldExtras{},
							},
						},
						Options: SchemaOptions{},
					},
					{
						Name: "Component",
						Fields: []*SchemaField{
							{
								Name:   "component_name",
								Type:   String,
								Extras: &SchemaFieldExtras{},
							},
						},
						Options: SchemaOptions{},
					},
					{
						Name: "Machine",
						Fields: []*SchemaField{
							{
								Name: "owner",
								Type: Reference,
								Extras: &SchemaFieldExtras{
									ReferenceTo: "Operator",
								},
							},
							// note: no components "field" here.
							{
								Name: "embedded_component",
								Type: Map,
								Extras: &SchemaFieldExtras{
									MapTo: "Component",
								},
							},
						},
						Options: SchemaOptions{},
						Collections: []*SchemaNestedCollection{
							{
								Name: "components",
								Type: "Component",
							},
						},
					},
				},
				Options: SchemaOptions{},
			},
		},
		{
			name: "language_options",
			want: &Schema{
				Options: SchemaOptions{
					"langname": {
						"optname": "optval",
					},
					"otherlang": {
						"str":       "hey",
						"number":    "1",
						"truth":     "true",
						"falseness": "false",
						"nullness":  "null",
					},
				},
			},
		},
		{
			name: "model_options",
			want: &Schema{
				Models: []*SchemaModel{
					{
						Name: "AnnotatedModel",
						Options: SchemaOptions{
							"lang": {
								"opt": "great",
							},
						},
					},
				},
				Options: SchemaOptions{},
			},
		},
		{
			name: "casing",
			want: &Schema{
				Models: []*SchemaModel{
					{
						Name: "NormalCase",
						Fields: []*SchemaField{
							{
								Name:   "foo_bar",
								Type:   String,
								Extras: &SchemaFieldExtras{},
							},
						},
						Options: SchemaOptions{},
					},
					{
						Name: "CamelCase",
						Fields: []*SchemaField{
							{
								Name:   "foo_bar",
								Type:   String,
								Extras: &SchemaFieldExtras{},
							},
						},
						Options: SchemaOptions{},
					},
					{
						Name: "TitleCase",
						Fields: []*SchemaField{
							{
								Name:   "foo_bar",
								Type:   String,
								Extras: &SchemaFieldExtras{},
							},
						},
						Options: SchemaOptions{},
					},
					{
						Name: "SnakeCase",
						Fields: []*SchemaField{
							{
								Name:   "foo_bar",
								Type:   String,
								Extras: &SchemaFieldExtras{},
							},
						},
						Options: SchemaOptions{},
					},
				},
				Options: SchemaOptions{},
			},
		},
		{
			name:    "reserved_model_name",
			wantErr: true,
		},
		{
			name:    "syntax_nonsense_2",
			wantErr: true,
		},
		{
			name: "model_named_user",
			want: &Schema{
				Models: []*SchemaModel{
					{
						Name:    "User",
						Comment: "Regression test.",
						Fields: []*SchemaField{
							{
								Name:   "name",
								Type:   String,
								Extras: &SchemaFieldExtras{},
							},
						},
						Options: SchemaOptions{},
					},
				},
				Options: SchemaOptions{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := os.Open(path.Join("testfixtures", tt.name+".firemodel"))
			if err != nil {
				t.Fatal(err)
			}
			got, err := ParseSchema(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Logf("ParseSchema() => %#v", spew.NewFormatter(got))
				for _, diff := range diff {
					t.Error(diff)
				}
			}
		})
	}
}
