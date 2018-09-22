package main

import (
	"fmt"
	"os"

	"github.com/PaulRosset/previs/api"
)

func main() {
	api.IsDockerInstall()
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encountered: %+v", err)
		os.Exit(2)
	}
	imgDocker, err := api.Writter()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encountered: %+v", err)
		os.Exit(2)
	}
	api.Start(imgDocker, cwd+"/")
}
