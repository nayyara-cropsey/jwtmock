package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/nayyara-cropsey/jwtmock/cmd"
)

const defaultConfigFile = "config.yaml"

func main() {
	configFile := defaultConfigFile
	flag.StringVar(&configFile, "config", defaultConfigFile, "config file")
	flag.Parse()

	// handle shutdown by user pressing Ctrl+C
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := cmd.Serve(ctx, configFile); err != nil {
		fmt.Println("Error in command:", err)
		defer os.Exit(1)
	}
}
