package firemodel

import (
	"runtime/debug"
	"testing"

	"gotest.tools/assert"

	"os"
	"path"
)

func TestParseSchema(t *testing.T) {
	tests := []struct {
		name    string
		want    *Schema
		wantErr bool
	}{
		{
			name: "empty",
			want: &Schema{},
		},
		{
			name: "empty_model",
			want: &Schema{
				Models: []*SchemaModel{
					{
						Name: "Empty",
					},
				},
			},
		},
		{
			name: "simple",
			want: &Schema{
				Models: []*SchemaModel{
					{
						Name: "SimpleModel",
						FirestorePath: SchemaModelPathTemplate{
							Pattern: "/models/{model_id}",
						},
						Fields: []*SchemaField{
							{
								Name: "foo",
								Type: &String{},
							},
						},
					},
				},
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
								Type: &Reference{},
							},
							{
								Name: "location",
								Type: &GeoPoint{},
							},
							{
								Name: "colors",
								Type: &Array{},
							},
							{
								Name: "meta",
								Type: &Map{},
							},
							{
								Name: "an_url",
								Type: &URL{},
							},
						},
					},
				},
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
								Name: "struct_ary",
								Type: &Array{T: &Struct{T: &SchemaStruct{Name: "TestStruct"}}},
							},
							{
								Name: "enum_ary",
								Type: &Array{T: &Enum{T: &SchemaEnum{Name: "TestEnum"}}},
							},
							{
								Name: "reference_ary",
								Type: &Array{T: &Reference{T: &SchemaModel{Name: "TestModel"}}},
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
								Name: "struct_map",
								Type: &Map{T: &Struct{T: &SchemaStruct{Name: "TestStruct"}}},
							},
							{
								Name: "enum_map",
								Type: &Map{T: &Enum{T: &SchemaEnum{Name: "TestEnum"}}},
							},
							{
								Name: "generic_map",
								Type: &Map{},
							},
						},
					},
				},
				Structs: []*SchemaStruct{
					{
						Name: "TestStruct",
					},
				},
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
					},
				},
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
								Type: &Enum{
									T: &SchemaEnum{
										Name:    "Direction",
										Comment: "A cardinal direction.",
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
							},
						},
					},
				},
			},
		},
		{
			name: "enums_with_associated_values",
			want: &Schema{
				Enums: []*SchemaEnum{
					{
						Name: "Frobnicator",
						Values: []*SchemaEnumValue{
							{
								Name: "up",
							},
							{
								Name: "down",
							},
						},
					},
					{
						Name: "Computer",
						Values: []*SchemaEnumValue{
							{
								Name: "off",
							},
							{
								Name: "on",
								AssociatedValue: &Struct{
									T: &SchemaStruct{
										Name: "ComputerOnState",
									},
								},
							},
						},
					},
				},
				Structs: []*SchemaStruct{
					{
						Name: "ComputerOnState",
						Fields: []*SchemaField{
							{
								Name: "processes",
								Type: &Integer{},
							},
							{
								Name: "frob",
								Type: &Enum{
									T: &SchemaEnum{
										Name: "Frobnicator",
										Values: []*SchemaEnumValue{
											{
												Name: "up",
											},
											{
												Name: "down",
											},
										},
									},
								},
							},
						},
					},
				},
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
					},
					{
						Name: "Component",
						Fields: []*SchemaField{
							{
								Name: "component_name",
								Type: &String{},
							},
						},
					},
					{
						Name: "Machine",
						Fields: []*SchemaField{
							{
								Name: "owner",
								Type: &Reference{T: &SchemaModel{Name: "Operator"}},
							},
						},
					},
				},
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
					},
					{
						Name: "CamelCase",
						Fields: []*SchemaField{
							{
								Name: "foo_bar",
								Type: &String{},
							},
						},
					},
					{
						Name: "TitleCase",
						Fields: []*SchemaField{
							{
								Name: "foo_bar",
								Type: &String{},
							},
						},
					},
					{
						Name: "SnakeCase",
						Fields: []*SchemaField{
							{
								Name: "foo_bar",
								Type: &String{},
							},
						},
					},
				},
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
			name:    "err_ary_embedded_model",
			wantErr: true,
		},
		{
			name:    "err_embedded_model",
			wantErr: true,
		},
		{
			name: "struct",
			want: &Schema{
				Structs: []*SchemaStruct{
					{
						Name:    "Person",
						Comment: "A sample struct",
						Fields: []*SchemaField{
							{
								Name: "display_name",
								Type: &String{},
							},
						},
					},
				},
			},
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
					},
				},
			},
		},
		{
			name: "comments",
			want: &Schema{
				Models: []*SchemaModel{
					{
						Name:    "FooModel",
						Comment: "Model Comments.",
						Fields: []*SchemaField{
							{
								Name:    "cool_field",
								Comment: "Field comment",
								Type:    &String{},
							},
						},
					},
				},
				Structs: []*SchemaStruct{
					{
						Name:    "Foo",
						Comment: "Struct Comment.",
						Fields: []*SchemaField{
							{
								Name:    "cool_field",
								Comment: "Field comment",
								Type:    &String{},
							},
						},
					},
					{
						Name:    "Bar",
						Comment: "Multi-line\n\n Struct comment.",
						Fields: []*SchemaField{
							{
								Name:    "field_two",
								Comment: "Multi-line\nField comment",
								Type:    &String{},
							},
						},
					},
				},
				Enums: []*SchemaEnum{
					{
						Name:    "FooEnum",
						Comment: "Enum Comments.",
						Values: []*SchemaEnumValue{
							{
								Name:            "one",
								Comment:         "Case comment",
								AssociatedValue: nil,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if p := recover(); p != nil && !tt.wantErr {
					t.Fatalf("panic: %s\n\n%s", p, debug.Stack())
				}
			}()
			r, err := os.Open(path.Join("testfixtures", "schema", tt.name+".firemodel"))
			if err != nil {
				t.Fatal(err)
			}
			got, err := ParseSchema(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSchema error: %v", err)
				return
			}

			assert.DeepEqual(t, got, tt.want)
		})
	}
}
