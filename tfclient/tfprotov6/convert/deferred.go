// Duplicating the logic from: github.com/hashicorp/terraform/internal/plugin6/convert/deferred.go

package convert

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/typ"
)

// ProtoToDeferred translates a proto.Deferred to a providers.Deferred.
func ProtoToDeferred(d *tfprotov6.Deferred) *typ.Deferred {
	if d == nil {
		return nil
	}

	var reason typ.DeferredReason
	switch d.Reason {
	case tfprotov6.DeferredReasonUnknown:
		reason = typ.DeferredReasonInvalid
	case tfprotov6.DeferredReasonResourceConfigUnknown:
		reason = typ.DeferredReasonResourceConfigUnknown
	case tfprotov6.DeferredReasonProviderConfigUnknown:
		reason = typ.DeferredReasonProviderConfigUnknown
	case tfprotov6.DeferredReasonAbsentPrereq:
		reason = typ.DeferredReasonAbsentPrereq
	default:
		reason = typ.DeferredReasonInvalid
	}

	return &typ.Deferred{
		Reason: reason,
	}
}
