// Logic is derived from: github.com/apparentlymart/terraform-provider/tfprovider/internal/protocol6/diagnostics.go (ba059c82d80b5c662af2b0e9de2416903ff05e44)

package convert

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/magodo/terraform-client-go/tfclient/client"
	"github.com/zclconf/go-cty/cty"
)

func DecodeDiagnostics(raws []*tfprotov6.Diagnostic) client.Diagnostics {
	if len(raws) == 0 {
		return nil
	}
	diags := make(client.Diagnostics, 0, len(raws))
	for _, raw := range raws {
		diag := client.Diagnostic{
			Summary:   raw.Summary,
			Detail:    raw.Detail,
			Attribute: DecodeAttributePath(raw.Attribute),
		}

		switch raw.Severity {
		case tfprotov6.DiagnosticSeverityError:
			diag.Severity = client.Error
		case tfprotov6.DiagnosticSeverityWarning:
			diag.Severity = client.Warning
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
