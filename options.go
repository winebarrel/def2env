package def2env

type Options struct {
	EcspressoOptions
	EnvFile []string `short:"e" help:"Environment variable file."`
	Command []string `arg:"" required:"" help:"Command and arguments."`
}

type EcspressoOptions struct {
	TaskDefPath  string `short:"t" required:"" help:"ECS task definition file path."`
	ContainerNum uint   `short:"n" default:"0" help:"Container definition index."`
}
