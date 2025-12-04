package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func ValidateActionConfigRequest(in *tfprotov6.ValidateActionConfigRequest) *tfplugin6.ValidateActionConfig_Request {
	if in == nil {
		return nil
	}

	return &tfplugin6.ValidateActionConfig_Request{
		ActionType: in.ActionType,
		Config:     DynamicValue(in.Config),
	}
}

func PlanActionRequest(in *tfprotov6.PlanActionRequest) *tfplugin6.PlanAction_Request {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.PlanAction_Request{
		ActionType:         in.ActionType,
		Config:             DynamicValue(in.Config),
		ClientCapabilities: PlanActionClientCapabilities(in.ClientCapabilities),
	}

	return resp
}

func InvokeActionRequest(in *tfprotov6.InvokeActionRequest) *tfplugin6.InvokeAction_Request {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.InvokeAction_Request{
		ActionType:         in.ActionType,
		Config:             DynamicValue(in.Config),
		ClientCapabilities: InvokeActionClientCapabilities(in.ClientCapabilities),
	}

	return resp
}
