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
	PluginPath      string
	LogLevel        string
	ProviderCfg     string
	TimeoutSec      int
	ResourceType    string
	Body            string
	IncludeResource bool
	Limit           int
}

func main() {
	var fset FlagSet
	flag.StringVar(&fset.PluginPath, "path", "", "The path to the plugin")
	flag.StringVar(&fset.LogLevel, "log-level", hclog.Error.String(), "Log level")
	flag.StringVar(&fset.ProviderCfg, "cfg", "{}", "The content of provider config block in JSON")
	flag.IntVar(&fset.TimeoutSec, "timeout", 0, "Timeout in second. Defaults to no timeout.")
	flag.StringVar(&fset.ResourceType, "type", "", "The resource type")
	flag.StringVar(&fset.Body, "body", "{}", "The block body for the list resource")
	flag.BoolVar(&fset.IncludeResource, "include-resource", false, "Should the provider include the full resource object for each result")
	flag.IntVar(&fset.Limit, "limit", 100, "The maximum number of results to return. Default: 100.")

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

	sch, ok := schResp.ListResourceTypes[fset.ResourceType]
	if !ok {
		return fmt.Errorf("no list resource named %q", fset.ResourceType)
	}

	resSch, ok := schResp.ResourceTypes[fset.ResourceType]
	if !ok {
		return fmt.Errorf("no resource named %q", fset.ResourceType)
	}
	resSchCty, ok := schResp.ResourceTypesCty[fset.ResourceType]
	if !ok {
		return fmt.Errorf("no resource named %q", fset.ResourceType)
	}

	body, err := ctyjson.Unmarshal([]byte(fset.Body), configschema.SchemaBlockImpliedType(sch.Block))
	if err != nil {
		return err
	}

	listResp, diags := c.ListResource(ctx, typ.ListResourceRequest{
		TypeName:              fset.ResourceType,
		Config:                cty.ObjectVal(map[string]cty.Value{"config": body}),
		IncludeResourceObject: fset.IncludeResource,
		Limit:                 int64(fset.Limit),
	})
	if err := showDiags(logger, diags); err != nil {
		return err
	}

	datas, ok := listResp.Result.AsValueMap()["data"]
	if !ok {
		return fmt.Errorf(`no "data" in the list resource`)
	}

	b, err := ctyjson.Marshal(datas, cty.List(cty.Object(map[string]cty.Type{
		"display_name": cty.String,
		"state":        resSchCty,
		"identity":     configschema.SchemaNestedAttributeTypeImpliedType(resSch.Identity),
	})))
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
