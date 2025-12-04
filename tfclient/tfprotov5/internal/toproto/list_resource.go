package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func ListResourceRequest(in *tfprotov5.ListResourceRequest) *tfplugin5.ListResource_Request {
	return &tfplugin5.ListResource_Request{
		TypeName:              in.TypeName,
		Config:                DynamicValue(in.Config),
		IncludeResourceObject: in.IncludeResource,
		Limit:                 in.Limit,
	}
}

func ValidateListResourceConfigRequest(in *tfprotov5.ValidateListResourceConfigRequest) *tfplugin5.ValidateListResourceConfig_Request {
	return &tfplugin5.ValidateListResourceConfig_Request{
		TypeName:              in.TypeName,
		Config:                DynamicValue(in.Config),
		IncludeResourceObject: DynamicValue(in.IncludeResourceObject),
		Limit:                 DynamicValue(in.Limit),
	}
}
