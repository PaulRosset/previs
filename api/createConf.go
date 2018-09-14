package api

import (
	"fmt"
	"os"
	"reflect"

	"github.com/satori/go.uuid"
)

// Config regroup the travis config
type Config struct {
	platform         string
	version          reflect.Value
	beforeInstall    reflect.Value
	install          reflect.Value
	beforeScript     reflect.Value
	script           reflect.Value
	dockerfileConfig string
}

func createDockerFile() (*os.File, string, error) {
	uid := uuid.Must(uuid.NewV4())
	imgDocker := ".dockerfile-" + uid.String()
	cwd, errOnCwd := os.Getwd()
	if errOnCwd != nil {
		return nil, "", errOnCwd
	}
	file, errOnOpen := os.OpenFile(cwd+"/"+imgDocker, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errOnOpen != nil {
		return nil, "", errOnOpen
	}
	return file, imgDocker, nil
}

func reflectInterface(t interface{}) reflect.Value {
	s := reflect.ValueOf(t)
	return s
}

func (c *Config) writterFrom() {
	firstVersion := c.version.Index(0)
	from := fmt.Sprintf("FROM treevis/%+v:%+v\n", c.platform, firstVersion)
	c.dockerfileConfig = c.dockerfileConfig + from
}

func (c *Config) writterAddConfig() {
	c.dockerfileConfig = c.dockerfileConfig + "COPY ./* ./\n"
}

func (c *Config) writterRunBeforeInstall() {
	var runBeforeInstall string
	for i := 0; i < c.beforeInstall.Len(); i++ {
		runBeforeInstall = runBeforeInstall + fmt.Sprintf("RUN %+v\n", c.beforeInstall.Index(i))
	}
	c.dockerfileConfig = c.dockerfileConfig + runBeforeInstall
}

func (c *Config) writterRunInstall() {
	var runInstall string
	for i := 0; i < c.install.Len(); i++ {
		runInstall = runInstall + fmt.Sprintf("RUN %+v\n", c.install.Index(i))
	}
	c.dockerfileConfig = c.dockerfileConfig + runInstall
}

func (c *Config) writterRunBeforeScript() {
	var runBeforeScript string
	for i := 0; i < c.beforeScript.Len(); i++ {
		runBeforeScript = runBeforeScript + fmt.Sprintf("RUN %+v\n", c.beforeScript.Index(i))
	}
	c.dockerfileConfig = c.dockerfileConfig + runBeforeScript
}

func (c *Config) writterRunScript() {
	entrypoint := c.script.Index(0)
	cmd := fmt.Sprintf("CMD %+v", entrypoint)
	c.dockerfileConfig = c.dockerfileConfig + cmd
}

// Writter is writting the config from travis to a new a dockerfile
func Writter() (string, error) {
	config, errOnGetConfigTravis := GetConfigFromTravis()
	if errOnGetConfigTravis != nil {
		return "", errOnGetConfigTravis
	}
	exploitConfig := &Config{
		platform:         config["language"].(string),
		version:          reflectInterface(config[config["language"].(string)]),
		beforeInstall:    reflectInterface(config["before_install"]),
		install:          reflectInterface(config["install"]),
		beforeScript:     reflectInterface(config["before_script"]),
		script:           reflectInterface(config["script"]),
		dockerfileConfig: "",
	}
	file, imgDocker, errOnCreationDockerfile := createDockerFile()
	if errOnCreationDockerfile != nil {
		return "", errOnCreationDockerfile
	}
	exploitConfig.writterFrom()
	exploitConfig.writterAddConfig()
	exploitConfig.writterRunBeforeInstall()
	exploitConfig.writterRunInstall()
	exploitConfig.writterRunBeforeScript()
	exploitConfig.writterRunScript()
	fmt.Println(exploitConfig.dockerfileConfig)
	file.WriteString(exploitConfig.dockerfileConfig)
	return imgDocker, nil
}
