package main

import (
	"context"
	"fmt"
	"jwt-mock/cmd"
	"jwt-mock/types"
	"os"
	"os/signal"
	"strings"

	"github.com/spf13/viper"
)

const (
	configEnvPrefix = "JWT_MOCK"
	configFormat    = "yaml"
	// using a hard-coded config to avoid shipping a config.yaml file
	// this just provides the defaults - which are overridden via ENV vars.
	defaultConfig = `
port: 80
key_length: 1024
cert_life_days: 1
`
)

func main() {
	if err := executeCmd(); err != nil {
		os.Exit(1)
	}
}

// executeCmd executes the main cobra command and returns any errors.
func executeCmd() error {
	rootCmd := cmd.RootCmd

	// load config
	config, err := initConfig()
	if err != nil {
		rootCmd.Println("Error loading config:", err)
		return err
	}

	// create context with config set on it
	configCtx := context.WithValue(context.Background(), cmd.ConfigKey, config)

	// handle shutdown by user pressing Ctrl+C
	ctx, cancel := signal.NotifyContext(configCtx, os.Interrupt)
	defer cancel()

	// executeCmd command
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		rootCmd.Println("Error in command:", err)
		return err
	}

	return nil
}

// initConfig initializes viper with default config and loads any overrides via ENV vars.
func initConfig() (*types.Config, error) {
	viper.AutomaticEnv()
	viper.SetConfigType(configFormat)
	viper.SetEnvPrefix(configEnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadConfig(strings.NewReader(defaultConfig)); err != nil {
		return nil, fmt.Errorf("initialize config: %w", err)
	}

	config := new(types.Config)
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return config, nil
}
