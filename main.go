package main

import (
	api "./api"
)

func main() {
	api.IsDockerInstall()
	// imgDocker, err := api.Writter()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "error encountered: %+v", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(imgDocker)
}
