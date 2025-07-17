package main

import (
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/kayac/ecspresso/v2"
	"github.com/winebarrel/def2env"
)

var (
	version          string
	ecspressoVersion string
)

func init() {
	ecspresso.Version = ecspressoVersion
	log.SetFlags(0)
}

func parseArgs() *def2env.Options {
	var cli struct {
		def2env.Options
		Version kong.VersionFlag
	}

	parser := kong.Must(&cli, kong.Vars{"version": version})
	parser.Model.HelpFlag.Help = "Show help."
	_, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)

	return &cli.Options
}

func main() {
	options := parseArgs()
	err := def2env.Run(options)

	if err != nil {
		log.Fatalf("def2env error: %s", err)
	}
}
