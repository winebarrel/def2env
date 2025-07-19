package def2env

import (
	"context"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/kayac/ecspresso/v2"
)

type Ecspresso struct {
	containerDef   types.ContainerDefinition
	secretsmanager *secretsmanager.Client
}

func NewEcspresso(options *EcspressoOptions) (*Ecspresso, error) {
	cliOpts := &ecspresso.CLIOptions{
		ConfigFilePath: options.Config,
	}

	app, err := ecspresso.New(context.Background(), cliOpts)

	if err != nil {
		return nil, err
	}

	tdPath := reflect.ValueOf(*app).FieldByName("config").Elem().FieldByName("TaskDefinitionPath").String()
	td, err := app.LoadTaskDefinition(tdPath)

	if err != nil {
		return nil, err
	}

	containerNum := int(options.ContainerNum)

	if containerNum >= len(td.ContainerDefinitions) {
		return nil, fmt.Errorf("invalid container index: %d", containerNum)
	}

	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		return nil, err
	}

	client := &Ecspresso{
		containerDef:   td.ContainerDefinitions[containerNum],
		secretsmanager: secretsmanager.NewFromConfig(cfg),
	}

	return client, nil
}

func (client *Ecspresso) Environ(allowlist AllowList) (map[string]string, error) {
	envs := map[string]string{}
	client.appendEnvironment(client.containerDef.Environment, allowlist, envs)
	err := client.appendSecrets(client.containerDef.Secrets, allowlist, envs)

	if err != nil {
		return nil, err
	}

	return envs, nil
}

func (*Ecspresso) appendEnvironment(environment []types.KeyValuePair, allowlist AllowList, envs map[string]string) {
	for _, e := range environment {
		name := aws.ToString(e.Name)
		value := aws.ToString(e.Value)

		if allowlist.IsAllowed(name) {
			envs[name] = value
		}
	}
}

func (client *Ecspresso) appendSecrets(secrets []types.Secret, allowlist AllowList, envs map[string]string) error {
	nameByArn := map[string]string{}
	arns := []string{}
	valueByArn := map[string]string{}

	for _, s := range secrets {
		name := aws.ToString(s.Name)
		arn := aws.ToString(s.ValueFrom)

		if allowlist.IsAllowed(name) {
			nameByArn[arn] = name
			arns = append(arns, arn)
		}
	}

	for {
		var secretIdList []string

		if len(arns) > 20 {
			secretIdList = arns[0:20]
			arns = arns[20:]
		} else {
			secretIdList = arns
			arns = nil
		}

		input := &secretsmanager.BatchGetSecretValueInput{
			SecretIdList: secretIdList,
		}

		output, err := client.secretsmanager.BatchGetSecretValue(context.Background(), input)

		if err != nil {
			return err
		}

		for _, sv := range output.SecretValues {
			value := aws.ToString(sv.SecretString)
			arn := aws.ToString(sv.ARN)
			valueByArn[arn] = value
		}

		if arns == nil {
			break
		}
	}

	for arn, name := range nameByArn {
		envs[name] = valueByArn[arn]
	}

	return nil
}
