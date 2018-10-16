package main

import (
	"fmt"
	"os"

	"github.com/srishanbhattarai/previs/api"
)

func whichConfig(args []string) string {
	if len(args) == 1 && args[0] == "-p" {
		fmt.Println("You are using '.previs.yml' file")
		return ".previs.yml"
	}
	fmt.Println("You are using '.travis.yml' file")
	return ".travis.yml"
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encountered: %+v\n", err)
		os.Exit(2)
	}

	configFile := whichConfig(os.Args[1:])

	runner, err := api.FromConfig(cwd + "/" + configFile)
	if err != nil {
		die("An error occurred: %+v\n", err)
	}

	imgDocker, envsVar, err := runner.WriteDockerfile()
	if err != nil {
		api.CleanUnusedDockerfile(cwd, imgDocker)
		die("An error occurred: %+v\n", err)
	}

	err = api.Start(imgDocker, cwd+"/", envsVar)
	if err != nil {
		api.CleanUnusedDockerfile(cwd, imgDocker)
		die("An error occurred: %+v\n", err)
	}
}

// Print a message and exit with a status code of 2.
func die(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(2)
}
