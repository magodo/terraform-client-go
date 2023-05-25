package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/magodo/terraform-client-go/tfclient"
	"github.com/magodo/terraform-client-go/tfclient/configschema"
	"github.com/magodo/terraform-client-go/tfclient/typ"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

func main() {
	pluginPath := flag.String("path", "", "The path to the plugin")
	resourceType := flag.String("type", "", "The resource type")
	resourceId := flag.String("id", "", "The resource id")
	logLevel := flag.String("log-level", hclog.Error.String(), "Log level")
	flag.Parse()

	logger := hclog.New(&hclog.LoggerOptions{
		Output: hclog.DefaultOutput,
		Level:  hclog.LevelFromString(*logLevel),
		Name:   filepath.Base(*pluginPath),
	})
	opts := tfclient.Option{
		Cmd:    exec.Command(*pluginPath),
		Logger: logger,
	}

	c, err := tfclient.New(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	ctx := context.TODO()

	schResp, diags := c.GetProviderSchema()
	showDiags(diags)

	config, err := ctyjson.Unmarshal([]byte(`{"features": []}`), configschema.SchemaBlockImpliedType(schResp.Provider.Block))
	if err != nil {
		log.Fatal(err)
	}

	_, diags = c.ConfigureProvider(ctx, typ.ConfigureProviderRequest{
		Config: config,
	})
	showDiags(diags)

	importResp, diags := c.ImportResourceState(ctx, typ.ImportResourceStateRequest{
		TypeName: *resourceType,
		ID:       *resourceId,
	})
	showDiags(diags)

	if len(importResp.ImportedResources) != 1 {
		log.Fatalf("expect 1 resource, got=%d", len(importResp.ImportedResources))
	}
	res := importResp.ImportedResources[0]
	readResp, diags := c.ReadResource(ctx, typ.ReadResourceRequest{
		TypeName:     res.TypeName,
		PriorState:   res.State,
		Private:      res.Private,
		ProviderMeta: cty.Value{},
	})
	showDiags(diags)

	b, err := ctyjson.Marshal(readResp.NewState, configschema.SchemaBlockImpliedType(schResp.ResourceTypes[*resourceType].Block))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

func showDiags(diags typ.Diagnostics) {
	for _, diag := range diags {
		if diag.Severity == typ.Error {
			log.Fatal(diag.Summary + ": " + diag.Detail)
		}
	}
	if len(diags) != 0 {
		fmt.Println(diags.Err())
	}
}
