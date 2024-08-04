package tfclient_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/magodo/terraform-client-go/tfclient"
	"github.com/magodo/terraform-client-go/tfclient/configschema"
	"github.com/magodo/terraform-client-go/tfclient/typ"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

const (
	EnvProviderPath = "TFCLIENTGO_TEST_PROVIDER_PATH"

	// If this is defined, will configure the provider
	EnvProviderConfig = "TFCLIENTGO_TEST_PROVIDER_CONFIG"

	// If both are defined, will import and read the resource
	EnvResourceType = "TFCLIENTGO_TEST_RESOURCE_TYPE"
	EnvResourceId   = "TFCLIENTGO_TEST_RESOURCE_ID"
)

func TestClient(t *testing.T) {
	providerPath, ok := os.LookupEnv(EnvProviderPath)
	if !ok {
		t.Skipf("%q not specified", EnvProviderPath)
	}
	opts := tfclient.Option{
		Cmd:    exec.Command(providerPath),
		Logger: hclog.NewNullLogger(),
	}

	c, err := tfclient.New(opts)
	if err != nil {
		t.Error(err)
	}
	defer c.Close()
}

func BenchmarkClient(b *testing.B) {
	providerPath, ok := os.LookupEnv(EnvProviderPath)
	if !ok {
		b.Skipf("%q not specified", EnvProviderPath)
	}

	b.Run("Spawn", func(b *testing.B) {
		schema := processBySpawn(b, providerPath, false, nil, nil)

		b.Run("WithoutSchema", func(b *testing.B) {
			b.Run("EnableLogStderr", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					processBySpawn(b, providerPath, false, nil, configureImportRead)
				}
			})
			b.Run("DisableLogStderr", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					processBySpawn(b, providerPath, true, nil, configureImportRead)
				}
			})
		})
		b.Run("WithSchema", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.Run("EnableLogStderr", func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						processBySpawn(b, providerPath, false, &schema, configureImportRead)
					}
				})
				b.Run("DisableLogStderr", func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						processBySpawn(b, providerPath, true, &schema, configureImportRead)
					}
				})
			}
		})
	})

	b.Run("Reattach", func(b *testing.B) {
		b.Run("WithoutSchema", func(b *testing.B) {
			reattach, killProvider := spawnProvider(b, providerPath)
			defer killProvider()
			for i := 0; i < b.N; i++ {
				processByReattach(b, reattach, nil, configureImportRead)
			}
		})
		b.Run("WithSchema", func(b *testing.B) {
			reattach, killProvider := spawnProvider(b, providerPath)
			defer killProvider()
			schema := processByReattach(b, reattach, nil, nil)
			for i := 0; i < b.N; i++ {
				processByReattach(b, reattach, &schema, configureImportRead)
			}
		})
	})
}

func processBySpawn(b *testing.B, providerPath string, disableLogStderr bool, schema *typ.GetProviderSchemaResponse, pf processFunc) typ.GetProviderSchemaResponse {
	cmd := exec.Cmd{
		Path: providerPath,
	}
	opts := tfclient.Option{
		Cmd:              &cmd,
		Logger:           hclog.NewNullLogger(),
		ProviderSchema:   schema,
		DisableLogStderr: disableLogStderr,
	}

	c, err := tfclient.New(opts)
	if err != nil {
		b.Error(err)
	}
	defer c.Close()

	schema, diags := c.GetProviderSchema()
	if diags.HasErrors() {
		b.Error(diags.Err())
	}

	if pf != nil {
		if err := pf(b, c); err != nil {
			b.Error(err)
		}
	}

	return *schema
}

func spawnProvider(b *testing.B, providerPath string) (*plugin.ReattachConfig, func()) {
	var stdout bytes.Buffer
	cmd := exec.Cmd{
		Path:   providerPath,
		Args:   []string{providerPath, "-debuggable"},
		Stdout: &stdout,
	}
	if err := cmd.Start(); err != nil {
		b.Error(err)
	}

	time.Sleep(time.Second)

	var reattachStr string
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "TF_REATTACH_PROVIDERS") {
			reattachStr = strings.Split(line, "=")[1]
			reattachStr = reattachStr[1 : len(reattachStr)-1]
			break
		}
	}
	if reattachStr == "" {
		b.Error("can't find reattach string")
	}
	reattach, err := tfclient.ParseReattach(reattachStr)
	if err != nil {
		b.Error(err)
	}
	return reattach, func() {
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			b.Error(err)
		}
		if err := cmd.Wait(); err != nil {
			b.Error(err)
		}
	}
}

func processByReattach(b *testing.B, reattach *plugin.ReattachConfig, schema *typ.GetProviderSchemaResponse, pf processFunc) typ.GetProviderSchemaResponse {
	opts := tfclient.Option{
		Reattach:       reattach,
		Logger:         hclog.NewNullLogger(),
		ProviderSchema: schema,
	}

	c, err := tfclient.New(opts)
	if err != nil {
		b.Error(err)
	}
	defer c.Close()

	schema, diags := c.GetProviderSchema()
	if diags.HasErrors() {
		b.Error(diags.Err())
	}

	if pf != nil {
		if err := pf(b, c); err != nil {
			b.Error(err)
		}
	}

	return *schema
}

type processFunc func(*testing.B, tfclient.Client) error

func configureImportRead(b *testing.B, c tfclient.Client) error {
	ctx := context.Background()

	providerConfig, ok := os.LookupEnv(EnvProviderConfig)
	if !ok {
		b.Logf("Skipping provider config and resource import, as %q not defined", EnvProviderConfig)
		return nil
	}

	schResp, _ := c.GetProviderSchema()
	config, err := ctyjson.Unmarshal([]byte(providerConfig), configschema.SchemaBlockImpliedType(schResp.Provider.Block))
	if err != nil {
		return err
	}

	if _, diags := c.ConfigureProvider(ctx, typ.ConfigureProviderRequest{
		Config: config,
	}); diags.HasErrors() {
		return diags.Err()
	}

	resourceType, ok := os.LookupEnv(EnvResourceType)
	if !ok {
		b.Logf("Skipping resource import, as %q not defined", EnvResourceType)
		return nil
	}
	resourceId, ok := os.LookupEnv(EnvResourceId)
	if !ok {
		b.Logf("Skipping resource import, as %q not defined", EnvResourceId)
		return nil
	}

	for i := 0; i < b.N; i++ {
		importResp, diags := c.ImportResourceState(ctx, typ.ImportResourceStateRequest{
			TypeName: resourceType,
			ID:       resourceId,
		})
		if diags.HasErrors() {
			return fmt.Errorf("importing: %v", diags.Err())
		}

		if len(importResp.ImportedResources) != 1 {
			return fmt.Errorf("expect 1 resource, got=%d", len(importResp.ImportedResources))
		}
		res := importResp.ImportedResources[0]

		if _, diags := c.ReadResource(ctx, typ.ReadResourceRequest{
			TypeName:     res.TypeName,
			PriorState:   res.State,
			Private:      res.Private,
			ProviderMeta: cty.Value{},
		}); diags.HasErrors() {
			return fmt.Errorf("reading: %v", diags.Err())
		}
	}
	return nil
}
