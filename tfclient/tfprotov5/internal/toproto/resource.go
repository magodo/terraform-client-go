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
		ClientCapabilities: ValidateResourceTypeConfigClientCapabilities(in.ClientCapabilities),
		TypeName:           in.TypeName,
		Config:             DynamicValue(in.Config),
	}

	return req
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

func UpgradeResourceIdentity_Request(in *tfprotov5.UpgradeResourceIdentityRequest) *tfplugin5.UpgradeResourceIdentity_Request {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.UpgradeResourceIdentity_Request{
		RawIdentity: RawState(in.RawIdentity),
		TypeName:    in.TypeName,
		Version:     in.Version,
	}

	return resp
}

func ReadResource_Request(in *tfprotov5.ReadResourceRequest) *tfplugin5.ReadResource_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.ReadResource_Request{
		TypeName:           in.TypeName,
		CurrentState:       DynamicValue(in.CurrentState),
		Private:            in.Private,
		ProviderMeta:       DynamicValue(in.ProviderMeta),
		CurrentIdentity:    ResourceIdentityData(in.CurrentIdentity),
		ClientCapabilities: ReadResourceClientCapabilities(in.ClientCapabilities),
	}

	return req
}

func PlanResourceChange_Request(in *tfprotov5.PlanResourceChangeRequest) *tfplugin5.PlanResourceChange_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.PlanResourceChange_Request{
		TypeName:           in.TypeName,
		PriorState:         DynamicValue(in.PriorState),
		ProposedNewState:   DynamicValue(in.ProposedNewState),
		Config:             DynamicValue(in.Config),
		PriorPrivate:       in.PriorPrivate,
		ProviderMeta:       DynamicValue(in.ProviderMeta),
		PriorIdentity:      ResourceIdentityData(in.PriorIdentity),
		ClientCapabilities: PlanResourceChangeClientCapabilities(in.ClientCapabilities),
	}

	return req
}

func ApplyResourceChange_Request(in *tfprotov5.ApplyResourceChangeRequest) *tfplugin5.ApplyResourceChange_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.ApplyResourceChange_Request{
		TypeName:        in.TypeName,
		PriorState:      DynamicValue(in.PriorState),
		PlannedState:    DynamicValue(in.PlannedState),
		Config:          DynamicValue(in.Config),
		PlannedPrivate:  in.PlannedPrivate,
		ProviderMeta:    DynamicValue(in.ProviderMeta),
		PlannedIdentity: ResourceIdentityData(in.PlannedIdentity),
	}

	return req
}

func ImportResourceState_Request(in *tfprotov5.ImportResourceStateRequest) *tfplugin5.ImportResourceState_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.ImportResourceState_Request{
		TypeName:           in.TypeName,
		Id:                 in.ID,
		Identity:           ResourceIdentityData(in.Identity),
		ClientCapabilities: ImportResourceStateClientCapabilities(in.ClientCapabilities),
	}

	return req
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
		SourceIdentity:        RawState(in.SourceIdentity),
	}

	return req
}
