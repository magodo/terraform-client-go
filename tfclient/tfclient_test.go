package tfclient_test

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/magodo/terraform-client-go/tfclient"
)

const (
	EnvProviderPath = "TFCLIENTGO_TEST_PROVIDER_PATH"
)

func BenchmarkSpawnProvider(b *testing.B) {
	providerPath, ok := os.LookupEnv(EnvProviderPath)
	if !ok {
		b.Skipf("%q not specified", EnvProviderPath)
	}

	for i := 0; i < b.N; i++ {
		opts := tfclient.Option{
			Cmd:    exec.Command(providerPath),
			Logger: hclog.NewNullLogger(),
		}

		c, err := tfclient.New(opts)
		if err != nil {
			b.Error(err)
		}

		_, diags := c.GetProviderSchema()
		if diags.HasErrors() {
			b.Error(diags.Err())
		}

		c.Close()
	}
}

func BenchmarkAttachProvider(b *testing.B) {
	providerPath, ok := os.LookupEnv(EnvProviderPath)
	if !ok {
		b.Skipf("%q not specified", EnvProviderPath)
	}

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

	b.ResetTimer()
	b.Log(b.N)
	for i := 0; i < b.N; i++ {
		opts := tfclient.Option{
			Reattach: reattach,
			Logger:   hclog.NewNullLogger(),
		}

		c, err := tfclient.New(opts)
		if err != nil {
			b.Error(err)
		}

		_, diags := c.GetProviderSchema()
		if diags.HasErrors() {
			b.Error(diags.Err())
		}
		c.Close()
	}
}
