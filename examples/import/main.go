package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/magodo/terraform-client-go/tfclient"
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

	client, err := tfclient.New(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Kill()

	// This example only targets to v5 provider
	c := client.AsV5Client()
	if c == nil {
		log.Fatal("not a v5 provider")
	}

	ctx := context.TODO()

	resp, err := c.GetProviderSchema(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	showDiags(resp.Diagnostics)

	providerConfigType := resp.Provider.Block.ValueType()
	providerConfigVal, err := tftypes.ValueFromJSONWithOpts([]byte(`{"features": []}`), providerConfigType, tftypes.ValueFromJSONOpts{})
	if err != nil {
		log.Fatal(err)
	}
	providerConfig, err := tfprotov5.NewDynamicValue(providerConfigType, providerConfigVal)
	if err != nil {
		log.Fatal(err)
	}
	configResp, err := c.ConfigureProvider(ctx, &tfprotov5.ConfigureProviderRequest{
		Config: &providerConfig,
	})
	if err != nil {
		log.Fatal(err)
	}
	showDiags(configResp.Diagnostics)

	importResp, err := c.ImportResourceState(ctx, &tfprotov5.ImportResourceStateRequest{
		TypeName: *resourceType,
		ID:       *resourceId,
	})
	if err != nil {
		log.Fatal(err)
	}
	showDiags(importResp.Diagnostics)
	if len(importResp.ImportedResources) != 1 {
		log.Fatalf("expect one resource imported, got=%d", len(importResp.ImportedResources))
	}
	state := importResp.ImportedResources[0].State

	readResp, err := c.ReadResource(ctx, &tfprotov5.ReadResourceRequest{
		TypeName:     *resourceType,
		CurrentState: state,
	})
	if err != nil {
		log.Fatal(err)
	}
	showDiags(readResp.Diagnostics)

	readVal, err := readResp.NewState.Unmarshal(resp.ResourceSchemas[*resourceType].ValueType())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(readVal.String())
	return
}

func showDiags(diags []*tfprotov5.Diagnostic) {
	for _, diag := range diags {
		msg := "[" + diag.Severity.String() + "]" + diag.Summary + ": " + diag.Detail + "\n"
		if diag.Severity.String() == "ERROR" {
			log.Fatal(msg)
		}
		fmt.Fprint(os.Stderr, msg)
	}
}
