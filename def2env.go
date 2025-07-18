package def2env

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-envparse"
	"github.com/kayac/ecspresso/v2"
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

	if !options.All {
		envs = filterEnvs(envs, options.Only)
	}

	return execCmd(options.Command, envs)
}

func filterEnvs(envs map[string]string, onlyFiles []string) map[string]string {
	newEnvs := map[string]string{}

	for _, file := range onlyFiles {
		f, err := os.Open(file)
		if err != nil {
			ecspresso.LogWarn("file loading skipped: %s", err)
			continue
		}

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			name := strings.TrimSpace(scanner.Text())

			if name == "" || strings.HasPrefix(name, "#") {
				continue
			}

			if value, ok := envs[name]; ok {
				newEnvs[name] = value
			}
		}

	}

	return newEnvs
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
