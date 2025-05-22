package fromproto

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func ResourceMetadata(in *tfplugin5.GetMetadata_ResourceMetadata) *tfprotov5.ResourceMetadata {
	if in == nil {
		return nil
	}

	return &tfprotov5.ResourceMetadata{
		TypeName: in.TypeName,
	}
}

func ValidateResourceTypeConfigResponse(in *tfplugin5.ValidateResourceTypeConfig_Response) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov5.ValidateResourceTypeConfigResponse{
		Diagnostics: diags,
	}

	return resp, nil
}

func UpgradeResourceStateResponse(in *tfplugin5.UpgradeResourceState_Response) (*tfprotov5.UpgradeResourceStateResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov5.UpgradeResourceStateResponse{
		UpgradedState: DynamicValue(in.UpgradedState),
		Diagnostics:   diags,
	}

	return resp, nil
}

func UpgradeResourceIdentityResponse(in *tfplugin5.UpgradeResourceIdentity_Response) (*tfprotov5.UpgradeResourceIdentityResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov5.UpgradeResourceIdentityResponse{
		Diagnostics:      diags,
		UpgradedIdentity: ResourceIdentityData(in.UpgradedIdentity),
	}

	return resp, nil
}

func ReadResourceResponse(in *tfplugin5.ReadResource_Response) (*tfprotov5.ReadResourceResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov5.ReadResourceResponse{
		NewState:    DynamicValue(in.NewState),
		Diagnostics: diags,
		Private:     in.Private,
		Deferred:    Deferred(in.Deferred),
		NewIdentity: ResourceIdentityData(in.NewIdentity),
	}

	return resp, nil
}

func PlanResourceChangeResponse(in *tfplugin5.PlanResourceChange_Response) (*tfprotov5.PlanResourceChangeResponse, error) {
	if in == nil {
		return nil, nil
	}

	requiresReplace, err := AttributePaths(in.RequiresReplace)
	if err != nil {
		return nil, err
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov5.PlanResourceChangeResponse{
		PlannedState:                DynamicValue(in.PlannedState),
		RequiresReplace:             requiresReplace,
		PlannedPrivate:              in.PlannedPrivate,
		Diagnostics:                 diags,
		UnsafeToUseLegacyTypeSystem: in.LegacyTypeSystem,
		Deferred:                    Deferred(in.Deferred),
		PlannedIdentity:             ResourceIdentityData(in.PlannedIdentity),
	}

	return resp, nil
}

func ApplyResourceChangeResponse(in *tfplugin5.ApplyResourceChange_Response) (*tfprotov5.ApplyResourceChangeResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov5.ApplyResourceChangeResponse{
		NewState:                    DynamicValue(in.NewState),
		Private:                     in.Private,
		Diagnostics:                 diags,
		UnsafeToUseLegacyTypeSystem: in.LegacyTypeSystem,
		NewIdentity:                 ResourceIdentityData(in.NewIdentity),
	}

	return resp, nil
}

func ImportResourceStateResponse(in *tfplugin5.ImportResourceState_Response) (*tfprotov5.ImportResourceStateResponse, error) {
	if in == nil {
		return nil, nil
	}

	importedResources, err := ImportedResources(in.ImportedResources)
	if err != nil {
		return nil, err
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov5.ImportResourceStateResponse{
		ImportedResources: importedResources,
		Diagnostics:       diags,
		Deferred:          Deferred(in.Deferred),
	}

	return resp, nil
}

func ImportedResource(in *tfplugin5.ImportResourceState_ImportedResource) (*tfprotov5.ImportedResource, error) {
	resp := &tfprotov5.ImportedResource{
		TypeName: in.TypeName,
		Private:  in.Private,
		Identity: ResourceIdentityData(in.Identity),
	}
	if in.State != nil {
		resp.State = DynamicValue(in.State)
	}
	return resp, nil
}

func ImportedResources(in []*tfplugin5.ImportResourceState_ImportedResource) ([]*tfprotov5.ImportedResource, error) {
	resp := make([]*tfprotov5.ImportedResource, 0, len(in))
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

func MoveResourceStateResponse(in *tfplugin5.MoveResourceState_Response) (*tfprotov5.MoveResourceStateResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov5.MoveResourceStateResponse{
		TargetPrivate:  in.TargetPrivate,
		TargetState:    DynamicValue(in.TargetState),
		Diagnostics:    diags,
		TargetIdentity: ResourceIdentityData(in.TargetIdentity),
	}

	return resp, nil
}
