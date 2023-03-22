package configschema

import (
	"fmt"
	"testing"

	"github.com/apparentlymart/go-dump/dump"
	"github.com/davecgh/go-spew/spew"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

// Mimic TestBlockEmptyValue
func TestSchemaBlockEmptyValue(t *testing.T) {
	tests := []struct {
		Schema *tfjson.SchemaBlock
		Want   cty.Value
	}{
		{
			&tfjson.SchemaBlock{},
			cty.EmptyObjectVal,
		},
		{
			&tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"str": {AttributeType: cty.String, Required: true},
				},
			},
			cty.ObjectVal(map[string]cty.Value{
				"str": cty.NullVal(cty.String),
			}),
		},
		{
			&tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"single": {
						NestingMode: tfjson.SchemaNestingModeSingle,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"str": {AttributeType: cty.String, Required: true},
							},
						},
					},
				},
			},
			cty.ObjectVal(map[string]cty.Value{
				"single": cty.NullVal(cty.Object(map[string]cty.Type{
					"str": cty.String,
				})),
			}),
		},
		{
			&tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"group": {
						NestingMode: tfjson.SchemaNestingModeGroup,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"str": {AttributeType: cty.String, Required: true},
							},
						},
					},
				},
			},
			cty.ObjectVal(map[string]cty.Value{
				"group": cty.ObjectVal(map[string]cty.Value{
					"str": cty.NullVal(cty.String),
				}),
			}),
		},
		{
			&tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"list": {
						NestingMode: tfjson.SchemaNestingModeList,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"str": {AttributeType: cty.String, Required: true},
							},
						},
					},
				},
			},
			cty.ObjectVal(map[string]cty.Value{
				"list": cty.ListValEmpty(cty.Object(map[string]cty.Type{
					"str": cty.String,
				})),
			}),
		},
		{
			&tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"list_dynamic": {
						NestingMode: tfjson.SchemaNestingModeList,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"str": {AttributeType: cty.DynamicPseudoType, Required: true},
							},
						},
					},
				},
			},
			cty.ObjectVal(map[string]cty.Value{
				"list_dynamic": cty.EmptyTupleVal,
			}),
		},
		{
			&tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"map": {
						NestingMode: tfjson.SchemaNestingModeMap,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"str": {AttributeType: cty.String, Required: true},
							},
						},
					},
				},
			},
			cty.ObjectVal(map[string]cty.Value{
				"map": cty.MapValEmpty(cty.Object(map[string]cty.Type{
					"str": cty.String,
				})),
			}),
		},
		{
			&tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"map_dynamic": {
						NestingMode: tfjson.SchemaNestingModeMap,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"str": {AttributeType: cty.DynamicPseudoType, Required: true},
							},
						},
					},
				},
			},
			cty.ObjectVal(map[string]cty.Value{
				"map_dynamic": cty.EmptyObjectVal,
			}),
		},
		{
			&tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"set": {
						NestingMode: tfjson.SchemaNestingModeSet,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"str": {AttributeType: cty.String, Required: true},
							},
						},
					},
				},
			},
			cty.ObjectVal(map[string]cty.Value{
				"set": cty.SetValEmpty(cty.Object(map[string]cty.Type{
					"str": cty.String,
				})),
			}),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.Schema), func(t *testing.T) {
			got := SchemaBlockEmptyValue(test.Schema)
			if !test.Want.RawEquals(got) {
				t.Errorf("wrong result\nschema: %s\ngot: %s\nwant: %s", spew.Sdump(test.Schema), dump.Value(got), dump.Value(test.Want))
			}
		})
	}
}

// Mimic TestAttributeEmptyValue
func TestSchemaAttributeEmptyValue(t *testing.T) {
	tests := []struct {
		Schema *tfjson.SchemaAttribute
		Want   cty.Value
	}{
		{
			&tfjson.SchemaAttribute{},
			cty.NilVal,
		},
		{
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
			},
			cty.NullVal(cty.String),
		},
		{
			&tfjson.SchemaAttribute{
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeSingle,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"str": {AttributeType: cty.String, Required: true},
					},
				},
			},
			cty.NullVal(cty.Object(map[string]cty.Type{
				"str": cty.String,
			})),
		},
		{
			&tfjson.SchemaAttribute{
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeList,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"str": {AttributeType: cty.String, Required: true},
					},
				},
			},
			cty.NullVal(cty.List(
				cty.Object(map[string]cty.Type{
					"str": cty.String,
				}),
			)),
		},
		{
			&tfjson.SchemaAttribute{
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeMap,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"str": {AttributeType: cty.String, Required: true},
					},
				},
			},
			cty.NullVal(cty.Map(
				cty.Object(map[string]cty.Type{
					"str": cty.String,
				}),
			)),
		},
		{
			&tfjson.SchemaAttribute{
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeSet,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"str": {AttributeType: cty.String, Required: true},
					},
				},
			},
			cty.NullVal(cty.Set(
				cty.Object(map[string]cty.Type{
					"str": cty.String,
				}),
			)),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.Schema), func(t *testing.T) {
			got := SchemaAttributeEmptyValue(test.Schema)
			if !test.Want.RawEquals(got) {
				t.Errorf("wrong result\nschema: %s\ngot: %s\nwant: %s", spew.Sdump(test.Schema), dump.Value(got), dump.Value(test.Want))
			}
		})
	}
}
