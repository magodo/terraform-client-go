package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func GetMetadata_ResourceMetadata(in *tfprotov5.ResourceMetadata) *tfplugin5.GetMetadata_ResourceMetadata {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.GetMetadata_ResourceMetadata{
		TypeName: in.TypeName,
	}

	return resp
}

func ValidateResourceTypeConfig_Request(in *tfprotov5.ValidateResourceTypeConfigRequest) *tfplugin5.ValidateResourceTypeConfig_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.ValidateResourceTypeConfig_Request{
		TypeName: in.TypeName,
		Config:   DynamicValue(in.Config),
	}

	return req
}

func ValidateResourceTypeConfig_Response(in *tfprotov5.ValidateResourceTypeConfigResponse) *tfplugin5.ValidateResourceTypeConfig_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ValidateResourceTypeConfig_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}

	return resp
}

func UpgradeResourceState_Request(in *tfprotov5.UpgradeResourceStateRequest) *tfplugin5.UpgradeResourceState_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.UpgradeResourceState_Request{
		TypeName: in.TypeName,
		Version:  in.Version,
		RawState: RawState(in.RawState),
	}

	return req
}

func UpgradeResourceState_Response(in *tfprotov5.UpgradeResourceStateResponse) *tfplugin5.UpgradeResourceState_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.UpgradeResourceState_Response{
		Diagnostics:   Diagnostics(in.Diagnostics),
		UpgradedState: DynamicValue(in.UpgradedState),
	}

	return resp
}

func ReadResource_Request(in *tfprotov5.ReadResourceRequest) *tfplugin5.ReadResource_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.ReadResource_Request{
		TypeName:     in.TypeName,
		CurrentState: DynamicValue(in.CurrentState),
		Private:      in.Private,
		ProviderMeta: DynamicValue(in.ProviderMeta),
	}

	if in.ClientCapabilities != nil {
		req.ClientCapabilities = &tfplugin5.ClientCapabilities{
			DeferralAllowed: in.ClientCapabilities.DeferralAllowed,
		}
	}

	return req
}

func ReadResource_Response(in *tfprotov5.ReadResourceResponse) *tfplugin5.ReadResource_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ReadResource_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
		NewState:    DynamicValue(in.NewState),
		Private:     in.Private,
		Deferred:    Deferred(in.Deferred),
	}

	return resp
}

func PlanResourceChange_Request(in *tfprotov5.PlanResourceChangeRequest) *tfplugin5.PlanResourceChange_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.PlanResourceChange_Request{
		TypeName:         in.TypeName,
		PriorState:       DynamicValue(in.PriorState),
		ProposedNewState: DynamicValue(in.ProposedNewState),
		Config:           DynamicValue(in.Config),
		PriorPrivate:     in.PriorPrivate,
		ProviderMeta:     DynamicValue(in.ProviderMeta),
	}

	if in.ClientCapabilities != nil {
		req.ClientCapabilities = &tfplugin5.ClientCapabilities{
			DeferralAllowed: in.ClientCapabilities.DeferralAllowed,
		}
	}

	return req
}

func PlanResourceChange_Response(in *tfprotov5.PlanResourceChangeResponse) *tfplugin5.PlanResourceChange_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.PlanResourceChange_Response{
		Diagnostics:      Diagnostics(in.Diagnostics),
		LegacyTypeSystem: in.UnsafeToUseLegacyTypeSystem, //nolint:staticcheck
		PlannedPrivate:   in.PlannedPrivate,
		PlannedState:     DynamicValue(in.PlannedState),
		RequiresReplace:  AttributePaths(in.RequiresReplace),
		Deferred:         Deferred(in.Deferred),
	}

	return resp
}

func ApplyResourceChange_Request(in *tfprotov5.ApplyResourceChangeRequest) *tfplugin5.ApplyResourceChange_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.ApplyResourceChange_Request{
		TypeName:       in.TypeName,
		PriorState:     DynamicValue(in.PriorState),
		PlannedState:   DynamicValue(in.PlannedState),
		Config:         DynamicValue(in.Config),
		PlannedPrivate: in.PlannedPrivate,
		ProviderMeta:   DynamicValue(in.ProviderMeta),
	}

	return req
}

func ApplyResourceChange_Response(in *tfprotov5.ApplyResourceChangeResponse) *tfplugin5.ApplyResourceChange_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ApplyResourceChange_Response{
		Diagnostics:      Diagnostics(in.Diagnostics),
		LegacyTypeSystem: in.UnsafeToUseLegacyTypeSystem, //nolint:staticcheck
		NewState:         DynamicValue(in.NewState),
		Private:          in.Private,
	}

	return resp
}

func ImportResourceState_Request(in *tfprotov5.ImportResourceStateRequest) *tfplugin5.ImportResourceState_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.ImportResourceState_Request{
		TypeName: in.TypeName,
		Id:       in.ID,
	}

	if in.ClientCapabilities != nil {
		req.ClientCapabilities = &tfplugin5.ClientCapabilities{
			DeferralAllowed: in.ClientCapabilities.DeferralAllowed,
		}
	}

	return req
}

func ImportResourceState_Response(in *tfprotov5.ImportResourceStateResponse) *tfplugin5.ImportResourceState_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ImportResourceState_Response{
		Diagnostics:       Diagnostics(in.Diagnostics),
		ImportedResources: ImportResourceState_ImportedResources(in.ImportedResources),
		Deferred:          Deferred(in.Deferred),
	}

	return resp
}

func ImportResourceState_ImportedResource(in *tfprotov5.ImportedResource) *tfplugin5.ImportResourceState_ImportedResource {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ImportResourceState_ImportedResource{
		Private:  in.Private,
		State:    DynamicValue(in.State),
		TypeName: in.TypeName,
	}

	return resp
}

func ImportResourceState_ImportedResources(in []*tfprotov5.ImportedResource) []*tfplugin5.ImportResourceState_ImportedResource {
	resp := make([]*tfplugin5.ImportResourceState_ImportedResource, 0, len(in))

	for _, i := range in {
		resp = append(resp, ImportResourceState_ImportedResource(i))
	}

	return resp
}

func MoveResourceState_Request(in *tfprotov5.MoveResourceStateRequest) *tfplugin5.MoveResourceState_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.MoveResourceState_Request{
		SourceProviderAddress: in.SourceProviderAddress,
		SourceTypeName:        in.SourceTypeName,
		SourceSchemaVersion:   in.SourceSchemaVersion,
		SourceState:           RawState(in.SourceState),
		TargetTypeName:        in.TargetTypeName,
		SourcePrivate:         in.SourcePrivate,
	}

	return req
}

func MoveResourceState_Response(in *tfprotov5.MoveResourceStateResponse) *tfplugin5.MoveResourceState_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.MoveResourceState_Response{
		Diagnostics:   Diagnostics(in.Diagnostics),
		TargetPrivate: in.TargetPrivate,
		TargetState:   DynamicValue(in.TargetState),
	}

	return resp
}
