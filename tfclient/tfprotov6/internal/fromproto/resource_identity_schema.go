package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func ResourceIdentitySchema(in *tfplugin6.ResourceIdentitySchema) (*tfprotov6.ResourceIdentitySchema, error) {
	if in == nil {
		return nil, nil
	}

	identityAttrs, err := ResourceIdentitySchema_IdentityAttributes(in.IdentityAttributes)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov6.ResourceIdentitySchema{
		Version:            in.Version,
		IdentityAttributes: identityAttrs,
	}

	return resp, nil
}

func ResourceIdentitySchema_IdentityAttribute(in *tfplugin6.ResourceIdentitySchema_IdentityAttribute) (*tfprotov6.ResourceIdentitySchemaAttribute, error) {
	if in == nil {
		return nil, nil
	}

	typ, err := CtyType(in.Type)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov6.ResourceIdentitySchemaAttribute{
		Name:              in.Name,
		Type:              typ,
		RequiredForImport: in.RequiredForImport,
		OptionalForImport: in.OptionalForImport,
		Description:       in.Description,
	}

	return resp, nil
}

func ResourceIdentitySchema_IdentityAttributes(in []*tfplugin6.ResourceIdentitySchema_IdentityAttribute) ([]*tfprotov6.ResourceIdentitySchemaAttribute, error) {
	if in == nil {
		return nil, nil
	}

	resp := make([]*tfprotov6.ResourceIdentitySchemaAttribute, 0, len(in))

	for _, a := range in {
		attr, err := ResourceIdentitySchema_IdentityAttribute(a)
		if err != nil {
			return nil, err
		}
		resp = append(resp, attr)
	}

	return resp, nil
}
