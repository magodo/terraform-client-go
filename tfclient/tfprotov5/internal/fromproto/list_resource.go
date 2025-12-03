package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func ListResourceMetadata(in *tfplugin5.GetMetadata_ListResourceMetadata) *tfprotov5.ListResourceMetadata {
	return &tfprotov5.ListResourceMetadata{
		TypeName: in.GetTypeName(),
	}
}

func ListResource_ListResourceEvent(in *tfplugin5.ListResource_Event) (*tfprotov5.ListResourceResult, error) {
	diags, err := Diagnostics(in.GetDiagnostic())
	if err != nil {
		return nil, err
	}
	resp := &tfprotov5.ListResourceResult{
		DisplayName: in.DisplayName,
		Resource:    DynamicValue(in.GetResourceObject()),
		Identity:    ResourceIdentityData(in.Identity),
		Diagnostics: diags,
	}
	return resp, nil
}

func ValidateListResourceConfig_Response(in *tfplugin5.ValidateListResourceConfig_Response) (*tfprotov5.ValidateListResourceConfigResponse, error) {
	diags, err := Diagnostics(in.GetDiagnostics())
	if err != nil {
		return nil, err
	}
	return &tfprotov5.ValidateListResourceConfigResponse{
		Diagnostics: diags,
	}, nil
}
