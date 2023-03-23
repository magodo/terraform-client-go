// Logic is derived from: github.com/apparentlymart/terraform-provider/tfprovider/internal/common/diagnostics.go (ba059c82d80b5c662af2b0e9de2416903ff05e44)

package typ

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/zclconf/go-cty/cty"
	grpcStatus "google.golang.org/grpc/status"
)

type Diagnostics []Diagnostic

type Diagnostic struct {
	Severity  DiagnosticSeverity
	Summary   string
	Detail    string
	Attribute cty.Path
}

type DiagnosticSeverity rune

const (
	Error   DiagnosticSeverity = 'E'
	Warning DiagnosticSeverity = 'W'
)

func (diags Diagnostics) HasErrors() bool {
	for _, diag := range diags {
		if diag.Severity == Error {
			return true
		}
	}
	return false
}

func (diags Diagnostics) Err() error {
	switch {
	case len(diags) == 0:
		// should never happen, since we don't create this wrapper if
		// there are no diagnostics in the list.
		return nil
	case len(diags) == 1:
		diag := diags[0]
		if diag.Detail == "" {
			return fmt.Errorf(diag.Summary)
		}
		return fmt.Errorf("%s: %s", diag.Summary, diag.Detail)
	default:
		var ret bytes.Buffer
		fmt.Fprintf(&ret, "%d problems:\n", len(diags))
		for _, diag := range diags {
			if diag.Detail == "" {
				fmt.Fprintf(&ret, "\n- %s", diag.Summary)
			} else {
				fmt.Fprintf(&ret, "\n- %s: %s", diag.Summary, diag.Detail)
			}
		}
		return fmt.Errorf(ret.String())
	}
}

func RPCErrorDiagnostics(err error) Diagnostics {
	if err == nil {
		return nil
	}
	var diags Diagnostics
	status, ok := grpcStatus.FromError(err)
	if !ok {
		diags = append(diags, Diagnostic{
			Severity: Error,
			Summary:  "Failed to call provider plugin",
			Detail:   fmt.Sprintf("Provider RPC call failed: %s.", err),
		})
	} else {
		diags = append(diags, Diagnostic{
			Severity: Error,
			Summary:  "Failed to call provider plugin",
			Detail:   fmt.Sprintf("Provider returned RPC error %s: %s.", status.Code(), status.Message()),
		})
	}
	return diags
}

func ErrorDiagnostics(summary string, err error) Diagnostics {
	switch err := err.(type) {
	case nil:
		return nil
	case cty.PathError:
		return Diagnostics{
			{
				Severity:  Error,
				Summary:   summary,
				Detail:    err.Error(),
				Attribute: err.Path,
			},
		}
	default:
		return Diagnostics{
			{
				Severity: Error,
				Summary:  summary,
				Detail:   err.Error(),
			},
		}
	}

}

func FormatError(err error) string {
	switch err := err.(type) {
	case cty.PathError:
		return fmt.Sprintf("%s: %s", FormatCtyPath(err.Path), err.Error())
	default:
		return err.Error()
	}
}

func FormatCtyPath(path cty.Path) string {
	var buf strings.Builder
	for _, step := range path {
		switch step := step.(type) {
		case cty.GetAttrStep:
			buf.WriteString("." + step.Name)
		case cty.IndexStep:
			switch step.Key.Type() {
			case cty.String:
				fmt.Fprintf(&buf, "[%q]", step.Key.AsString())
			case cty.Number:
				fmt.Fprintf(&buf, "[%s]", step.Key.AsBigFloat().Text('f', 0))
			default:
				buf.WriteString("[...]")
			}
		default:
			buf.WriteString("[...]")
		}
	}
	return buf.String()
}
