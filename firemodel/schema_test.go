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
						Options: SchemaModelOptions{},
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
								Name: "foo",
								Type: &String{},
							},
						},
						Options: SchemaModelOptions{},
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
							{
								Name:    "name",
								Comment: "The name.",
								Type:    &String{},
							},
							{
								Name:    "age",
								Comment: "The age.",
								Type:    &Integer{},
							},
							{
								Name:    "pi",
								Comment: "The number pi.",
								Type:    &Double{},
							},
							{
								Name:    "birthdate",
								Comment: "The birth date.",
								Type:    &Timestamp{},
							},
							{
								Name:    "is_good",
								Comment: "True if it is good.",
								Type:    &Boolean{},
							},
							{
								Name: "data",
								Type: &Bytes{},
							},
							{
								Name: "friend",
								Type: &Reference{T: &SchemaModel{Name: "Friend"}},
							},
							{
								Name: "location",
								Type: &GeoPoint{},
							},
							{
								Name: "colors",
								Type: &Array{T: &String{}},
							},
							{
								Name: "meta",
								Type: &Map{},
							},
							{
								Name:    "a_file",
								Comment: "Fake types...",
								Type:    &File{},
							},
							{
								Name: "an_url",
								Type: &URL{},
							},
						},
						Options: SchemaModelOptions{},
					},
				},
				Options: SchemaOptions{},
			},
		},
		{
			name: "extras",
			want: &Schema{
				Enums: []*SchemaEnum{
					{Name: "TestEnum"},
				},
				Models: []*SchemaModel{
					{
						Name: "TestModel",
						Fields: []*SchemaField{
							{
								Name: "other",
								Type: &Reference{T: &SchemaModel{Name: "TestModel"}},
							},
							{
								Name: "unspecified_other",
								Type: &Reference{},
							},
							{
								Name: "primative_ary",
								Type: &Array{T: &String{}},
							},
							{
								Name: "model_ary",
								Type: &Array{T: &Model{ /*TestModel*/ }},
							},
							{
								Name: "enum_ary",
								Type: &Array{T: &Enum{ /*TestModel*/ }},
							},
							{
								Name: "reference_ary",
								Type: &Array{T: &Reference{ /* TestModel*/ }},
							},
							{
								Name: "nested_ary",
								Type: &Array{T: &Array{&String{}}},
							},
							{
								Name: "generic_ary",
								Type: &Array{},
							},
							{
								Name: "primative_map",
								Type: &Map{T: &String{}},
							},
							{
								Name: "model_map",
								Type: &Map{T: &Model{ /*TestModel*/ }},
							},
							{
								Name: "enum_map",
								Type: &Map{T: &Enum{ /*TestModel*/ }},
							},
							{
								Name: "generic_map",
								Type: &Map{},
							},
						},
						Options: SchemaModelOptions{},
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
								Name: "url",
								Type: &URL{},
							},
						},
						Options: SchemaModelOptions{},
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
								Type:    &Enum{ /*Direction*/ },
							},
						},
						Options: SchemaModelOptions{},
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
								Name: "operator_name",
								Type: &String{},
							},
						},
						Options: SchemaModelOptions{},
					},
					{
						Name: "Component",
						Fields: []*SchemaField{
							{
								Name: "component_name",
								Type: &String{},
							},
						},
						Options: SchemaModelOptions{},
					},
					{
						Name: "Machine",
						Fields: []*SchemaField{
							{
								Name: "owner",
								Type: &Reference{ /*Operator*/ },
							},
							// note: no components "field" here.
							{
								Name: "embedded_component",
								Type: &Model{T: &SchemaModel{Name: "Component"}},
							},
						},
						Options: SchemaModelOptions{},
						Collections: []*SchemaNestedCollection{
							{
								Name: "components",
								Type: &Model{ /* "Component"*/ },
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
						Options: SchemaModelOptions{
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
								Name: "foo_bar",
								Type: &String{},
							},
						},
						Options: SchemaModelOptions{},
					},
					{
						Name: "CamelCase",
						Fields: []*SchemaField{
							{
								Name: "foo_bar",
								Type: &String{},
							},
						},
						Options: SchemaModelOptions{},
					},
					{
						Name: "TitleCase",
						Fields: []*SchemaField{
							{
								Name: "foo_bar",
								Type: &String{},
							},
						},
						Options: SchemaModelOptions{},
					},
					{
						Name: "SnakeCase",
						Fields: []*SchemaField{
							{
								Name: "foo_bar",
								Type: &String{},
							},
						},
						Options: SchemaModelOptions{},
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
								Name: "name",
								Type: &String{},
							},
						},
						Options: SchemaModelOptions{},
					},
				},
				Options: SchemaOptions{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := os.Open(path.Join("testfixtures", "schema", tt.name+".firemodel"))
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
