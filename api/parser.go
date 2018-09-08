package api

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	language      string   `yaml:"language"`
	versions      []string `yaml:"versions"`
	beforeInstall []string `yaml:"before_install"`
	install       []string `yaml:"install"`
	beforeScript  []string `yaml:"before_script"`
	script        []string `yaml:"script"`
}

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

func getConfigFromTravis() error {
	config, errOnLoad := openConfigTravis()
	if errOnLoad != nil {
		return errOnLoad
	}
	confStructured := Config{}
	errOnYml := yaml.Unmarshal(config, &confStructured)
	if errOnYml != nil {
		return errOnYml
	}
	return nil
}
