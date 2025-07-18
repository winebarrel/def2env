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
  -c, --config="ecspresso.yml"    ecspresso config file path ($ECSPRESSO_CONFIG).
  -n, --container-num=0           Container definition index.
  -e, --env-file=ENV-FILE         A file listing environment variables to override.
      --only=ONLY,...             A file containing a list of environment variable names to pass to the command.
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

$ def2env --all -- bash -c 'echo $USERNAME'
scott
```
