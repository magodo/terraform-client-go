package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfprotov5/internal/tfplugin5"
)

func RawState(in *tfprotov5.RawState) *tfplugin5.RawState {
	return &tfplugin5.RawState{
		Json:    in.JSON,
		Flatmap: in.Flatmap,
	}
}

// we have to say this next thing to get golint to stop yelling at us about the
// underscores in the function names. We want the function names to match
// actually-generated code, so it feels like fair play. It's just a shame we
// lose golint for the entire file.
//
// This file is not actually generated. You can edit it. Ignore this next line.
// Code generated by hand ignore this next bit DO NOT EDIT.
