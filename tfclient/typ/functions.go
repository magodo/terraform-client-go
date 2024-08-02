// This is derived from github.com/hashicorp/terraform/internal/providers/functions.go

package typ

import (
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

type FunctionDecl struct {
	Parameters        []FunctionParam
	VariadicParameter *FunctionParam
	ReturnType        cty.Type

	Description        string
	DescriptionKind    tfjson.SchemaDescriptionKind
	Summary            string
	DeprecationMessage string
}

type FunctionParam struct {
	Name string // Only for documentation and UI, because arguments are positional
	Type cty.Type

	AllowNullValue     bool
	AllowUnknownValues bool

	Description     string
	DescriptionKind tfjson.SchemaDescriptionKind
}
