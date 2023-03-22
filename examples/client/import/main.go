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
	"github.com/magodo/terraform-client-go/tfclient/client"
	"github.com/magodo/terraform-client-go/tfclient/configschema"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

func main() {
	pluginPath := flag.String("path", "", "The path to the plugin")
	resourceType := flag.String("type", "", "The resource type")
	resourceId := flag.String("id", "", "The resource id")
	flag.Parse()

	logger := hclog.New(&hclog.LoggerOptions{
		Output: hclog.DefaultOutput,
		Level:  hclog.Info,
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

	_, diags = c.ConfigureProvider(ctx, client.ConfigureProviderRequest{
		Config: config,
	})
	showDiags(diags)

	importResp, diags := c.ImportResourceState(ctx, client.ImportResourceStateRequest{
		TypeName: *resourceType,
		ID:       *resourceId,
	})
	showDiags(diags)

	if len(importResp.ImportedResources) != 1 {
		log.Fatalf("expect 1 resource, got=%d", len(importResp.ImportedResources))
	}
	res := importResp.ImportedResources[0]
	readResp, diags := c.ReadResource(ctx, client.ReadResourceRequest{
		TypeName:   res.TypeName,
		PriorState: res.State,
		//Private:      res.Private,
		//ProviderMeta: cty.Value{},
	})
	showDiags(diags)
	fmt.Println(readResp.NewState.GoString())
}

func showDiags(diags client.Diagnostics) {
	for _, diag := range diags {
		if diag.Severity == client.Error {
			log.Fatal(diag.Summary + ": " + diag.Detail)
		}
	}
	if len(diags) != 0 {
		fmt.Println(diags.Err())
	}
}
