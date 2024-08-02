// Duplicating the logic from: github.com/hashicorp/terraform/internal/plugin/convert/deferred.go

package convert

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/typ"
)

// ProtoToDeferred translates a proto.Deferred to a providers.Deferred.
func ProtoToDeferred(d *tfprotov5.Deferred) *typ.Deferred {
	if d == nil {
		return nil
	}

	var reason typ.DeferredReason
	switch d.Reason {
	case tfprotov5.DeferredReasonUnknown:
		reason = typ.DeferredReasonInvalid
	case tfprotov5.DeferredReasonResourceConfigUnknown:
		reason = typ.DeferredReasonResourceConfigUnknown
	case tfprotov5.DeferredReasonProviderConfigUnknown:
		reason = typ.DeferredReasonProviderConfigUnknown
	case tfprotov5.DeferredReasonAbsentPrereq:
		reason = typ.DeferredReasonAbsentPrereq
	default:
		reason = typ.DeferredReasonInvalid
	}

	return &typ.Deferred{
		Reason: reason,
	}
}
