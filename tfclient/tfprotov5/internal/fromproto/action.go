package fromproto

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func ActionMetadata(in *tfplugin5.GetMetadata_ActionMetadata) *tfprotov5.ActionMetadata {
	if in == nil {
		return nil
	}

	return &tfprotov5.ActionMetadata{
		TypeName: in.TypeName,
	}
}

func ValidateActionConfig_Response(in *tfplugin5.ValidateActionConfig_Response) (*tfprotov5.ValidateActionConfigResponse, error) {
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov5.ValidateActionConfigResponse{
		Diagnostics: diags,
	}

	return resp, nil
}

func PlanAction_Response(in *tfplugin5.PlanAction_Response) (*tfprotov5.PlanActionResponse, error) {
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov5.PlanActionResponse{
		Diagnostics: diags,
		Deferred:    Deferred(in.Deferred),
	}

	return resp, nil
}

func InvokeAction_InvokeActionEvent(in *tfplugin5.InvokeAction_Event) (*tfprotov5.InvokeActionEvent, error) {
	switch event := (in.Type).(type) {
	case *tfplugin5.InvokeAction_Event_Progress_:
		return &tfprotov5.InvokeActionEvent{
			Type: &tfprotov5.ProgressInvokeActionEventType{
				Message: event.Progress.GetMessage(),
			},
		}, nil
	case *tfplugin5.InvokeAction_Event_Completed_:
		diags, err := Diagnostics(event.Completed.GetDiagnostics())
		if err != nil {
			return nil, err
		}
		return &tfprotov5.InvokeActionEvent{
			Type: &tfprotov5.CompletedInvokeActionEventType{
				Diagnostics: diags,
			},
		}, nil
	}

	// It is not currently possible to create tfprotov5.InvokeActionEventType
	// implementations outside the tfprotov5 package. If this panic was reached,
	// it implies that a new event type was introduced and needs to be implemented
	// as a new case above.
	panic(fmt.Sprintf("unimplemented tfprotov5.InvokeActionEventType type: %T", in.Type))
}
