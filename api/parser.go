package api

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func openConfigTravis() ([]byte, error) {
	cwd, errOnPath := os.Getwd()
	if errOnPath != nil {
		return nil, errOnPath
	}
	fmt.Println(cwd)
	file, errOnOpen := os.Open(cwd + "/.travis.yml")
	defer file.Close()
	if errOnOpen != nil {
		return nil, errOnOpen
	}
	fileInfo, errStat := file.Stat()
	if errStat != nil {
		return nil, errStat
	}
	data := make([]byte, fileInfo.Size())
	_, errRead := file.Read(data)
	if errRead != nil {
		return nil, errRead
	}
	return data, nil
}

// GetConfigFromTravis Get the config inside the travis file and transform it in a format that can be readable
func GetConfigFromTravis() (map[string]interface{}, error) {
	config, errOnLoad := openConfigTravis()
	if errOnLoad != nil {
		return nil, errOnLoad
	}
	confStructured := make(map[string]interface{})
	errOnYml := yaml.Unmarshal(config, &confStructured)
	if errOnYml != nil {
		return nil, errOnYml
	}
	return confStructured, nil
}
