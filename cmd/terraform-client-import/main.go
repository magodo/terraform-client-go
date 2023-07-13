package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
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
	providerCfg := flag.String("cfg", "{}", "The content of provider config block in JSON")
	statePatch := flag.String("state-patch", "", "The patch to the state after importing, which will then be used as the prior state for reading")
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

	reattach, err := parseReattach(os.Getenv("TF_REATTACH_PROVIDERS"))
	if err != nil {
		log.Fatal(err)
	}
	if reattach != nil {
		opts.Cmd = nil
		opts.Reattach = reattach
	}

	c, err := tfclient.New(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	ctx := context.TODO()

	schResp, diags := c.GetProviderSchema()
	showDiags(diags)

	config, err := ctyjson.Unmarshal([]byte(*providerCfg), configschema.SchemaBlockImpliedType(schResp.Provider.Block))
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

	state := res.State
	if *statePatch != "" {
		b, err := ctyjson.Marshal(state, state.Type())
		if err != nil {
			log.Fatalf("marshalling the state: %v", err)
		}
		b, err = jsonpatch.MergePatch(b, []byte(*statePatch))
		if err != nil {
			log.Fatalf("patching the state: %v", err)
		}
		state, err = ctyjson.Unmarshal(b, state.Type())
		if err != nil {
			log.Fatalf("unmarshalling the patched state: %v", err)
		}
	}

	readResp, diags := c.ReadResource(ctx, typ.ReadResourceRequest{
		TypeName:     res.TypeName,
		PriorState:   state,
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

func parseReattach(in string) (*plugin.ReattachConfig, error) {
	if in == "" {
		return nil, nil
	}

	type reattachConfig struct {
		Protocol        string
		ProtocolVersion int
		Addr            struct {
			Network string
			String  string
		}
		Pid  int
		Test bool
	}
	var m map[string]reattachConfig
	err := json.Unmarshal([]byte(in), &m)
	if err != nil {
		return nil, fmt.Errorf("Invalid format for TF_REATTACH_PROVIDERS: %w", err)
	}
	if len(m) != 1 {
		return nil, fmt.Errorf("expect only one of provider specified in the TF_REATTACH_PROVIDERS, got=%d", len(m))
	}

	var c reattachConfig
	for _, v := range m {
		c = v
	}

	var addr net.Addr
	switch c.Addr.Network {
	case "unix":
		addr, err = net.ResolveUnixAddr("unix", c.Addr.String)
		if err != nil {
			return nil, fmt.Errorf("Invalid unix socket path %q: %w", c.Addr.String, err)
		}
	case "tcp":
		addr, err = net.ResolveTCPAddr("tcp", c.Addr.String)
		if err != nil {
			return nil, fmt.Errorf("Invalid TCP address %q: %w", c.Addr.String, err)
		}
	default:
		return nil, fmt.Errorf("Unknown address type %q for %q", c.Addr.Network)
	}
	return &plugin.ReattachConfig{
		Protocol:        plugin.Protocol(c.Protocol),
		ProtocolVersion: c.ProtocolVersion,
		Pid:             c.Pid,
		Test:            c.Test,
		Addr:            addr,
	}, nil
}
