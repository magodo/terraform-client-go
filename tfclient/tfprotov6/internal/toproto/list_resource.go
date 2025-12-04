package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func ListResourceRequest(in *tfprotov6.ListResourceRequest) *tfplugin6.ListResource_Request {
	if in == nil {
		return nil
	}

	return &tfplugin6.ListResource_Request{
		TypeName:              in.TypeName,
		Config:                DynamicValue(in.Config),
		IncludeResourceObject: in.IncludeResource,
		Limit:                 in.Limit,
	}
}

func ValidateListResourceConfigRequest(in *tfprotov6.ValidateListResourceConfigRequest) *tfplugin6.ValidateListResourceConfig_Request {
	if in == nil {
		return nil
	}

	return &tfplugin6.ValidateListResourceConfig_Request{
		TypeName:              in.TypeName,
		Config:                DynamicValue(in.Config),
		IncludeResourceObject: DynamicValue(in.IncludeResourceObject),
		Limit:                 DynamicValue(in.Limit),
	}
}
