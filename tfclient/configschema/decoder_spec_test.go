package configschema

import (
	"sort"
	"testing"

	"github.com/apparentlymart/go-dump/dump"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hcltest"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

func TestBlockDecoderSpec(t *testing.T) {
	tests := map[string]struct {
		Schema    *tfjson.SchemaBlock
		TestBody  hcl.Body
		Want      cty.Value
		DiagCount int
	}{
		"empty": {
			&tfjson.SchemaBlock{},
			hcl.EmptyBody(),
			cty.EmptyObjectVal,
			0,
		},
		"nil": {
			nil,
			hcl.EmptyBody(),
			cty.EmptyObjectVal,
			0,
		},
		"attributes": {
			&tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"optional": {
						AttributeType: cty.Number,
						Optional:      true,
					},
					"required": {
						AttributeType: cty.String,
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
					"optional_computed_overridden": {
						AttributeType: cty.Bool,
						Optional:      true,
						Computed:      true,
					},
					"optional_computed_unknown": {
						AttributeType: cty.String,
						Optional:      true,
						Computed:      true,
					},
				},
			},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"required": {
						Name: "required",
						Expr: hcltest.MockExprLiteral(cty.NumberIntVal(5)),
					},
					"optional_computed_overridden": {
						Name: "optional_computed_overridden",
						Expr: hcltest.MockExprLiteral(cty.True),
					},
					"optional_computed_unknown": {
						Name: "optional_computed_overridden",
						Expr: hcltest.MockExprLiteral(cty.UnknownVal(cty.String)),
					},
				},
			}),
			cty.ObjectVal(map[string]cty.Value{
				"optional":                     cty.NullVal(cty.Number),
				"required":                     cty.StringVal("5"), // converted from number to string
				"computed":                     cty.NullVal(cty.List(cty.Bool)),
				"optional_computed":            cty.NullVal(cty.Map(cty.Bool)),
				"optional_computed_overridden": cty.True,
				"optional_computed_unknown":    cty.UnknownVal(cty.String),
			}),
			0,
		},
		"dynamically-typed attribute": {
			&tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"foo": {
						AttributeType: cty.DynamicPseudoType, // any type is permitted
						Required:      true,
					},
				},
			},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"foo": {
						Name: "foo",
						Expr: hcltest.MockExprLiteral(cty.True),
					},
				},
			}),
			cty.ObjectVal(map[string]cty.Value{
				"foo": cty.True,
			}),
			0,
		},
		"dynamically-typed attribute omitted": {
			&tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"foo": {
						AttributeType: cty.DynamicPseudoType, // any type is permitted
						Optional:      true,
					},
				},
			},
			hcltest.MockBody(&hcl.BodyContent{}),
			cty.ObjectVal(map[string]cty.Value{
				"foo": cty.NullVal(cty.DynamicPseudoType),
			}),
			0,
		},
		"required attribute omitted": {
			&tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"foo": {
						AttributeType: cty.Bool,
						Required:      true,
					},
				},
			},
			hcltest.MockBody(&hcl.BodyContent{}),
			cty.ObjectVal(map[string]cty.Value{
				"foo": cty.NullVal(cty.Bool),
			}),
			1, // missing required attribute
		},
		"wrong attribute type": {
			&tfjson.SchemaBlock{
				Attributes: map[string]*tfjson.SchemaAttribute{
					"optional": {
						AttributeType: cty.Number,
						Optional:      true,
					},
				},
			},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"optional": {
						Name: "optional",
						Expr: hcltest.MockExprLiteral(cty.True),
					},
				},
			}),
			cty.ObjectVal(map[string]cty.Value{
				"optional": cty.UnknownVal(cty.Number),
			}),
			1, // incorrect type; number required
		},
		"blocks": {
			&tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"single": {
						NestingMode: tfjson.SchemaNestingModeSingle,
						Block:       &tfjson.SchemaBlock{},
					},
					"list": {
						NestingMode: tfjson.SchemaNestingModeList,
						Block:       &tfjson.SchemaBlock{},
					},
					"set": {
						NestingMode: tfjson.SchemaNestingModeSet,
						Block:       &tfjson.SchemaBlock{},
					},
					"map": {
						NestingMode: tfjson.SchemaNestingModeMap,
						Block:       &tfjson.SchemaBlock{},
					},
				},
			},
			hcltest.MockBody(&hcl.BodyContent{
				Blocks: hcl.Blocks{
					&hcl.Block{
						Type: "list",
						Body: hcl.EmptyBody(),
					},
					&hcl.Block{
						Type: "single",
						Body: hcl.EmptyBody(),
					},
					&hcl.Block{
						Type: "list",
						Body: hcl.EmptyBody(),
					},
					&hcl.Block{
						Type: "set",
						Body: hcl.EmptyBody(),
					},
					&hcl.Block{
						Type:        "map",
						Labels:      []string{"foo"},
						LabelRanges: []hcl.Range{{}},
						Body:        hcl.EmptyBody(),
					},
					&hcl.Block{
						Type:        "map",
						Labels:      []string{"bar"},
						LabelRanges: []hcl.Range{{}},
						Body:        hcl.EmptyBody(),
					},
					&hcl.Block{
						Type: "set",
						Body: hcl.EmptyBody(),
					},
				},
			}),
			cty.ObjectVal(map[string]cty.Value{
				"single": cty.EmptyObjectVal,
				"list": cty.ListVal([]cty.Value{
					cty.EmptyObjectVal,
					cty.EmptyObjectVal,
				}),
				"set": cty.SetVal([]cty.Value{
					cty.EmptyObjectVal,
					cty.EmptyObjectVal,
				}),
				"map": cty.MapVal(map[string]cty.Value{
					"foo": cty.EmptyObjectVal,
					"bar": cty.EmptyObjectVal,
				}),
			}),
			0,
		},
		"blocks with dynamically-typed attributes": {
			&tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"single": {
						NestingMode: tfjson.SchemaNestingModeSingle,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"a": {
									AttributeType: cty.DynamicPseudoType,
									Optional:      true,
								},
							},
						},
					},
					"list": {
						NestingMode: tfjson.SchemaNestingModeList,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"a": {
									AttributeType: cty.DynamicPseudoType,
									Optional:      true,
								},
							},
						},
					},
					"map": {
						NestingMode: tfjson.SchemaNestingModeMap,
						Block: &tfjson.SchemaBlock{
							Attributes: map[string]*tfjson.SchemaAttribute{
								"a": {
									AttributeType: cty.DynamicPseudoType,
									Optional:      true,
								},
							},
						},
					},
				},
			},
			hcltest.MockBody(&hcl.BodyContent{
				Blocks: hcl.Blocks{
					&hcl.Block{
						Type: "list",
						Body: hcl.EmptyBody(),
					},
					&hcl.Block{
						Type: "single",
						Body: hcl.EmptyBody(),
					},
					&hcl.Block{
						Type: "list",
						Body: hcl.EmptyBody(),
					},
					&hcl.Block{
						Type:        "map",
						Labels:      []string{"foo"},
						LabelRanges: []hcl.Range{{}},
						Body:        hcl.EmptyBody(),
					},
					&hcl.Block{
						Type:        "map",
						Labels:      []string{"bar"},
						LabelRanges: []hcl.Range{{}},
						Body:        hcl.EmptyBody(),
					},
				},
			}),
			cty.ObjectVal(map[string]cty.Value{
				"single": cty.ObjectVal(map[string]cty.Value{
					"a": cty.NullVal(cty.DynamicPseudoType),
				}),
				"list": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"a": cty.NullVal(cty.DynamicPseudoType),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"a": cty.NullVal(cty.DynamicPseudoType),
					}),
				}),
				"map": cty.ObjectVal(map[string]cty.Value{
					"foo": cty.ObjectVal(map[string]cty.Value{
						"a": cty.NullVal(cty.DynamicPseudoType),
					}),
					"bar": cty.ObjectVal(map[string]cty.Value{
						"a": cty.NullVal(cty.DynamicPseudoType),
					}),
				}),
			}),
			0,
		},
		"too many list items": {
			&tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"foo": {
						NestingMode: tfjson.SchemaNestingModeList,
						Block:       &tfjson.SchemaBlock{},
						MaxItems:    1,
					},
				},
			},
			hcltest.MockBody(&hcl.BodyContent{
				Blocks: hcl.Blocks{
					&hcl.Block{
						Type: "foo",
						Body: hcl.EmptyBody(),
					},
					&hcl.Block{
						Type: "foo",
						Body: unknownBody{hcl.EmptyBody()},
					},
				},
			}),
			cty.ObjectVal(map[string]cty.Value{
				"foo": cty.UnknownVal(cty.List(cty.EmptyObject)),
			}),
			0, // max items cannot be validated during decode
		},
		// dynamic blocks may fulfill MinItems, but there is only one block to
		// decode.
		"required MinItems": {
			&tfjson.SchemaBlock{
				NestedBlocks: map[string]*tfjson.SchemaBlockType{
					"foo": {
						NestingMode: tfjson.SchemaNestingModeList,
						Block:       &tfjson.SchemaBlock{},
						MinItems:    2,
					},
				},
			},
			hcltest.MockBody(&hcl.BodyContent{
				Blocks: hcl.Blocks{
					&hcl.Block{
						Type: "foo",
						Body: unknownBody{hcl.EmptyBody()},
					},
				},
			}),
			cty.ObjectVal(map[string]cty.Value{
				"foo": cty.UnknownVal(cty.List(cty.EmptyObject)),
			}),
			0,
		},
		"extraneous attribute": {
			&tfjson.SchemaBlock{},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"extra": {
						Name: "extra",
						Expr: hcltest.MockExprLiteral(cty.StringVal("hello")),
					},
				},
			}),
			cty.EmptyObjectVal,
			1, // extraneous attribute
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			spec := DecoderSpec(test.Schema)

			got, diags := hcldec.Decode(test.TestBody, spec, nil)
			if len(diags) != test.DiagCount {
				t.Errorf("wrong number of diagnostics %d; want %d", len(diags), test.DiagCount)
				for _, diag := range diags {
					t.Logf("- %s", diag.Error())
				}
			}

			if !got.RawEquals(test.Want) {
				t.Logf("[INFO] implied schema is %s", spew.Sdump(hcldec.ImpliedSchema(spec)))
				t.Errorf("wrong result\ngot:  %s\nwant: %s", dump.Value(got), dump.Value(test.Want))
			}

			// Double-check that we're producing consistent results for DecoderSpec
			// and ImpliedType.
			impliedType := SchemaBlockImpliedType(test.Schema)
			if errs := got.Type().TestConformance(impliedType); len(errs) != 0 {
				t.Errorf("result does not conform to the schema's implied type")
				for _, err := range errs {
					t.Logf("- %s", err.Error())
				}
			}
		})
	}
}

