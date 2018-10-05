package main

import (
	"fmt"
	"os"

	"github.com/PaulRosset/previs/api"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encountered: %+v\n", err)
		os.Exit(2)
	}
	imgDocker, envsVar, err := api.Writter()
	if err != nil {
		api.CleanUnusedDockerfile(cwd, imgDocker)
		fmt.Fprintf(os.Stderr, "error encountered: %+v\n", err)
		os.Exit(2)
	}
	err = api.Start(imgDocker, cwd+"/", envsVar)
	if err != nil {
		api.CleanUnusedDockerfile(cwd, imgDocker)
		fmt.Fprintf(os.Stderr, "error encountered: %+v\n", err)
		os.Exit(2)
	}
}
