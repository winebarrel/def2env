package def2env

import (
	"os"
	"strings"
)

type Options struct {
	EcspressoOptions
	Command []string `arg:"" required:"" help:"Command and arguments."`
	Only    []string `xor:"only_all" env:"DEF2ENV_ONLY" required:"" help:"Name of environment variable to pass to the command. Read from file if prefix is 'file://'."`
	All     bool     `xor:"only_all" required:"" help:"Pass all environment variables to the command."`
}

type EcspressoOptions struct {
	Config       string `short:"c" env:"ECSPRESSO_CONFIG" default:"ecspresso.yml" type:"existingfile" help:"ecspresso config file path."`
	ContainerNum uint   `short:"n" default:"0" help:"Container definition index."`
}

func (options *Options) AfterApply() error {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	if strings.HasPrefix(options.Config, "~/") {
		options.Config = strings.Replace(options.Config, "~", homeDir, 1)
	}

	return nil
}
