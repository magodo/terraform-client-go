package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func ValidateResourceTypeConfigClientCapabilities(in *tfprotov5.ValidateResourceTypeConfigClientCapabilities) *tfplugin5.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ClientCapabilities{
		WriteOnlyAttributesAllowed: in.WriteOnlyAttributesAllowed,
	}

	return resp
}

func OpenEphemeralResourceClientCapabilities(in *tfprotov5.OpenEphemeralResourceClientCapabilities) *tfplugin5.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}
