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

func ConfigureProviderClientCapabilities(in *tfprotov6.ConfigureProviderClientCapabilities) *tfplugin6.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func ReadDataSourceClientCapabilities(in *tfprotov6.ReadDataSourceClientCapabilities) *tfplugin6.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func ReadResourceClientCapabilities(in *tfprotov6.ReadResourceClientCapabilities) *tfplugin6.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func PlanResourceChangeClientCapabilities(in *tfprotov6.PlanResourceChangeClientCapabilities) *tfplugin6.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func ImportResourceStateClientCapabilities(in *tfprotov6.ImportResourceStateClientCapabilities) *tfplugin6.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func PlanActionClientCapabilities(in *tfprotov6.PlanActionClientCapabilities) *tfplugin6.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func InvokeActionClientCapabilities(in *tfprotov6.InvokeActionClientCapabilities) *tfplugin6.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ClientCapabilities{}

	return resp
}
