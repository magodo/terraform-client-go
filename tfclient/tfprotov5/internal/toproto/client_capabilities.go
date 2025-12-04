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

func ConfigureProviderClientCapabilities(in *tfprotov5.ConfigureProviderClientCapabilities) *tfplugin5.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func ReadDataSourceClientCapabilities(in *tfprotov5.ReadDataSourceClientCapabilities) *tfplugin5.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func ReadResourceClientCapabilities(in *tfprotov5.ReadResourceClientCapabilities) *tfplugin5.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func PlanResourceChangeClientCapabilities(in *tfprotov5.PlanResourceChangeClientCapabilities) *tfplugin5.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func ImportResourceStateClientCapabilities(in *tfprotov5.ImportResourceStateClientCapabilities) *tfplugin5.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
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

func PlanActionClientCapabilities(in *tfprotov5.PlanActionClientCapabilities) *tfplugin5.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func InvokeActionClientCapabilities(in *tfprotov5.InvokeActionClientCapabilities) *tfplugin5.ClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ClientCapabilities{}

	return resp
}
