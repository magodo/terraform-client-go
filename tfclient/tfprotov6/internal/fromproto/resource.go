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

func ValidateResourceConfigRequest(in *tfplugin6.ValidateResourceConfig_Request) *tfprotov6.ValidateResourceConfigRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ValidateResourceConfigRequest{
		Config:   DynamicValue(in.Config),
		TypeName: in.TypeName,
	}

	return resp
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

func UpgradeResourceStateRequest(in *tfplugin6.UpgradeResourceState_Request) *tfprotov6.UpgradeResourceStateRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.UpgradeResourceStateRequest{
		RawState: RawState(in.RawState),
		TypeName: in.TypeName,
		Version:  in.Version,
	}

	return resp
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

func ReadResourceRequest(in *tfplugin6.ReadResource_Request) *tfprotov6.ReadResourceRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ReadResourceRequest{
		CurrentState:       DynamicValue(in.CurrentState),
		Private:            in.Private,
		ProviderMeta:       DynamicValue(in.ProviderMeta),
		TypeName:           in.TypeName,
		ClientCapabilities: ReadResourceClientCapabilities(in.ClientCapabilities),
	}

	return resp
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
	}
	return resp, nil
}

func PlanResourceChangeRequest(in *tfplugin6.PlanResourceChange_Request) *tfprotov6.PlanResourceChangeRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.PlanResourceChangeRequest{
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
	}
	return resp, nil
}

func ApplyResourceChangeRequest(in *tfplugin6.ApplyResourceChange_Request) *tfprotov6.ApplyResourceChangeRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ApplyResourceChangeRequest{
		Config:         DynamicValue(in.Config),
		PlannedPrivate: in.PlannedPrivate,
		PlannedState:   DynamicValue(in.PlannedState),
		PriorState:     DynamicValue(in.PriorState),
		ProviderMeta:   DynamicValue(in.ProviderMeta),
		TypeName:       in.TypeName,
	}

	return resp
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
	}
	return resp, nil
}

func ImportResourceStateRequest(in *tfplugin6.ImportResourceState_Request) *tfprotov6.ImportResourceStateRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ImportResourceStateRequest{
		TypeName:           in.TypeName,
		ID:                 in.Id,
		ClientCapabilities: ImportResourceStateClientCapabilities(in.ClientCapabilities),
	}

	return resp
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

func MoveResourceStateRequest(in *tfplugin6.MoveResourceState_Request) *tfprotov6.MoveResourceStateRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.MoveResourceStateRequest{
		SourcePrivate:         in.SourcePrivate,
		SourceProviderAddress: in.SourceProviderAddress,
		SourceSchemaVersion:   in.SourceSchemaVersion,
		SourceState:           RawState(in.SourceState),
		SourceTypeName:        in.SourceTypeName,
		TargetTypeName:        in.TargetTypeName,
	}

	return resp
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
		TargetPrivate: in.TargetPrivate,
		TargetState:   DynamicValue(in.TargetState),
		Diagnostics:   diags,
	}

	return resp, nil
}
