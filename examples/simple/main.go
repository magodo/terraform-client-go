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
)

func main() {
	pluginPath := flag.String("path", "", "The path to the plugin")
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

	client, err := tfclient.New(context.TODO(), opts)

	if err != nil {
		log.Fatal(err)
	}
	defer client.Kill()

	if client := client.AsV5Client(); client != nil {
		resp, err := client.GetProviderSchema(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Provider version: %d\n", resp.Provider.Version)
		return
	}

	if client := client.AsV6Client(); client != nil {
		resp, err := client.GetProviderSchema(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Provider version: %d\n", resp.Provider.Version)
		return
	}
}
