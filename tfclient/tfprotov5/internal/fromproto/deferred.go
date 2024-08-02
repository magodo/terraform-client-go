package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func Deferred(in *tfplugin5.Deferred) *tfprotov5.Deferred {
	if in == nil {
		return nil
	}
	return &tfprotov5.Deferred{
		Reason: tfprotov5.DeferredReason(in.Reason),
	}
}