// this satisfies hcldec.UnknownBody to simulate a dynamic block with an
// unknown number of values.
type unknownBody struct {
	hcl.Body
}

func (b unknownBody) Unknown() bool {
	return true
}

// Mimic TestAttributeDecoderSpec
func TestSchemaAttributeDecoderSpec(t *testing.T) {
	tests := map[string]struct {
		Schema    *tfjson.SchemaAttribute
		TestBody  hcl.Body
		Want      cty.Value
		DiagCount int
	}{
		"empty": {
			&tfjson.SchemaAttribute{},
			hcl.EmptyBody(),
			cty.NilVal,
			0,
		},
		"nil": {
			nil,
			hcl.EmptyBody(),
			cty.NilVal,
			0,
		},
		"optional string (null)": {
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Optional:      true,
			},
			hcltest.MockBody(&hcl.BodyContent{}),
			cty.NullVal(cty.String),
			0,
		},
		"optional string": {
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Optional:      true,
			},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"attr": {
						Name: "attr",
						Expr: hcltest.MockExprLiteral(cty.StringVal("bar")),
					},
				},
			}),
			cty.StringVal("bar"),
			0,
		},
		"NestedType with required string": {
			&tfjson.SchemaAttribute{
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeSingle,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
						},
					},
				},
				Optional: true,
			},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"attr": {
						Name: "attr",
						Expr: hcltest.MockExprLiteral(cty.ObjectVal(map[string]cty.Value{
							"foo": cty.StringVal("bar"),
						})),
					},
				},
			}),
			cty.ObjectVal(map[string]cty.Value{
				"foo": cty.StringVal("bar"),
			}),
			0,
		},
		"NestedType with optional attributes": {
			&tfjson.SchemaAttribute{
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeSingle,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Optional:      true,
						},
						"bar": {
							AttributeType: cty.String,
							Optional:      true,
						},
					},
				},
				Optional: true,
			},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"attr": {
						Name: "attr",
						Expr: hcltest.MockExprLiteral(cty.ObjectVal(map[string]cty.Value{
							"foo": cty.StringVal("bar"),
						})),
					},
				},
			}),
			cty.ObjectVal(map[string]cty.Value{
				"foo": cty.StringVal("bar"),
				"bar": cty.NullVal(cty.String),
			}),
			0,
		},
		"NestedType with missing required string": {
			&tfjson.SchemaAttribute{
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeSingle,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
						},
					},
				},
				Optional: true,
			},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"attr": {
						Name: "attr",
						Expr: hcltest.MockExprLiteral(cty.EmptyObjectVal),
					},
				},
			}),
			cty.UnknownVal(cty.Object(map[string]cty.Type{
				"foo": cty.String,
			})),
			1,
		},
		// NestedModes
		"NestedType NestingModeList valid": {
			&tfjson.SchemaAttribute{
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeList,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
						},
					},
				},
				Optional: true,
			},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"attr": {
						Name: "attr",
						Expr: hcltest.MockExprLiteral(cty.ListVal([]cty.Value{
							cty.ObjectVal(map[string]cty.Value{
								"foo": cty.StringVal("bar"),
							}),
							cty.ObjectVal(map[string]cty.Value{
								"foo": cty.StringVal("baz"),
							}),
						})),
					},
				},
			}),
			cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{"foo": cty.StringVal("bar")}),
				cty.ObjectVal(map[string]cty.Value{"foo": cty.StringVal("baz")}),
			}),
			0,
		},
		"NestedType NestingModeList invalid": {
			&tfjson.SchemaAttribute{
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeList,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
						},
					},
				},
				Optional: true,
			},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"attr": {
						Name: "attr",
						Expr: hcltest.MockExprLiteral(cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
							// "foo" should be a string, not a list
							"foo": cty.ListVal([]cty.Value{cty.StringVal("bar"), cty.StringVal("baz")}),
						})})),
					},
				},
			}),
			cty.UnknownVal(cty.List(cty.Object(map[string]cty.Type{"foo": cty.String}))),
			1,
		},
		"NestedType NestingModeSet valid": {
			&tfjson.SchemaAttribute{
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeSet,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
						},
					},
				},
				Optional: true,
			},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"attr": {
						Name: "attr",
						Expr: hcltest.MockExprLiteral(cty.SetVal([]cty.Value{
							cty.ObjectVal(map[string]cty.Value{
								"foo": cty.StringVal("bar"),
							}),
							cty.ObjectVal(map[string]cty.Value{
								"foo": cty.StringVal("baz"),
							}),
						})),
					},
				},
			}),
			cty.SetVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{"foo": cty.StringVal("bar")}),
				cty.ObjectVal(map[string]cty.Value{"foo": cty.StringVal("baz")}),
			}),
			0,
		},
		"NestedType NestingModeSet invalid": {
			&tfjson.SchemaAttribute{
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeSet,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
						},
					},
				},
				Optional: true,
			},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"attr": {
						Name: "attr",
						Expr: hcltest.MockExprLiteral(cty.SetVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
							// "foo" should be a string, not a list
							"foo": cty.ListVal([]cty.Value{cty.StringVal("bar"), cty.StringVal("baz")}),
						})})),
					},
				},
			}),
			cty.UnknownVal(cty.Set(cty.Object(map[string]cty.Type{"foo": cty.String}))),
			1,
		},
		"NestedType NestingModeMap valid": {
			&tfjson.SchemaAttribute{
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeMap,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
						},
					},
				},
				Optional: true,
			},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"attr": {
						Name: "attr",
						Expr: hcltest.MockExprLiteral(cty.MapVal(map[string]cty.Value{
							"one": cty.ObjectVal(map[string]cty.Value{
								"foo": cty.StringVal("bar"),
							}),
							"two": cty.ObjectVal(map[string]cty.Value{
								"foo": cty.StringVal("baz"),
							}),
						})),
					},
				},
			}),
			cty.MapVal(map[string]cty.Value{
				"one": cty.ObjectVal(map[string]cty.Value{"foo": cty.StringVal("bar")}),
				"two": cty.ObjectVal(map[string]cty.Value{"foo": cty.StringVal("baz")}),
			}),
			0,
		},
		"NestedType NestingModeMap invalid": {
			&tfjson.SchemaAttribute{
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeMap,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
						},
					},
				},
				Optional: true,
			},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"attr": {
						Name: "attr",
						Expr: hcltest.MockExprLiteral(cty.MapVal(map[string]cty.Value{
							"one": cty.ObjectVal(map[string]cty.Value{
								// "foo" should be a string, not a list
								"foo": cty.ListVal([]cty.Value{cty.StringVal("bar"), cty.StringVal("baz")}),
							}),
						})),
					},
				},
			}),
			cty.UnknownVal(cty.Map(cty.Object(map[string]cty.Type{"foo": cty.String}))),
			1,
		},
		"deeply NestedType NestingModeList valid": {
			&tfjson.SchemaAttribute{
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeList,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeNestedType: &tfjson.SchemaNestedAttributeType{
								NestingMode: tfjson.SchemaNestingModeList,
								Attributes: map[string]*tfjson.SchemaAttribute{
									"bar": {
										AttributeType: cty.String,
										Required:      true,
									},
								},
							},
							Required: true,
						},
					},
				},
				Optional: true,
			},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"attr": {
						Name: "attr",
						Expr: hcltest.MockExprLiteral(cty.ListVal([]cty.Value{
							cty.ObjectVal(map[string]cty.Value{
								"foo": cty.ListVal([]cty.Value{
									cty.ObjectVal(map[string]cty.Value{"bar": cty.StringVal("baz")}),
									cty.ObjectVal(map[string]cty.Value{"bar": cty.StringVal("boz")}),
								}),
							}),
							cty.ObjectVal(map[string]cty.Value{
								"foo": cty.ListVal([]cty.Value{
									cty.ObjectVal(map[string]cty.Value{"bar": cty.StringVal("biz")}),
									cty.ObjectVal(map[string]cty.Value{"bar": cty.StringVal("buz")}),
								}),
							}),
						})),
					},
				},
			}),
			cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{"foo": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{"bar": cty.StringVal("baz")}),
					cty.ObjectVal(map[string]cty.Value{"bar": cty.StringVal("boz")}),
				})}),
				cty.ObjectVal(map[string]cty.Value{"foo": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{"bar": cty.StringVal("biz")}),
					cty.ObjectVal(map[string]cty.Value{"bar": cty.StringVal("buz")}),
				})}),
			}),
			0,
		},
		"deeply NestedType NestingList invalid": {
			&tfjson.SchemaAttribute{
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeList,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeNestedType: &tfjson.SchemaNestedAttributeType{
								NestingMode: tfjson.SchemaNestingModeList,
								Attributes: map[string]*tfjson.SchemaAttribute{
									"bar": {
										AttributeType: cty.Number,
										Required:      true,
									},
								},
							},
							Required: true,
						},
					},
				},
				Optional: true,
			},
			hcltest.MockBody(&hcl.BodyContent{
				Attributes: hcl.Attributes{
					"attr": {
						Name: "attr",
						Expr: hcltest.MockExprLiteral(cty.ListVal([]cty.Value{
							cty.ObjectVal(map[string]cty.Value{
								"foo": cty.ListVal([]cty.Value{
									// bar should be a Number
									cty.ObjectVal(map[string]cty.Value{"bar": cty.True}),
									cty.ObjectVal(map[string]cty.Value{"bar": cty.False}),
								}),
							}),
						})),
					},
				},
			}),
			cty.UnknownVal(cty.List(cty.Object(map[string]cty.Type{
				"foo": cty.List(cty.Object(map[string]cty.Type{"bar": cty.Number})),
			}))),
			1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			spec := decoderSpec(test.Schema, "attr")
			got, diags := hcldec.Decode(test.TestBody, spec, nil)
			if len(diags) != test.DiagCount {
				t.Errorf("wrong number of diagnostics %d; want %d", len(diags), test.DiagCount)
				for _, diag := range diags {
					t.Logf("- %s", diag.Error())
				}
			}

			if !got.RawEquals(test.Want) {
				t.Logf("[INFO] implied schema is %s", spew.Sdump(hcldec.ImpliedSchema(spec)))
				t.Errorf("wrong result\ngot:  %s\nwant: %s", dump.Value(got), dump.Value(test.Want))
			}
		})
	}
}

