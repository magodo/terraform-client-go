// Duplicating the logic from: github.com/hashicorp/terraform/internal/plugin6/convert/functions.go

package convert

import (
	"encoding/json"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/typ"
)

func FunctionDeclsFromProto(protoFuncs map[string]*tfprotov6.Function) (map[string]typ.FunctionDecl, error) {
	if len(protoFuncs) == 0 {
		return nil, nil
	}

	ret := make(map[string]typ.FunctionDecl, len(protoFuncs))
	for name, protoFunc := range protoFuncs {
		decl, err := FunctionDeclFromProto(protoFunc)
		if err != nil {
			return nil, fmt.Errorf("invalid declaration for function %q: %s", name, err)
		}
		ret[name] = decl
	}
	return ret, nil
}

func FunctionDeclFromProto(fun *tfprotov6.Function) (typ.FunctionDecl, error) {
	var ret typ.FunctionDecl

	ret.Description = fun.Description
	ret.DescriptionKind = tfjson.SchemaDescriptionKind(fun.DescriptionKind)
	ret.Summary = fun.Summary
	ret.DeprecationMessage = fun.DeprecationMessage

	rettyp, err := fun.Return.Type.MarshalJSON()
	if err != nil {
		return ret, fmt.Errorf("marshaling return type: %v", err)
	}
	if err := json.Unmarshal(rettyp, &ret.ReturnType); err != nil {
		return ret, fmt.Errorf("invalid return type constraint: %s", err)
	}

	if len(fun.Parameters) != 0 {
		ret.Parameters = make([]typ.FunctionParam, len(fun.Parameters))
		for i, protoParam := range fun.Parameters {
			param, err := functionParamFromProto(protoParam)
			if err != nil {
				return ret, fmt.Errorf("invalid parameter %d (%q): %s", i, protoParam.Name, err)
			}
			ret.Parameters[i] = param
		}
	}
	if fun.VariadicParameter != nil {
		param, err := functionParamFromProto(fun.VariadicParameter)
		if err != nil {
			return ret, fmt.Errorf("invalid variadic parameter (%q): %s", fun.VariadicParameter.Name, err)
		}
		ret.VariadicParameter = &param
	}

	return ret, nil
}

func functionParamFromProto(param *tfprotov6.FunctionParameter) (typ.FunctionParam, error) {
	var ret typ.FunctionParam
	ret.Name = param.Name
	ret.Description = param.Description
	ret.DescriptionKind = tfjson.SchemaDescriptionKind(param.DescriptionKind)
	ret.AllowNullValue = param.AllowNullValue
	ret.AllowUnknownValues = param.AllowUnknownValues

	paramtype, err := param.Type.MarshalJSON()
	if err != nil {
		return ret, fmt.Errorf("marshaling param type: %v", err)
	}
	if err := json.Unmarshal(paramtype, &ret.Type); err != nil {
		return ret, fmt.Errorf("invalid param type constraint: %s", err)
	}

	return ret, nil
}
