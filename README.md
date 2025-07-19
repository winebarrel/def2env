# def2env

[![CI](https://github.com/winebarrel/def2env/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/def2env/actions/workflows/ci.yml)

def2env is a tool that extracts environment variables from ECS task definitions and executes commands.
(Supports AWS Secrets Manager variable expansion)

## Usage

```
Usage: def2env --only=ONLY,... --all <command> ... [flags]

Arguments:
  <command> ...    Command and arguments.

Flags:
  -h, --help                      Show help.
  -c, --config="ecspresso.yml"    ecspresso config file path
                                  ($ECSPRESSO_CONFIG).
  -n, --container-num=0           Container definition index.
      --only=ONLY,...             Name of environment variable to pass to the
                                  command. Read from file if prefix is 'file://'
                                  ($DEF2ENV_ONLY).
      --all                       Pass all environment variables to the command.
      --version
```

## Example

```
$ ls
ecs-service-def.jsonnet
ecs-task-def.jsonnet
ecspresso.yml

$ grep environment -A 2 ecs-task-def.jsonnet
  environment: [
    { name: 'USERNAME', value: 'scott' },

$ def2env --only USERNAME -- bash -c 'echo $USERNAME'
scott
```
