package def2env

import (
	"os"
	"os/exec"
)

func Run(options *Options) error {
	ecspressoClient, err := NewEcspresso(&options.EcspressoOptions)

	if err != nil {
		return err
	}

	allowlist, err := NewAllowList(&options.AllowListOptions)

	if err != nil {
		return err
	}

	envs, err := ecspressoClient.Environ(allowlist)

	if err != nil {
		return err
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
