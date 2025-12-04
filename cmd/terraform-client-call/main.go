package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/magodo/terraform-client-go/tfclient"
	"github.com/magodo/terraform-client-go/tfclient/typ"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type FlagSet struct {
	PluginPath   string
	LogLevel     string
	TimeoutSec   int
	FunctionName string
	FunctionArgs stringSlice
}

type stringSlice []string

func (l *stringSlice) String() string {
	return fmt.Sprintf("[%s]", strings.Join(*l, ", "))
}

func (l *stringSlice) Set(value string) error {
	*l = append(*l, value)
	return nil
}

func main() {
	var fset FlagSet
	flag.StringVar(&fset.PluginPath, "path", "", "The path to the plugin")
	flag.StringVar(&fset.LogLevel, "log-level", hclog.Error.String(), "Log level")
	flag.IntVar(&fset.TimeoutSec, "timeout", 0, "Timeout in second. Defaults to no timeout.")
	flag.StringVar(&fset.FunctionName, "func", "", "The name of the function")
	flag.Var(&fset.FunctionArgs, "arg", "The argument of the function (can be specified multiple times)")

	flag.Parse()

	logger := hclog.New(&hclog.LoggerOptions{
		Output: hclog.DefaultOutput,
		Level:  hclog.LevelFromString(fset.LogLevel),
		Name:   filepath.Base(fset.PluginPath),
	})

	if err := realMain(logger, fset); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func realMain(logger hclog.Logger, fset FlagSet) error {
	opts := tfclient.Option{
		Cmd:    exec.Command(fset.PluginPath),
		Logger: logger,
	}

	reattach, err := tfclient.ParseReattach(os.Getenv("TF_REATTACH_PROVIDERS"))
	if err != nil {
		return err
	}
	if reattach != nil {
		opts.Cmd = nil
		opts.Reattach = reattach
	}

	c, err := tfclient.New(opts)
	if err != nil {
		return err
	}
	defer c.Close()

	ctx := context.TODO()
	var cancel context.CancelFunc
	if fset.TimeoutSec > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(fset.TimeoutSec))
		defer cancel()
	}

	schResp, diags := c.GetProviderSchema()
	if err := showDiags(logger, diags); err != nil {
		return err
	}

	fsch, ok := schResp.Functions[fset.FunctionName]
	if !ok {
		return fmt.Errorf("no function named %q", fset.FunctionName)
	}

	var args []cty.Value
	if len(fset.FunctionArgs) != 0 {
		args = make([]cty.Value, len(fset.FunctionArgs))
		for i, rawArg := range fset.FunctionArgs {
			idx := int64(i)

			var argTy cty.Type
			if i < len(fsch.Parameters) {
				argTy = fsch.Parameters[i].Type
			} else {
				if fsch.VariadicParameter == nil {
					return fmt.Errorf("too many arguments for non-variadic function beyond %d-th argument", idx)
				}
				argTy = fsch.VariadicParameter.Type
			}

			argVal, err := ctyjson.Unmarshal([]byte(rawArg), argTy)
			if err != nil {
				return fmt.Errorf("parsing %d-th argument: %v", idx, err)
			}

			args[i] = argVal
		}
	}

	resp, diags := c.CallFunction(ctx, typ.CallFunctionRequest{
		FunctionName: fset.FunctionName,
		Arguments:    args,
	})
	if err := showDiags(logger, diags); err != nil {
		return err
	}

	if resp.Err != nil {
		fmt.Printf("Function returns error: %v\n", resp.Err)
		return nil
	}

	b, err := ctyjson.Marshal(resp.Result, fsch.ReturnType)
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}

func showDiags(logger hclog.Logger, diags typ.Diagnostics) error {
	for _, diag := range diags {
		if diag.Severity == typ.Error {
			return fmt.Errorf("%s: %s", diag.Summary, diag.Detail)
		}
	}
	if len(diags) != 0 {
		logger.Warn(diags.Err().Error())
	}
	return nil
}
