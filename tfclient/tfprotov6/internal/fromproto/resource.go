package fromproto

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func ResourceMetadata(in *tfplugin6.GetMetadata_ResourceMetadata) *tfprotov6.ResourceMetadata {
	if in == nil {
		return nil
	}

	return &tfprotov6.ResourceMetadata{
		TypeName: in.TypeName,
	}
}

func ValidateResourceConfigResponse(in *tfplugin6.ValidateResourceConfig_Response) (*tfprotov6.ValidateResourceConfigResponse, error) {
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}
	return &tfprotov6.ValidateResourceConfigResponse{
		Diagnostics: diags,
	}, nil
}

func UpgradeResourceStateResponse(in *tfplugin6.UpgradeResourceState_Response) (*tfprotov6.UpgradeResourceStateResponse, error) {
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}
	resp := &tfprotov6.UpgradeResourceStateResponse{
		Diagnostics:   diags,
		UpgradedState: DynamicValue(in.UpgradedState),
	}
	return resp, nil
}

func UpgradeResourceIdentityResponse(in *tfplugin6.UpgradeResourceIdentity_Response) (*tfprotov6.UpgradeResourceIdentityResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov6.UpgradeResourceIdentityResponse{
		Diagnostics:      diags,
		UpgradedIdentity: ResourceIdentityData(in.UpgradedIdentity),
	}

	return resp, nil
}

func ReadResourceResponse(in *tfplugin6.ReadResource_Response) (*tfprotov6.ReadResourceResponse, error) {
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}
	resp := &tfprotov6.ReadResourceResponse{
		Private:     in.Private,
		NewState:    DynamicValue(in.NewState),
		Diagnostics: diags,
		NewIdentity: ResourceIdentityData(in.NewIdentity),
	}
	return resp, nil
}

func PlanResourceChangeResponse(in *tfplugin6.PlanResourceChange_Response) (*tfprotov6.PlanResourceChangeResponse, error) {
	requireReplace, err := AttributePaths(in.RequiresReplace)
	if err != nil {
		return nil, err
	}
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}
	resp := &tfprotov6.PlanResourceChangeResponse{
		RequiresReplace:             requireReplace,
		Diagnostics:                 diags,
		PlannedPrivate:              in.PlannedPrivate,
		UnsafeToUseLegacyTypeSystem: in.LegacyTypeSystem,
		PlannedState:                DynamicValue(in.PlannedState),
		Deferred:                    Deferred(in.Deferred),
		PlannedIdentity:             ResourceIdentityData(in.PlannedIdentity),
	}
	return resp, nil
}

func ApplyResourceChangeResponse(in *tfplugin6.ApplyResourceChange_Response) (*tfprotov6.ApplyResourceChangeResponse, error) {
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}
	resp := &tfprotov6.ApplyResourceChangeResponse{
		Private:                     in.Private,
		UnsafeToUseLegacyTypeSystem: in.LegacyTypeSystem,
		Diagnostics:                 diags,
		NewState:                    DynamicValue(in.NewState),
		NewIdentity:                 ResourceIdentityData(in.NewIdentity),
	}
	return resp, nil
}

func ImportResourceStateResponse(in *tfplugin6.ImportResourceState_Response) (*tfprotov6.ImportResourceStateResponse, error) {
	imported, err := ImportedResources(in.ImportedResources)
	if err != nil {
		return nil, err
	}
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}
	return &tfprotov6.ImportResourceStateResponse{
		ImportedResources: imported,
		Diagnostics:       diags,
		Deferred:          Deferred(in.Deferred),
	}, nil
}

func ImportedResource(in *tfplugin6.ImportResourceState_ImportedResource) (*tfprotov6.ImportedResource, error) {
	resp := &tfprotov6.ImportedResource{
		TypeName: in.TypeName,
		Private:  in.Private,
		Identity: ResourceIdentityData(in.Identity),
	}
	if in.State != nil {
		resp.State = DynamicValue(in.State)
	}
	return resp, nil
}

func ImportedResources(in []*tfplugin6.ImportResourceState_ImportedResource) ([]*tfprotov6.ImportedResource, error) {
	resp := make([]*tfprotov6.ImportedResource, 0, len(in))
	for pos, i := range in {
		if i == nil {
			resp = append(resp, nil)
			continue
		}
		r, err := ImportedResource(i)
		if err != nil {
			return resp, fmt.Errorf("Error converting imported resource %d/%d: %w", pos+1, len(in), err)
		}
		resp = append(resp, r)
	}
	return resp, nil
}

func MoveResourceStateResponse(in *tfplugin6.MoveResourceState_Response) (*tfprotov6.MoveResourceStateResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov6.MoveResourceStateResponse{
		TargetPrivate:  in.TargetPrivate,
		TargetState:    DynamicValue(in.TargetState),
		Diagnostics:    diags,
		TargetIdentity: ResourceIdentityData(in.TargetIdentity),
	}

	return resp, nil
}
