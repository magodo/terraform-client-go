package configschema

import (
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

// Mimic TestBlockImpliedType
func TestSchemaBlockImpliedType(t *testing.T) {
	tests := map[string]struct {
		Schema *tfjson.SchemaBlock
		Want   cty.Type
	}{
		"nil": {
			nil,
			cty.EmptyObject,
		},
		"empty": {
			&tfjson.SchemaBlock{},
			cty.EmptyObject,
		},
		"attributes": {
			&tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"optional": {
						AttributeType: cty.String,
						Optional:      true,
					},
					"required": {
						AttributeType: cty.Number,
						Required:      true,
					},
					"computed": {
						AttributeType: cty.List(cty.Bool),
						Computed:      true,
					},
					"optional_computed": {
						AttributeType: cty.Map(cty.Bool),
						Optional:      true,
						Computed:      true,
					},
				},
			},
			cty.Object(map[string]cty.Type{
				"optional":          cty.String,
				"required":          cty.Number,
				"computed":          cty.List(cty.Bool),
				"optional_computed": cty.Map(cty.Bool),
			}),
		},
		"blocks": {
			&tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"single": {
						NestingMode: tfjson.SchemaNestingModeSingle,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"foo": {
									AttributeType: cty.DynamicPseudoType,
									Required:      true,
								},
							},
						},
					},
					"list": {
						NestingMode: tfjson.SchemaNestingModeList,
					},
					"set": {
						NestingMode: tfjson.SchemaNestingModeSet,
					},
					"map": {
						NestingMode: tfjson.SchemaNestingModeMap,
					},
				},
			},
			cty.Object(map[string]cty.Type{
				"single": cty.Object(map[string]cty.Type{
					"foo": cty.DynamicPseudoType,
				}),
				"list": cty.List(cty.EmptyObject),
				"set":  cty.Set(cty.EmptyObject),
				"map":  cty.Map(cty.EmptyObject),
			}),
		},
		"deep block nesting": {
			&tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"single": {
						NestingMode: tfjson.SchemaNestingModeSingle,
						Block: &tfjson.SchemaBlock{
							NestedBlocks: map[string]*tfjson.SchemaBlockType{
								"list": {
									NestingMode: tfjson.SchemaNestingModeList,
									Block: &tfjson.SchemaBlock{
										NestedBlocks: map[string]*tfjson.SchemaBlockType{
											"set": {
												NestingMode: tfjson.SchemaNestingModeSet,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			cty.Object(map[string]cty.Type{
				"single": cty.Object(map[string]cty.Type{
					"list": cty.List(cty.Object(map[string]cty.Type{
						"set": cty.Set(cty.EmptyObject),
					})),
				}),
			}),
		},
		"nested objects with optional attrs": {
			&tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"map": {
						Optional: true,
						AttributeNestedType: &tfjson.SchemaNestedAttributeType{
							NestingMode: tfjson.SchemaNestingModeMap,
							Attributes: map[string]*tfjson.SchemaAttribute{
								"optional":          {AttributeType: cty.String, Optional: true},
								"required":          {AttributeType: cty.Number, Required: true},
								"computed":          {AttributeType: cty.List(cty.Bool), Computed: true},
								"optional_computed": {AttributeType: cty.Map(cty.Bool), Optional: true, Computed: true},
							},
						},
					},
				},
			},
			// The ImpliedType from the type-level block should not contain any
			// optional attributes.
			cty.Object(map[string]cty.Type{
				"map": cty.Map(cty.Object(
					map[string]cty.Type{
						"optional":          cty.String,
						"required":          cty.Number,
						"computed":          cty.List(cty.Bool),
						"optional_computed": cty.Map(cty.Bool),
					},
				)),
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := SchemaBlockImpliedType(test.Schema)
			if !got.Equals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

// Mimic TestObjectImpliedType
func TestSchemaNestedAttributeTypeImpliedType(t *testing.T) {
	tests := map[string]struct {
		Schema *tfjson.SchemaNestedAttributeType
		Want   cty.Type
	}{
		"nil": {
			nil,
			cty.EmptyObject,
		},
		"empty": {
			&tfjson.SchemaNestedAttributeType{},
			cty.EmptyObject,
		},
		"attributes": {
			&tfjson.SchemaNestedAttributeType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Attributes: map[string]*tfjson.SchemaAttribute{
					"optional":          {AttributeType: cty.String, Optional: true},
					"required":          {AttributeType: cty.Number, Required: true},
					"computed":          {AttributeType: cty.List(cty.Bool), Computed: true},
					"optional_computed": {AttributeType: cty.Map(cty.Bool), Optional: true, Computed: true},
				},
			},
			cty.Object(
				map[string]cty.Type{
					"optional":          cty.String,
					"required":          cty.Number,
					"computed":          cty.List(cty.Bool),
					"optional_computed": cty.Map(cty.Bool),
				},
			),
		},
		"nested attributes": {
			&tfjson.SchemaNestedAttributeType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Attributes: map[string]*tfjson.SchemaAttribute{
					"nested_type": {
						AttributeNestedType: &tfjson.SchemaNestedAttributeType{
							NestingMode: tfjson.SchemaNestingModeSingle,
							Attributes: map[string]*tfjson.SchemaAttribute{
								"optional":          {AttributeType: cty.String, Optional: true},
								"required":          {AttributeType: cty.Number, Required: true},
								"computed":          {AttributeType: cty.List(cty.Bool), Computed: true},
								"optional_computed": {AttributeType: cty.Map(cty.Bool), Optional: true, Computed: true},
							},
						},
						Optional: true,
					},
				},
			},
			cty.Object(map[string]cty.Type{
				"nested_type": cty.Object(map[string]cty.Type{
					"optional":          cty.String,
					"required":          cty.Number,
					"computed":          cty.List(cty.Bool),
					"optional_computed": cty.Map(cty.Bool),
				}),
			}),
		},
		"nested object-type attributes": {
			&tfjson.SchemaNestedAttributeType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Attributes: map[string]*tfjson.SchemaAttribute{
					"nested_type": {
						AttributeNestedType: &tfjson.SchemaNestedAttributeType{
							NestingMode: tfjson.SchemaNestingModeSingle,
							Attributes: map[string]*tfjson.SchemaAttribute{
								"optional":          {AttributeType: cty.String, Optional: true},
								"required":          {AttributeType: cty.Number, Required: true},
								"computed":          {AttributeType: cty.List(cty.Bool), Computed: true},
								"optional_computed": {AttributeType: cty.Map(cty.Bool), Optional: true, Computed: true},
								"object": {
									AttributeType: cty.ObjectWithOptionalAttrs(map[string]cty.Type{
										"optional": cty.String,
										"required": cty.Number,
									}, []string{"optional"}),
								},
							},
						},
						Optional: true,
					},
				},
			},
			cty.Object(map[string]cty.Type{
				"nested_type": cty.Object(map[string]cty.Type{
					"optional":          cty.String,
					"required":          cty.Number,
					"computed":          cty.List(cty.Bool),
					"optional_computed": cty.Map(cty.Bool),
					"object":            cty.Object(map[string]cty.Type{"optional": cty.String, "required": cty.Number}),
				}),
			}),
		},
		"NestingList": {
			&tfjson.SchemaNestedAttributeType{
				NestingMode: tfjson.SchemaNestingModeList,
				Attributes: map[string]*tfjson.SchemaAttribute{
					"foo": {AttributeType: cty.String, Optional: true},
				},
			},
			cty.List(cty.Object(map[string]cty.Type{"foo": cty.String})),
		},
		"NestingMap": {
			&tfjson.SchemaNestedAttributeType{
				NestingMode: tfjson.SchemaNestingModeMap,
				Attributes: map[string]*tfjson.SchemaAttribute{
					"foo": {AttributeType: cty.String},
				},
			},
			cty.Map(cty.Object(map[string]cty.Type{"foo": cty.String})),
		},
		"NestingSet": {
			&tfjson.SchemaNestedAttributeType{
				NestingMode: tfjson.SchemaNestingModeSet,
				Attributes: map[string]*tfjson.SchemaAttribute{
					"foo": {AttributeType: cty.String},
				},
			},
			cty.Set(cty.Object(map[string]cty.Type{"foo": cty.String})),
		},
		"deeply nested NestingList": {
			&tfjson.SchemaNestedAttributeType{
				NestingMode: tfjson.SchemaNestingModeList,
				Attributes: map[string]*tfjson.SchemaAttribute{
					"foo": {
						AttributeNestedType: &tfjson.SchemaNestedAttributeType{
							NestingMode: tfjson.SchemaNestingModeList,
							Attributes: map[string]*tfjson.SchemaAttribute{
								"bar": {AttributeType: cty.String},
							},
						},
					},
				},
			},
			cty.List(cty.Object(map[string]cty.Type{"foo": cty.List(cty.Object(map[string]cty.Type{"bar": cty.String}))})),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := SchemaNestedAttributeTypeImpliedType(test.Schema)
			if !got.Equals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

// Mimic TestObjectSpecType
func TestSchemaNestedAttributeTypeSpecType(t *testing.T) {
	tests := map[string]struct {
		Schema *tfjson.SchemaNestedAttributeType
		Want   cty.Type
	}{
		"attributes": {
			&tfjson.SchemaNestedAttributeType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Attributes: map[string]*tfjson.SchemaAttribute{
					"optional":          {AttributeType: cty.String, Optional: true},
					"required":          {AttributeType: cty.Number, Required: true},
					"computed":          {AttributeType: cty.List(cty.Bool), Computed: true},
					"optional_computed": {AttributeType: cty.Map(cty.Bool), Optional: true, Computed: true},
				},
			},
			cty.ObjectWithOptionalAttrs(
				map[string]cty.Type{
					"optional":          cty.String,
					"required":          cty.Number,
					"computed":          cty.List(cty.Bool),
					"optional_computed": cty.Map(cty.Bool),
				},
				[]string{"optional", "computed", "optional_computed"},
			),
		},
		"nested attributes": {
			&tfjson.SchemaNestedAttributeType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Attributes: map[string]*tfjson.SchemaAttribute{
					"nested_type": {
						AttributeNestedType: &tfjson.SchemaNestedAttributeType{
							NestingMode: tfjson.SchemaNestingModeSingle,
							Attributes: map[string]*tfjson.SchemaAttribute{
								"optional":          {AttributeType: cty.String, Optional: true},
								"required":          {AttributeType: cty.Number, Required: true},
								"computed":          {AttributeType: cty.List(cty.Bool), Computed: true},
								"optional_computed": {AttributeType: cty.Map(cty.Bool), Optional: true, Computed: true},
							},
						},
						Optional: true,
					},
				},
			},
			cty.ObjectWithOptionalAttrs(map[string]cty.Type{
				"nested_type": cty.ObjectWithOptionalAttrs(map[string]cty.Type{
					"optional":          cty.String,
					"required":          cty.Number,
					"computed":          cty.List(cty.Bool),
					"optional_computed": cty.Map(cty.Bool),
				}, []string{"optional", "computed", "optional_computed"}),
			}, []string{"nested_type"}),
		},
		"nested object-type attributes": {
			&tfjson.SchemaNestedAttributeType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Attributes: map[string]*tfjson.SchemaAttribute{
					"nested_type": {
						AttributeNestedType: &tfjson.SchemaNestedAttributeType{
							NestingMode: tfjson.SchemaNestingModeSingle,
							Attributes: map[string]*tfjson.SchemaAttribute{
								"optional":          {AttributeType: cty.String, Optional: true},
								"required":          {AttributeType: cty.Number, Required: true},
								"computed":          {AttributeType: cty.List(cty.Bool), Computed: true},
								"optional_computed": {AttributeType: cty.Map(cty.Bool), Optional: true, Computed: true},
								"object": {
									AttributeType: cty.ObjectWithOptionalAttrs(map[string]cty.Type{
										"optional": cty.String,
										"required": cty.Number,
									}, []string{"optional"}),
								},
							},
						},
						Optional: true,
					},
				},
			},
			cty.ObjectWithOptionalAttrs(map[string]cty.Type{
				"nested_type": cty.ObjectWithOptionalAttrs(map[string]cty.Type{
					"optional":          cty.String,
					"required":          cty.Number,
					"computed":          cty.List(cty.Bool),
					"optional_computed": cty.Map(cty.Bool),
					"object":            cty.ObjectWithOptionalAttrs(map[string]cty.Type{"optional": cty.String, "required": cty.Number}, []string{"optional"}),
				}, []string{"optional", "computed", "optional_computed"}),
			}, []string{"nested_type"}),
		},
		"NestingList": {
			&tfjson.SchemaNestedAttributeType{
				NestingMode: tfjson.SchemaNestingModeList,
				Attributes: map[string]*tfjson.SchemaAttribute{
					"foo": {AttributeType: cty.String, Optional: true},
				},
			},
			cty.List(cty.ObjectWithOptionalAttrs(map[string]cty.Type{"foo": cty.String}, []string{"foo"})),
		},
		"NestingMap": {
			&tfjson.SchemaNestedAttributeType{
				NestingMode: tfjson.SchemaNestingModeMap,
				Attributes: map[string]*tfjson.SchemaAttribute{
					"foo": {AttributeType: cty.String},
				},
			},
			cty.Map(cty.Object(map[string]cty.Type{"foo": cty.String})),
		},
		"NestingSet": {
			&tfjson.SchemaNestedAttributeType{
				NestingMode: tfjson.SchemaNestingModeSet,
				Attributes: map[string]*tfjson.SchemaAttribute{
					"foo": {AttributeType: cty.String},
				},
			},
			cty.Set(cty.Object(map[string]cty.Type{"foo": cty.String})),
		},
		"deeply nested NestingList": {
			&tfjson.SchemaNestedAttributeType{
				NestingMode: tfjson.SchemaNestingModeList,
				Attributes: map[string]*tfjson.SchemaAttribute{
					"foo": {
						AttributeNestedType: &tfjson.SchemaNestedAttributeType{
							NestingMode: tfjson.SchemaNestingModeList,
							Attributes: map[string]*tfjson.SchemaAttribute{
								"bar": {AttributeType: cty.String},
							},
						},
					},
				},
			},
			cty.List(cty.Object(map[string]cty.Type{"foo": cty.List(cty.Object(map[string]cty.Type{"bar": cty.String}))})),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := schemaNestedAttributeTypeSpecType(test.Schema)
			if !got.Equals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}
