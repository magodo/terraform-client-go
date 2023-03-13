// Duplicating the logic from: github.com/hashicorp/terraform/internal/plugin6/convert (15ecdb66c84cd8202b0ae3d34c44cb4bbece5444)

package convert

import (
	"encoding/json"
	"reflect"
	"sort"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/zclconf/go-cty/cty"
)

func ConfigSchemaToProto(b *tfjson.SchemaBlock) *tfprotov6.SchemaBlock {
	block := &tfprotov6.SchemaBlock{
		Description:     b.Description,
		DescriptionKind: protoStringKind(b.DescriptionKind),
		Deprecated:      b.Deprecated,
	}

	for _, name := range sortedKeys(b.Attributes) {
		a := b.Attributes[name]

		attr := &tfprotov6.SchemaAttribute{
			Name:            name,
			Description:     a.Description,
			DescriptionKind: protoStringKind(a.DescriptionKind),
			Optional:        a.Optional,
			Computed:        a.Computed,
			Required:        a.Required,
			Sensitive:       a.Sensitive,
			Deprecated:      a.Deprecated,
		}

		if a.AttributeType != cty.NilType {
			ty, err := json.Marshal(a.AttributeType)
			if err != nil {
				panic(err)
			}
			tftype, err := tftypes.ParseJSONType(ty)
			if err != nil {
				panic(err)
			}
			attr.Type = tftype
		}

		if a.AttributeNestedType != nil {
			attr.NestedType = tfjsonObjectToProto(a.AttributeNestedType)
		}

		block.Attributes = append(block.Attributes, attr)
	}

	for _, name := range sortedKeys(b.NestedBlocks) {
		b := b.NestedBlocks[name]
		block.BlockTypes = append(block.BlockTypes, protoSchemaNestedBlock(name, b))
	}

	return block
}

func protoStringKind(k tfjson.SchemaDescriptionKind) tfprotov6.StringKind {
	switch k {
	default:
		return tfprotov6.StringKindPlain
	case tfjson.SchemaDescriptionKindMarkdown:
		return tfprotov6.StringKindMarkdown
	}
}

func protoSchemaNestedBlock(name string, b *tfjson.SchemaBlockType) *tfprotov6.SchemaNestedBlock {
	var nesting tfprotov6.SchemaNestedBlockNestingMode
	switch b.NestingMode {
	case tfjson.SchemaNestingModeSingle:
		nesting = tfprotov6.SchemaNestedBlockNestingModeSingle
	case tfjson.SchemaNestingModeGroup:
		nesting = tfprotov6.SchemaNestedBlockNestingModeGroup
	case tfjson.SchemaNestingModeList:
		nesting = tfprotov6.SchemaNestedBlockNestingModeList
	case tfjson.SchemaNestingModeSet:
		nesting = tfprotov6.SchemaNestedBlockNestingModeSet
	case tfjson.SchemaNestingModeMap:
		nesting = tfprotov6.SchemaNestedBlockNestingModeMap
	default:
		nesting = tfprotov6.SchemaNestedBlockNestingModeInvalid
	}
	return &tfprotov6.SchemaNestedBlock{
		TypeName: name,
		Block:    ConfigSchemaToProto(b.Block),
		Nesting:  nesting,
		MinItems: int64(b.MinItems),
		MaxItems: int64(b.MaxItems),
	}
}

func ProtoToProviderSchema(s *tfprotov6.Schema) tfjson.Schema {
	return tfjson.Schema{
		Version: uint64(s.Version),
		Block:   ProtoToConfigSchema(s.Block),
	}
}

func ProtoToConfigSchema(b *tfprotov6.SchemaBlock) *tfjson.SchemaBlock {
	block := &tfjson.SchemaBlock{
		Attributes:   make(map[string]*tfjson.SchemaAttribute),
		NestedBlocks: make(map[string]*tfjson.SchemaBlockType),

		Description:     b.Description,
		DescriptionKind: schemaStringKind(b.DescriptionKind),
		Deprecated:      b.Deprecated,
	}

	for _, a := range b.Attributes {
		attr := &tfjson.SchemaAttribute{
			Description:     a.Description,
			DescriptionKind: schemaStringKind(a.DescriptionKind),
			Required:        a.Required,
			Optional:        a.Optional,
			Computed:        a.Computed,
			Sensitive:       a.Sensitive,
			Deprecated:      a.Deprecated,
		}

		if a.Type != nil {
			b, err := a.Type.MarshalJSON()
			if err != nil {
				panic(err)
			}
			if err := json.Unmarshal(b, &attr.AttributeType); err != nil {
				panic(err)
			}
		}

		if a.NestedType != nil {
			attr.AttributeNestedType = protoObjectToConfigSchema(a.NestedType)
		}

		block.Attributes[a.Name] = attr
	}

	for _, b := range b.BlockTypes {
		block.NestedBlocks[b.TypeName] = schemaNestedBlock(b)
	}

	return block
}

func schemaStringKind(k tfprotov6.StringKind) tfjson.SchemaDescriptionKind {
	switch k {
	default:
		return tfjson.SchemaDescriptionKindPlain
	case tfprotov6.StringKindMarkdown:
		return tfjson.SchemaDescriptionKindMarkdown
	}
}

