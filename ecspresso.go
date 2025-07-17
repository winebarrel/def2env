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

func (client *Ecspresso) Environ() (map[string]string, error) {
	envs := map[string]string{}
	client.appendEnvironment(client.containerDef.Environment, envs)
	err := client.appendSecrets(client.containerDef.Secrets, envs)

	if err != nil {
		return nil, err
	}

	return envs, nil
}

func (*Ecspresso) appendEnvironment(environment []types.KeyValuePair, envs map[string]string) {
	for _, e := range environment {
		envs[*e.Name] = *e.Value
	}
}

func (client *Ecspresso) appendSecrets(secrets []types.Secret, envs map[string]string) error {
	names := []string{}
	valueFroms := []string{}
	values := []string{}

	for _, s := range secrets {
		names = append(names, *s.Name)
		valueFroms = append(valueFroms, *s.ValueFrom)
	}

	for {
		var secretIdList []string

		if len(valueFroms) > 20 {
			secretIdList = valueFroms[0:20]
			valueFroms = valueFroms[20:]
		} else {
			secretIdList = valueFroms
			valueFroms = nil
		}

		input := &secretsmanager.BatchGetSecretValueInput{
			SecretIdList: secretIdList,
		}

		output, err := client.secretsmanager.BatchGetSecretValue(context.Background(), input)

		if err != nil {
			return err
		}

		for _, sv := range output.SecretValues {
			values = append(values, aws.ToString(sv.SecretString))
		}

		if valueFroms == nil {
			break
		}
	}

	for i, n := range names {
		envs[n] = values[i]
	}

	return nil
}