// Mimic TestAttributeDecoderSpec_panic
func TestSchemaAttributeDecoderSpec_panic(t *testing.T) {
	attrS := &tfjson.SchemaAttribute{
		AttributeType: cty.Object(map[string]cty.Type{
			"nested_attribute": cty.String,
		}),
		AttributeNestedType: &tfjson.SchemaNestedAttributeType{},
		Optional:            true,
	}

	defer func() { recover() }()
	decoderSpec(attrS, "attr")
	t.Errorf("expected panic")
}

// Mimic TestListOptionalAttrsFromObject
func TestListOptionalAttrsFromSchemaNestedAttributeType(t *testing.T) {
	tests := []struct {
		input *tfjson.SchemaNestedAttributeType
		want  []string
	}{
		{
			nil,
			[]string{},
		},
		{
			&tfjson.SchemaNestedAttributeType{},
			[]string{},
		},
		{
			&tfjson.SchemaNestedAttributeType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Attributes: map[string]*tfjson.SchemaAttribute{
					"optional":          {AttributeType: cty.String, Optional: true},
					"required":          {AttributeType: cty.Number, Required: true},
					"computed":          {AttributeType: cty.List(cty.Bool), Computed: true},
					"optional_computed": {AttributeType: cty.Map(cty.Bool), Optional: true, Computed: true},
				},
			},
			[]string{"optional", "computed", "optional_computed"},
		},
	}

	for _, test := range tests {
		got := listOptionalAttrsFromObject(test.input)

		// order is irrelevant
		sort.Strings(got)
		sort.Strings(test.want)

		if diff := cmp.Diff(got, test.want); diff != "" {
			t.Fatalf("wrong result: %s\n", diff)
		}
	}
}
