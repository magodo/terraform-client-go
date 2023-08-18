package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/magodo/terraform-client-go/tfclient"
	"github.com/magodo/terraform-client-go/tfclient/configschema"
	"github.com/magodo/terraform-client-go/tfclient/typ"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type JSONPatches []jsonpatch.Patch

func (pl *JSONPatches) String() string {
	return fmt.Sprint(*pl)
}

func (pl *JSONPatches) Set(value string) error {
	p, err := jsonpatch.DecodePatch([]byte(value))
	if err != nil {
		return fmt.Errorf("decoding patch %s: %v", value, err)
	}
	*pl = append(*pl, p)
	return nil

}

type FlagSet struct {
	PluginPath   string
	ResourceType string
	ResourceId   string
	LogLevel     string
	ProviderCfg  string
	StatePatches JSONPatches
	TimeoutSec   int
}

func main() {
	var fset FlagSet
	flag.StringVar(&fset.PluginPath, "path", "", "The path to the plugin")
	flag.StringVar(&fset.ResourceType, "type", "", "The resource type")
	flag.StringVar(&fset.ResourceId, "id", "", "The resource id")
	flag.StringVar(&fset.LogLevel, "log-level", hclog.Error.String(), "Log level")
	flag.StringVar(&fset.ProviderCfg, "cfg", "{}", "The content of provider config block in JSON")
	flag.Var(&fset.StatePatches, "state-patch", "The JSON patch to the state after importing, which will then be used as the prior state for reading. Can be specified multiple times")
	flag.IntVar(&fset.TimeoutSec, "timeout", 0, "Timeout in second. Defaults to no timeout.")

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

	reattach, err := parseReattach(os.Getenv("TF_REATTACH_PROVIDERS"))
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

	importResp, diags := c.ImportResourceState(ctx, typ.ImportResourceStateRequest{
		TypeName: fset.ResourceType,
		ID:       fset.ResourceId,
	})
	if err := showDiags(logger, diags); err != nil {
		return err
	}

	if len(importResp.ImportedResources) != 1 {
		return fmt.Errorf("expect 1 resource, got=%d", len(importResp.ImportedResources))
	}
	res := importResp.ImportedResources[0]

	state := res.State
	if fset.StatePatches != nil {
		for _, patch := range fset.StatePatches {
			b, err := ctyjson.Marshal(state, state.Type())
			if err != nil {
				return fmt.Errorf("marshalling the state: %v", err)
			}
			nb, err := patch.Apply(b)
			if err != nil {
				return fmt.Errorf("patching the state %s: %v", string(b), err)
			}
			state, err = ctyjson.Unmarshal(nb, state.Type())
			if err != nil {
				return fmt.Errorf("unmarshalling the patched state: %v", err)
			}
		}
	}

	readResp, diags := c.ReadResource(ctx, typ.ReadResourceRequest{
		TypeName:     res.TypeName,
		PriorState:   state,
		Private:      res.Private,
		ProviderMeta: cty.Value{},
	})
	if err := showDiags(logger, diags); err != nil {
		return err
	}

	b, err := ctyjson.Marshal(readResp.NewState, configschema.SchemaBlockImpliedType(schResp.ResourceTypes[fset.ResourceType].Block))
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
	var p string
	for k, v := range m {
		c = v
		p = k
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
		return nil, fmt.Errorf("Unknown address type %q for %q", c.Addr.Network, p)
	}
	return &plugin.ReattachConfig{
		Protocol:        plugin.Protocol(c.Protocol),
		ProtocolVersion: c.ProtocolVersion,
		Pid:             c.Pid,
		Test:            c.Test,
		Addr:            addr,
	}, nil
}
