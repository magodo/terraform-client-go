package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func ValidateEphemeralResourceConfigRequest(in *tfprotov6.ValidateEphemeralResourceConfigRequest) *tfplugin6.ValidateEphemeralResourceConfig_Request {
	if in == nil {
		return nil
	}

	return &tfplugin6.ValidateEphemeralResourceConfig_Request{
		TypeName: in.TypeName,
		Config:   DynamicValue(in.Config),
	}
}

func OpenEphemeralResourceRequest(in *tfprotov6.OpenEphemeralResourceRequest) *tfplugin6.OpenEphemeralResource_Request {
	if in == nil {
		return nil
	}

	return &tfplugin6.OpenEphemeralResource_Request{
		TypeName:           in.TypeName,
		Config:             DynamicValue(in.Config),
		ClientCapabilities: OpenEphemeralResourceClientCapabilities(in.ClientCapabilities),
	}
}

func RenewEphemeralResourceRequest(in *tfprotov6.RenewEphemeralResourceRequest) *tfplugin6.RenewEphemeralResource_Request {
	if in == nil {
		return nil
	}

	return &tfplugin6.RenewEphemeralResource_Request{
		TypeName: in.TypeName,
		Private:  in.Private,
	}
}

func CloseEphemeralResourceRequest(in *tfprotov6.CloseEphemeralResourceRequest) *tfplugin6.CloseEphemeralResource_Request {
	if in == nil {
		return nil
	}

	return &tfplugin6.CloseEphemeralResource_Request{
		TypeName: in.TypeName,
		Private:  in.Private,
	}
}
