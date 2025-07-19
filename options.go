package def2env

import (
	"fmt"
	"os"
	"path"
	"strings"
)

var (
	defaultEcspressoConfigs = []string{
		"ecspresso.yml",
		"ecspresso.yaml",
		"ecspresso.json",
		"ecspresso.jsonnet",
	}
)

type Options struct {
	EcspressoOptions
	AllowListOptions
	Command []string `arg:"" required:"" help:"Command and arguments."`
}

type EcspressoOptions struct {
	Config       string `short:"c" env:"ECSPRESSO_CONFIG" default:"ecspresso.yml" help:"ecspresso config file or directory path."`
	ContainerNum uint   `short:"n" default:"0" help:"Container definition index."`
}

type AllowListOptions struct {
	Only []string `xor:"only_all" env:"DEF2ENV_ONLY" required:"" help:"Name of environment variable to pass to the command. Read from file if prefix is 'file://'."`
	All  bool     `xor:"only_all" required:"" help:"Pass all environment variables to the command."`
}

func (options *Options) AfterApply() error {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	if strings.HasPrefix(options.Config, "~/") {
		options.Config = strings.Replace(options.Config, "~", homeDir, 1)
	}

	fi, err := os.Stat(options.Config)

	if err != nil {
		return err
	}

	if fi.IsDir() {
		var err error

		for _, name := range defaultEcspressoConfigs {
			confPath := path.Join(options.Config, name)

			if _, err = os.Stat(confPath); err == nil {
				options.Config = confPath
				break
			}
		}

		if err != nil {
			return fmt.Errorf("%s: no ecspresso config file in directory", options.Config)
		}
	}

	return nil
}
