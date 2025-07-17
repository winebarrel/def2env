package def2env

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/hashicorp/go-envparse"
)

func Run(options *Options) error {
	ecspressoClient, err := NewEcspresso(&options.EcspressoOptions)

	if err != nil {
		return err
	}

	envs, err := ecspressoClient.Environ()

	if err != nil {
		return err
	}

	for _, envFile := range options.EnvFile {
		content, err := os.ReadFile(envFile)

		if err != nil {
			return err
		}

		reader := bytes.NewReader(content)
		envsFromFile, err := envparse.Parse(reader)

		if err != nil {
			return err
		}

		for name, value := range envsFromFile {
			envs[name] = value
		}
	}

	return execCmd(options.Command, envs)
}

func execCmd(cmdArgs []string, extraEnv map[string]string) error {
	name := cmdArgs[0]
	args := []string{}

	if len(cmdArgs) >= 2 {
		args = cmdArgs[1:]
	}

	env := os.Environ()

	for name, value := range extraEnv {
		env = append(env, name+"="+value)
	}

	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	return cmd.Run()
}
