// Logic is derived from: github.com/apparentlymart/terraform-provider/tfprovider/internal/protocol5/diagnostics.go (ba059c82d80b5c662af2b0e9de2416903ff05e44)

package convert

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/magodo/terraform-client-go/tfclient/typ"
	"github.com/zclconf/go-cty/cty"
)

func DecodeDiagnostics(raws []*tfprotov5.Diagnostic) typ.Diagnostics {
	if len(raws) == 0 {
		return nil
	}
	diags := make(typ.Diagnostics, 0, len(raws))
	for _, raw := range raws {
		diag := typ.Diagnostic{
			Summary:   raw.Summary,
			Detail:    raw.Detail,
			Attribute: DecodeAttributePath(raw.Attribute),
		}

		switch raw.Severity {
		case tfprotov5.DiagnosticSeverityError:
			diag.Severity = typ.Error
		case tfprotov5.DiagnosticSeverityWarning:
			diag.Severity = typ.Warning
		}

		diags = append(diags, diag)
	}
	return diags
}

func DecodeAttributePath(raws *tftypes.AttributePath) cty.Path {
	if raws == nil || len(raws.Steps()) == 0 {
		return nil
	}
	ret := make(cty.Path, 0, len(raws.Steps()))
	for _, raw := range raws.Steps() {
		switch s := raw.(type) {
		case *tftypes.AttributeName:
			ret = ret.GetAttr(string(*s))
		case *tftypes.ElementKeyString:
			ret = ret.Index(cty.StringVal(string(*s)))
		case *tftypes.ElementKeyInt:
			ret = ret.Index(cty.NumberIntVal(int64(*s)))
		default:
			ret = append(ret, nil)
		}
	}
	return ret
}
