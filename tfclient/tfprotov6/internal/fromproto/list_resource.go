package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func ListResourceMetadata(in *tfplugin6.GetMetadata_ListResourceMetadata) *tfprotov6.ListResourceMetadata {
	if in == nil {
		return nil
	}

	return &tfprotov6.ListResourceMetadata{
		TypeName: in.TypeName,
	}
}

func ListResource_ListResourceEvent(in *tfplugin6.ListResource_Event) (*tfprotov6.ListResourceResult, error) {
	diags, err := Diagnostics(in.GetDiagnostic())
	if err != nil {
		return nil, err
	}
	return &tfprotov6.ListResourceResult{
		DisplayName: in.DisplayName,
		Resource:    DynamicValue(in.ResourceObject),
		Identity:    ResourceIdentityData(in.Identity),
		Diagnostics: diags,
	}, nil
}

func ValidateListResourceConfig_Response(in *tfplugin6.ValidateListResourceConfig_Response) (*tfprotov6.ValidateListResourceConfigResponse, error) {
	diags, err := Diagnostics(in.GetDiagnostics())
	if err != nil {
		return nil, err
	}
	return &tfprotov6.ValidateListResourceConfigResponse{
		Diagnostics: diags,
	}, nil
}
