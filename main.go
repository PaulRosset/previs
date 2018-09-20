package main

import (
	"log"
	"os"

	api "./api"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	api.Start("Dockerfile", cwd+"/")
	//api.IsDockerInstall()
	// imgDocker, err := api.Writter()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "error encountered: %+v", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(imgDocker)
}
