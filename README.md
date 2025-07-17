# def2env

def2env is a tool that extracts environment variables from ECS task definitions and executes commands.
(Supports AWS Secrets Manager variable expansion)

## Usage

```
Usage: def2env <command> ... [flags]

Arguments:
  <command> ...    Command and arguments.

Flags:
  -h, --help                      Show help.
  -c, --config="ecspresso.yml"    ecspresso config file path.
  -n, --container-num=0           Container definition index.
  -e, --env-file=ENV-FILE,...     Environment variable file.
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

$ def2env -- bash -c 'echo $USERNAME'
scott
```
