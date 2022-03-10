package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/nayyara-cropsey/jwt-mock/cmd"
	"github.com/nayyara-cropsey/jwt-mock/log"
	"github.com/nayyara-cropsey/jwt-mock/types"
)

const defaultConfigFile = "config.yaml"

var configFile = defaultConfigFile

func main() {
	flag.StringVar(&configFile, "config", defaultConfigFile, "config file")
	flag.Parse()

	if err := executeCmd(); err != nil {
		os.Exit(1)
	}
}

// executeCmd executes the main cobra command and returns any errors.
func executeCmd() error {
	cfg, err := types.Load(configFile)
	if err != nil {
		return err
	}

	// handle shutdown by user pressing Ctrl+C
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger := log.NewLogger(log.WithLevelStr(cfg.LogLevel))
	if err := cmd.Execute(ctx, cfg, logger); err != nil {
		fmt.Println("Error in command:", err)
		return err
	}

	return nil
}
