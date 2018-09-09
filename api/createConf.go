package api

import (
	"fmt"
	"os"

	"github.com/satori/go.uuid"
)

func createDockerFile() (*os.File, string) {
	uid := uuid.Must(uuid.NewV4())
	imgDocker := ".dockerfile" + uid.String()
	cwd, errOnCwd := os.Getwd()
	if errOnCwd != nil {
		panic(errOnCwd)
	}
	file, err := os.OpenFile(cwd+"/"+imgDocker, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return file, imgDocker
}

//config map[interface{}]interface{}

func Writter() {
	config, _ := GetConfigFromTravis()
	fmt.Println(config)
	file, _ := createDockerFile()
	_, err := file.WriteString("FROM " + config["language"].(string))
	fmt.Println(err)
}
