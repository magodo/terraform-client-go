package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func Deferred(in *tfplugin6.Deferred) *tfprotov6.Deferred {
	if in == nil {
		return nil
	}
	return &tfprotov6.Deferred{
		Reason: tfprotov6.DeferredReason(in.Reason),
	}
}
