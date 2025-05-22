package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func ValidateEphemeralResourceConfigRequest(in *tfprotov5.ValidateEphemeralResourceConfigRequest) *tfplugin5.ValidateEphemeralResourceConfig_Request {
	if in == nil {
		return nil
	}

	return &tfplugin5.ValidateEphemeralResourceConfig_Request{
		TypeName: in.TypeName,
		Config:   DynamicValue(in.Config),
	}
}

func OpenEphemeralResourceRequest(in *tfprotov5.OpenEphemeralResourceRequest) *tfplugin5.OpenEphemeralResource_Request {
	if in == nil {
		return nil
	}

	return &tfplugin5.OpenEphemeralResource_Request{
		TypeName:           in.TypeName,
		Config:             DynamicValue(in.Config),
		ClientCapabilities: OpenEphemeralResourceClientCapabilities(in.ClientCapabilities),
	}
}

func RenewEphemeralResourceRequest(in *tfprotov5.RenewEphemeralResourceRequest) *tfplugin5.RenewEphemeralResource_Request {
	if in == nil {
		return nil
	}

	return &tfplugin5.RenewEphemeralResource_Request{
		TypeName: in.TypeName,
		Private:  in.Private,
	}
}

func CloseEphemeralResourceRequest(in *tfprotov5.CloseEphemeralResourceRequest) *tfplugin5.CloseEphemeralResource_Request {
	if in == nil {
		return nil
	}

	return &tfplugin5.CloseEphemeralResource_Request{
		TypeName: in.TypeName,
		Private:  in.Private,
	}
}
