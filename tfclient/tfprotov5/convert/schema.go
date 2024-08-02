// Duplicating the logic from: github.com/hashicorp/terraform/internal/plugin/convert/schema.go

package convert

import (
	"encoding/json"
	"reflect"
	"sort"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func ConfigSchemaToProto(b *tfjson.SchemaBlock) *tfprotov5.SchemaBlock {
	block := &tfprotov5.SchemaBlock{
		Description:     b.Description,
		DescriptionKind: protoStringKind(b.DescriptionKind),
		Deprecated:      b.Deprecated,
	}

	for _, name := range sortedKeys(b.Attributes) {
		a := b.Attributes[name]

		attr := &tfprotov5.SchemaAttribute{
			Name:            name,
			Description:     a.Description,
			DescriptionKind: protoStringKind(a.DescriptionKind),
			Optional:        a.Optional,
			Computed:        a.Computed,
			Required:        a.Required,
			Sensitive:       a.Sensitive,
			Deprecated:      a.Deprecated,
		}

		ty, err := json.Marshal(a.AttributeType)
		if err != nil {
			panic(err)
		}

		tftype, err := tftypes.ParseJSONType(ty)
		if err != nil {
			panic(err)
		}

		attr.Type = tftype

		block.Attributes = append(block.Attributes, attr)
	}

	for _, name := range sortedKeys(b.NestedBlocks) {
		b := b.NestedBlocks[name]
		block.BlockTypes = append(block.BlockTypes, protoSchemaNestedBlock(name, b))
	}

	return block
}

func protoStringKind(k tfjson.SchemaDescriptionKind) tfprotov5.StringKind {
	switch k {
	default:
		return tfprotov5.StringKindPlain
	case tfjson.SchemaDescriptionKindMarkdown:
		return tfprotov5.StringKindMarkdown
	}
}

func protoSchemaNestedBlock(name string, b *tfjson.SchemaBlockType) *tfprotov5.SchemaNestedBlock {
	var nesting tfprotov5.SchemaNestedBlockNestingMode
	switch b.NestingMode {
	case tfjson.SchemaNestingModeSingle:
		nesting = tfprotov5.SchemaNestedBlockNestingModeSingle
	case tfjson.SchemaNestingModeGroup:
		nesting = tfprotov5.SchemaNestedBlockNestingModeGroup
	case tfjson.SchemaNestingModeList:
		nesting = tfprotov5.SchemaNestedBlockNestingModeList
	case tfjson.SchemaNestingModeSet:
		nesting = tfprotov5.SchemaNestedBlockNestingModeSet
	case tfjson.SchemaNestingModeMap:
		nesting = tfprotov5.SchemaNestedBlockNestingModeMap
	default:
		nesting = tfprotov5.SchemaNestedBlockNestingModeInvalid
	}
	return &tfprotov5.SchemaNestedBlock{
		TypeName: name,
		Block:    ConfigSchemaToProto(b.Block),
		Nesting:  nesting,
		MinItems: int64(b.MinItems),
		MaxItems: int64(b.MaxItems),
	}
}

func ProtoToProviderSchema(s *tfprotov5.Schema) tfjson.Schema {
	return tfjson.Schema{
		Version: uint64(s.Version),
		Block:   ProtoToConfigSchema(s.Block),
	}
}

func ProtoToConfigSchema(b *tfprotov5.SchemaBlock) *tfjson.SchemaBlock {
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

		b, err := a.Type.MarshalJSON()
		if err != nil {
			panic(err)
		}

		if err := json.Unmarshal(b, &attr.AttributeType); err != nil {
			panic(err)
		}

		block.Attributes[a.Name] = attr
	}

	for _, b := range b.BlockTypes {
		block.NestedBlocks[b.TypeName] = schemaNestedBlock(b)
	}

	return block
}

func schemaStringKind(k tfprotov5.StringKind) tfjson.SchemaDescriptionKind {
	switch k {
	default:
		return tfjson.SchemaDescriptionKindPlain
	case tfprotov5.StringKindMarkdown:
		return tfjson.SchemaDescriptionKindMarkdown
	}
}

func schemaNestedBlock(b *tfprotov5.SchemaNestedBlock) *tfjson.SchemaBlockType {
	var nesting tfjson.SchemaNestingMode
	switch b.Nesting {
	case tfprotov5.SchemaNestedBlockNestingModeSingle:
		nesting = tfjson.SchemaNestingModeSingle
	case tfprotov5.SchemaNestedBlockNestingModeGroup:
		nesting = tfjson.SchemaNestingModeGroup
	case tfprotov5.SchemaNestedBlockNestingModeList:
		nesting = tfjson.SchemaNestingModeList
	case tfprotov5.SchemaNestedBlockNestingModeMap:
		nesting = tfjson.SchemaNestingModeMap
	case tfprotov5.SchemaNestedBlockNestingModeSet:
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
