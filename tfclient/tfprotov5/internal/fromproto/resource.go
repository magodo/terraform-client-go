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

func ValidateResourceTypeConfigRequest(in *tfplugin5.ValidateResourceTypeConfig_Request) *tfprotov5.ValidateResourceTypeConfigRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.ValidateResourceTypeConfigRequest{
		Config:   DynamicValue(in.Config),
		TypeName: in.TypeName,
	}

	return resp
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

func UpgradeResourceStateRequest(in *tfplugin5.UpgradeResourceState_Request) *tfprotov5.UpgradeResourceStateRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.UpgradeResourceStateRequest{
		RawState: RawState(in.RawState),
		TypeName: in.TypeName,
		Version:  in.Version,
	}

	return resp
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

func ReadResourceRequest(in *tfplugin5.ReadResource_Request) *tfprotov5.ReadResourceRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.ReadResourceRequest{
		CurrentState:       DynamicValue(in.CurrentState),
		Private:            in.Private,
		ProviderMeta:       DynamicValue(in.ProviderMeta),
		TypeName:           in.TypeName,
		ClientCapabilities: ReadResourceClientCapabilities(in.ClientCapabilities),
	}

	return resp
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
	}

	return resp, nil
}

func PlanResourceChangeRequest(in *tfplugin5.PlanResourceChange_Request) *tfprotov5.PlanResourceChangeRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.PlanResourceChangeRequest{
		Config:             DynamicValue(in.Config),
		PriorPrivate:       in.PriorPrivate,
		PriorState:         DynamicValue(in.PriorState),
		ProposedNewState:   DynamicValue(in.ProposedNewState),
		ProviderMeta:       DynamicValue(in.ProviderMeta),
		TypeName:           in.TypeName,
		ClientCapabilities: PlanResourceChangeClientCapabilities(in.ClientCapabilities),
	}

	return resp
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
	}

	return resp, nil
}

func ApplyResourceChangeRequest(in *tfplugin5.ApplyResourceChange_Request) *tfprotov5.ApplyResourceChangeRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.ApplyResourceChangeRequest{
		Config:         DynamicValue(in.Config),
		PlannedPrivate: in.PlannedPrivate,
		PlannedState:   DynamicValue(in.PlannedState),
		PriorState:     DynamicValue(in.PriorState),
		ProviderMeta:   DynamicValue(in.ProviderMeta),
		TypeName:       in.TypeName,
	}

	return resp
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
	}

	return resp, nil
}

func ImportResourceStateRequest(in *tfplugin5.ImportResourceState_Request) *tfprotov5.ImportResourceStateRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.ImportResourceStateRequest{
		TypeName:           in.TypeName,
		ID:                 in.Id,
		ClientCapabilities: ImportResourceStateClientCapabilities(in.ClientCapabilities),
	}

	return resp
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

func MoveResourceStateRequest(in *tfplugin5.MoveResourceState_Request) *tfprotov5.MoveResourceStateRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.MoveResourceStateRequest{
		SourcePrivate:         in.SourcePrivate,
		SourceProviderAddress: in.SourceProviderAddress,
		SourceSchemaVersion:   in.SourceSchemaVersion,
		SourceState:           RawState(in.SourceState),
		SourceTypeName:        in.SourceTypeName,
		TargetTypeName:        in.TargetTypeName,
	}

	return resp
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
		TargetPrivate: in.TargetPrivate,
		TargetState:   DynamicValue(in.TargetState),
		Diagnostics:   diags,
	}

	return resp, nil
}
