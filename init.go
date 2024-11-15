package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/nomad/api"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	flag "github.com/spf13/pflag"
	"github.com/zerodha/logf"
)

// initLogger initializes logger instance.
func initLogger(ko *koanf.Koanf) logf.Logger {
	opts := logf.Opts{EnableCaller: true}
	if ko.String("app.log_level") == "debug" {
		opts.Level = logf.DebugLevel
	}
	if ko.String("app.env") == "dev" {
		opts.EnableColor = true
	}
	return logf.New(opts)
}

// initConfig loads config to `ko` object.
func initConfig(cfgDefault string, envPrefix string) (*koanf.Koanf, error) {
	var (
		ko = koanf.New(".")
		f  = flag.NewFlagSet("front", flag.ContinueOnError)
	)

	// Configure Flags.
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}

	// Register `--config` flag.
	cfgPath := f.String("config", cfgDefault, "Path to a config file to load.")

	// Parse and Load Flags.
	err := f.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	// Load the config files from the path provided.
	err = ko.Load(file.Provider(*cfgPath), toml.Parser())
	if err != nil {
		return nil, err
	}

	// Load environment variables if the key is given
	// and merge into the loaded config.
	if envPrefix != "" {
		err = ko.Load(env.Provider(envPrefix, ".", func(s string) string {
			return strings.Replace(strings.ToLower(
				strings.TrimPrefix(s, envPrefix)), "__", ".", -1)
		}), nil)
		if err != nil {
			return nil, err
		}
	}

	return ko, nil
}

// initNomadClient initialises a Nomad API client.
func initNomadClient(secretId string) (*api.Client, error) {
	config := api.DefaultConfig()

	// add secretId to config if it's not empty
	if len(secretId) > 0 && secretId[0] != '<' {
		config.SecretID = secretId
	} else {
		fmt.Printf("WARNING: Initializing Nomad client WITHOUT SecretID since it's value is '%s'. Nomad API requests will be unauthenticated\n", secretId)
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func initOpts(ko *koanf.Koanf) Opts {
	return Opts{
		nomadSecretId:          ko.MustString("app.nomad_secret_id"),
		refreshInterval:        ko.MustDuration("app.refresh_interval"),
		removeAllocInterval:    ko.MustDuration("app.remove_alloc_interval"),
		nomadDataDir:           ko.MustString("app.nomad_data_dir"),
		vectorConfigDir:        ko.MustString("app.vector_config_dir"),
		extraTemplatesDir:      ko.String("app.extra_templates_dir"),
		nomadOutspeedServerJob: ko.MustString("app.nomad_outspeed_server_job"),
		nomadJobIdPrefix:       ko.String("app.nomad_job_id_prefix"),
		lokiEndpoint:           ko.String("app.loki.endpoint"),
		lokiUser:               ko.String("app.loki.user"),
		lokiPassword:           ko.String("app.loki.password"),
	}
}