func schemaNestedBlock(b *tfprotov6.SchemaNestedBlock) *tfjson.SchemaBlockType {
	var nesting tfjson.SchemaNestingMode
	switch b.Nesting {
	case tfprotov6.SchemaNestedBlockNestingModeSingle:
		nesting = tfjson.SchemaNestingModeSingle
	case tfprotov6.SchemaNestedBlockNestingModeGroup:
		nesting = tfjson.SchemaNestingModeGroup
	case tfprotov6.SchemaNestedBlockNestingModeList:
		nesting = tfjson.SchemaNestingModeList
	case tfprotov6.SchemaNestedBlockNestingModeMap:
		nesting = tfjson.SchemaNestingModeMap
	case tfprotov6.SchemaNestedBlockNestingModeSet:
		nesting = tfjson.SchemaNestingModeSet
	default:
		// In all other cases we'll leave it as the zero value (invalid) and
		// let the caller validate it and deal with this.
	}

	nb := &tfjson.SchemaBlockType{
		NestingMode: nesting,
		MinItems:    uint64(b.MinItems),
		MaxItems:    uint64(b.MaxItems),
	}

	nested := ProtoToConfigSchema(b.Block)
	nb.Block = nested
	return nb
}

func protoObjectToConfigSchema(b *tfprotov6.SchemaObject) *tfjson.SchemaNestedAttributeType {
	var nesting tfjson.SchemaNestingMode
	switch b.Nesting {
	case tfprotov6.SchemaObjectNestingModeSingle:
		nesting = tfjson.SchemaNestingModeSingle
	case tfprotov6.SchemaObjectNestingModeList:
		nesting = tfjson.SchemaNestingModeList
	case tfprotov6.SchemaObjectNestingModeMap:
		nesting = tfjson.SchemaNestingModeMap
	case tfprotov6.SchemaObjectNestingModeSet:
		nesting = tfjson.SchemaNestingModeSet
	default:
		// In all other cases we'll leave it as the zero value (invalid) and
		// let the caller validate it and deal with this.
	}

	object := &tfjson.SchemaNestedAttributeType{
		Attributes:  make(map[string]*tfjson.SchemaAttribute),
		NestingMode: nesting,
	}

	for _, a := range b.Attributes {
		attr := &tfjson.SchemaAttribute{
			Description:     a.Description,
			DescriptionKind: schemaStringKind(a.DescriptionKind),
			Required:        a.Required,
			Optional:        a.Optional,
			Computed:        a.Computed,
			Sensitive:       a.Sensitive,
			Deprecated:      a.Deprecated,
		}

		if a.Type != nil {
			b, err := a.Type.MarshalJSON()
			if err != nil {
				panic(err)
			}
			if err := json.Unmarshal(b, &attr.AttributeType); err != nil {
				panic(err)
			}
		}

		if a.NestedType != nil {
			attr.AttributeNestedType = protoObjectToConfigSchema(a.NestedType)
		}

		object.Attributes[a.Name] = attr
	}

	return object
}

// sortedKeys returns the lexically sorted keys from the given map. This is
// used to make schema conversions are deterministic. This panics if map keys
// are not a string.
func sortedKeys(m interface{}) []string {
	v := reflect.ValueOf(m)
	keys := make([]string, v.Len())

	mapKeys := v.MapKeys()
	for i, k := range mapKeys {
		keys[i] = k.Interface().(string)
	}

	sort.Strings(keys)
	return keys
}

func tfjsonObjectToProto(b *tfjson.SchemaNestedAttributeType) *tfprotov6.SchemaObject {
	var nesting tfprotov6.SchemaObjectNestingMode
	switch b.NestingMode {
	case tfjson.SchemaNestingModeSingle:
		nesting = tfprotov6.SchemaObjectNestingModeSingle
	case tfjson.SchemaNestingModeList:
		nesting = tfprotov6.SchemaObjectNestingModeList
	case tfjson.SchemaNestingModeSet:
		nesting = tfprotov6.SchemaObjectNestingModeSet
	case tfjson.SchemaNestingModeMap:
		nesting = tfprotov6.SchemaObjectNestingModeMap
	default:
		nesting = tfprotov6.SchemaObjectNestingModeInvalid
	}

	attributes := make([]*tfprotov6.SchemaAttribute, 0, len(b.Attributes))

	for _, name := range sortedKeys(b.Attributes) {
		a := b.Attributes[name]

		attr := &tfprotov6.SchemaAttribute{
			Name:            name,
			Description:     a.Description,
			DescriptionKind: protoStringKind(a.DescriptionKind),
			Optional:        a.Optional,
			Computed:        a.Computed,
			Required:        a.Required,
			Sensitive:       a.Sensitive,
			Deprecated:      a.Deprecated,
		}

		if a.AttributeType != cty.NilType {
			ty, err := json.Marshal(a.AttributeType)
			if err != nil {
				panic(err)
			}
			tftype, err := tftypes.ParseJSONType(ty)
			if err != nil {
				panic(err)
			}
			attr.Type = tftype
		}

		if a.AttributeNestedType != nil {
			attr.NestedType = tfjsonObjectToProto(a.AttributeNestedType)
		}

		attributes = append(attributes, attr)
	}

	return &tfprotov6.SchemaObject{
		Attributes: attributes,
		Nesting:    nesting,
	}
}
