package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func ValidateResourceConfigClientCapabilities(in *tfprotov6.ValidateResourceConfigClientCapabilities) *tfplugin6.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ClientCapabilities{
		WriteOnlyAttributesAllowed: in.WriteOnlyAttributesAllowed,
	}

	return resp
}

func OpenEphemeralResourceClientCapabilities(in *tfprotov6.OpenEphemeralResourceClientCapabilities) *tfplugin6.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}
