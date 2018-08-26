package firemodel

import (
	"testing"
	"strings"
	"github.com/go-test/deep"
	"github.com/davecgh/go-spew/spew"
)

func TestParseSchema(t *testing.T) {
	tests := []struct {
		name    string
		schema  string
		want    *Schema
		wantErr bool
	}{
		{
			name:   "empty",
			schema: ``,
			want: &Schema{
				Options: SchemaOptions{},
			},
		},
		{
			name:   "empty model",
			schema: `model Empty {}`,
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
			schema: `
model SimpleModel {
  string foo;
}
`,
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
			schema: `
// A Test is a test model.
model TestModel {
	// The name.
	string name;
    // The age.
	integer age;
    // The number pi.
    double pi;
    // The birth date.
    timestamp birthdate;
    // True if it is good.
    boolean is_good;
	
    bytes data;
	
    reference friend;
    geopoint location;
    array colors;
    map meta;

	// Fake types...
	File aFile;
	URL anURL;
}
`,
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
			name: "schemaExtras",
			schema: `
model TestModel {
    reference<TestModel> other;
    reference unspecified_other;
    array<string> str_ary;
    array<TestModel> model_ary;
    map<string> str_map;
    map<TestModel> model_map;
}
`,
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
			schema: `
model TestModel {
    URL url;
}
`,
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
			schema: `
// A cardinal direction.
enum Direction {
	// Leftwards.
	left,
	right,
	up,
	down,
}

model TestModel {
	// The direction.
    Direction dir;
}
`,
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
			name: "error: nonsense",
			schema: `
oijasef oijasef
ijef98 aw3raw 3f98asjf oaoeifj 
}
`,
			wantErr: true,
		},
		{
			name: "relational",
			schema: `
model Operator {
	string operator_name;
}
model Component {
	string component_name;
}
model Machine {
	reference<Operator> owner;
	collection<Component> components;
	Component embedded_component;
}
`,
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
			name: "language options",
			schema: `
option langname.optname = "optval";
option otherlang.str = "hey";
option otherlang.number = 1;
option otherlang.truth = true;
option otherlang.falseness = false;
option otherlang.nullness = null;
`,
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
			name: "model options",
			schema: `
model AnnotatedModel {
	option lang.opt = "great";
}
`,
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
			schema: `
model NormalCase {
  string foo_bar;
}
model camelCase {
  string fooBar;
}
model TitleCase {
  string FooBar;
}
model snake_case {
  string foo_bar;
}
`,
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
			name: "syntaxNonsense2",
			schema: `
model {
	// missing title
}
`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.schema)
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
