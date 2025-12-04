package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/magodo/terraform-client-go/tfclient"
	"github.com/magodo/terraform-client-go/tfclient/configschema"
	"github.com/magodo/terraform-client-go/tfclient/typ"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type FlagSet struct {
	PluginPath  string
	LogLevel    string
	ProviderCfg string
	TimeoutSec  int
	ActionType  string
	Body        string
}

func main() {
	var fset FlagSet
	flag.StringVar(&fset.PluginPath, "path", "", "The path to the plugin")
	flag.StringVar(&fset.LogLevel, "log-level", hclog.Error.String(), "Log level")
	flag.StringVar(&fset.ProviderCfg, "cfg", "{}", "The content of provider config block in JSON")
	flag.IntVar(&fset.TimeoutSec, "timeout", 0, "Timeout in second. Defaults to no timeout.")
	flag.StringVar(&fset.ActionType, "type", "", "The action type")
	flag.StringVar(&fset.Body, "body", "{}", "The block body for the action")

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

	config, err := ctyjson.Unmarshal([]byte(fset.ProviderCfg), configschema.SchemaBlockImpliedType(schResp.Provider.Block))
	if err != nil {
		return err
	}

	_, diags = c.ConfigureProvider(ctx, typ.ConfigureProviderRequest{
		Config: config,
	})
	if err := showDiags(logger, diags); err != nil {
		return err
	}

	sch, ok := schResp.ActionsCty[fset.ActionType]
	if !ok {
		return fmt.Errorf("no action named %q", fset.ActionType)
	}

	body, err := ctyjson.Unmarshal([]byte(fset.Body), sch)
	if err != nil {
		return err
	}

	resp, diags := c.InvokeAction(ctx, typ.InvokeActionRequest{
		ActionType:        fset.ActionType,
		PlannedActionData: cty.ObjectVal(map[string]cty.Value{"config": body}),
	})
	if err := showDiags(logger, diags); err != nil {
		return err
	}

	for evt := range resp.Events {
		switch evt := evt.(type) {
		case typ.InvokeActionEvent_Progress:
			fmt.Println(evt.Message)
		case typ.InvokeActionEvent_Completed:
			if evt.Diagnostics.HasErrors() {
				fmt.Println(evt.Diagnostics.Err())
			}
		}
	}
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
