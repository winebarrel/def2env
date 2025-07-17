package def2env

import (
	"os"
	"strings"
)

type Options struct {
	EcspressoOptions
	EnvFile []string `short:"e" help:"Environment variable file."`
	Command []string `arg:"" required:"" help:"Command and arguments."`
}

type EcspressoOptions struct {
	Config       string `short:"c" default:"ecspresso.yml" help:"ecspresso config file path."`
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
