package def2env

import (
	"bufio"
	"net/url"
	"os"
	"os/exec"
	"strings"

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

	if !options.All {
		envs = filterEnvs(envs, options.Only)
	}

	return execCmd(options.Command, envs)
}

func filterEnvs(envs map[string]string, only []string) map[string]string {
	newEnvs := map[string]string{}

	for _, fileOrName := range only {
		if u, err := url.Parse(fileOrName); err == nil && u.Scheme == "file" {
			f, err := os.Open(u.Host)

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
		} else {
			if value, ok := envs[fileOrName]; ok {
				newEnvs[fileOrName] = value
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
