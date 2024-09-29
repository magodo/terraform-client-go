package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/magodo/terraform-client-go/internal/find"
	"github.com/magodo/terraform-client-go/tfclient"
	"github.com/magodo/terraform-client-go/tfclient/typ"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type FlagSet struct {
	PluginPath   string
	LogLevel     string
	ProviderCfg  string
	TimeoutSec   int
	ModuleDir    string
	ResourceAddr string
	ModuleAddr   string
}

func main() {
	var fset FlagSet
	flag.StringVar(&fset.PluginPath, "path", "", "The path to the plugin")
	flag.StringVar(&fset.LogLevel, "log-level", hclog.Error.String(), "Log level")
	flag.StringVar(&fset.ProviderCfg, "cfg", "{}", "The content of provider config block in JSON")
	flag.IntVar(&fset.TimeoutSec, "timeout", 0, "Timeout in second. Defaults to no timeout.")
	flag.StringVar(&fset.ModuleDir, "module-dir", "", "Path to the root module")
	flag.StringVar(&fset.ResourceAddr, "resource-addr", "", "The resource address (e.g. azurerm_resource_group.test)")
	flag.StringVar(&fset.ModuleAddr, "module-addr", "", "The module address (e.g. mod1.mod2). Defaults to the root module")

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
	ctx := context.Background()

	// Reading the state file, via Terraform
	tfpath, err := find.FindTF(ctx, version.MustConstraints(version.NewConstraint(">=1.0.0")))
	if err != nil {
		return fmt.Errorf("finding terraform executable: %v", err)
	}
	tf, err := tfexec.NewTerraform(fset.ModuleDir, tfpath)
	if err != nil {
		return fmt.Errorf("new terraform via terraform-exec: %v", err)
	}
	state, err := tf.ShowStateFile(ctx, "terraform.tfstate")
	if err != nil {
		return fmt.Errorf("show state file: %v", err)
	}

	// Find the resource state value
	if state == nil {
		return fmt.Errorf("state is nil")
	}
	if state.Values == nil {
		return fmt.Errorf("state.Values is nil")
	}
	if state.Values.RootModule == nil {
		return fmt.Errorf("state.Values.RootModule is nil")
	}

	var moduleAddrsRev []string
	if fset.ModuleAddr != "" {
		moduleAddrsRev = strings.Split(fset.ModuleAddr, ".")
		slices.Reverse(moduleAddrsRev)
	}

	var moduleAddressSegs []string
	stateModule := state.Values.RootModule
	for _, maddr := range moduleAddrsRev {
		moduleAddressSegs = append([]string{"module", maddr}, moduleAddressSegs...)
		var found bool
		for _, cm := range stateModule.ChildModules {
			if cm.Address == strings.Join(moduleAddressSegs, ".") {
				stateModule = cm
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("failed to find module %q in state", strings.Join(moduleAddressSegs, "."))
		}
	}
	rt, rn, ok := strings.Cut(fset.ResourceAddr, ".")
	if !ok {
		return fmt.Errorf("invalid resource address specified: %s", fset.ResourceAddr)
	}

	var stateVal map[string]interface{}
	var schemaVersion uint64
	for _, rs := range stateModule.Resources {
		if rs.Type != rt {
			continue
		}
		if rs.Name != rn {
			continue
		}
		stateVal = rs.AttributeValues
		schemaVersion = rs.SchemaVersion
		break
	}
	if stateVal == nil {
		return fmt.Errorf("failed to find resource %s in module %s", fset.ResourceAddr, fset.ModuleAddr)
	}
	stateValJSON, err := json.Marshal(stateVal)
	if err != nil {
		return fmt.Errorf("JSON marshal the state value: %v", err)
	}

	// Upgrade state
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

	var cancel context.CancelFunc
	if fset.TimeoutSec > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(fset.TimeoutSec))
		defer cancel()
	}

	schResp, diags := c.GetProviderSchema()
	if err := showDiags(logger, diags); err != nil {
		return err
	}

	req := typ.UpgradeResourceStateRequest{
		TypeName:     rt,
		Version:      int64(schemaVersion),
		RawStateJSON: stateValJSON,
	}
	resp, diags := c.UpgradeResourceState(ctx, req)
	if err := showDiags(logger, diags); err != nil {
		return err
	}

	b, err := ctyjson.Marshal(resp.UpgradedState, schResp.ResourceTypesCty[rt])
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
