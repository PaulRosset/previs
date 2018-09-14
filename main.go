package main

import (
	"fmt"
	"os"

	api "./api"
)

func main() {
	// config, err := api.GetConfigFromTravis()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	imgDocker, err := api.Writter()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encountered: %+v", err)
		os.Exit(1)
	}
	fmt.Println(imgDocker)
}
