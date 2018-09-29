package api

import (
	"fmt"
	"os"
	"reflect"
)

// Config regroup the travis config
type Config struct {
	platform         string
	version          reflect.Value
	beforeInstall    reflect.Value
	install          reflect.Value
	beforeScript     reflect.Value
	script           reflect.Value
	afterScript      reflect.Value
	dockerfileConfig string
}

func createDockerFile() (*os.File, string, error) {
	imgDocker := "Dockerfile"
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
	if c.version.IsValid() {
		firstVersion := c.version.Index(0)
		from := fmt.Sprintf("FROM treevis/%+v:%+v\n", c.platform, firstVersion)
		c.dockerfileConfig = c.dockerfileConfig + from
	} else {
		fmt.Fprintf(os.Stderr, "The <from> and <version> directive is mandatory in your travis config")
		os.Exit(2)
	}
}

func (c *Config) writterAddConfig() {
	c.dockerfileConfig = c.dockerfileConfig + "COPY ./ /home/app\n"
}

func (c *Config) writterRunBeforeInstall() {
	if c.beforeInstall.IsValid() {
		var runBeforeInstall string
		for i := 0; i < c.beforeInstall.Len(); i++ {
			runBeforeInstall = runBeforeInstall + fmt.Sprintf("RUN %+v\n", c.beforeInstall.Index(i))
		}
		c.dockerfileConfig = c.dockerfileConfig + runBeforeInstall
	}
}

func (c *Config) writterRunInstall() {
	if c.install.IsValid() {
		var runInstall string
		for i := 0; i < c.install.Len(); i++ {
			runInstall = runInstall + fmt.Sprintf("RUN %+v\n", c.install.Index(i))
		}
		c.dockerfileConfig = c.dockerfileConfig + runInstall
	}
}

func (c *Config) writterRunBeforeScript() {
	if c.beforeScript.IsValid() {
		var runBeforeScript string
		for i := 0; i < c.beforeScript.Len(); i++ {
			runBeforeScript = runBeforeScript + fmt.Sprintf("RUN %+v\n", c.beforeScript.Index(i))
		}
		c.dockerfileConfig = c.dockerfileConfig + runBeforeScript
	}
}

func (c *Config) writterRunScript() {
	if c.script.IsValid() {
		entrypoint := c.script.Index(0)
		cmd := fmt.Sprintf("CMD %+v\n", entrypoint)
		c.dockerfileConfig = c.dockerfileConfig + cmd
	} else {
		fmt.Fprintln(os.Stderr, "The <script> directive is mandatory in your travis config")
		os.Exit(2)
	}
}

func (c *Config) writterRunAfterScript() {
	if c.afterScript.IsValid() {
		var runAfterScript string
		for i := 0; i < c.afterScript.Len(); i++ {
			runAfterScript = runAfterScript + fmt.Sprintf("RUN %+v\n", c.afterScript.Index(i))
		}
		c.dockerfileConfig = c.dockerfileConfig + runAfterScript
	}
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
		afterScript:      reflectInterface(config["after_script"]),
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
	exploitConfig.writterRunAfterScript()
	file.WriteString(exploitConfig.dockerfileConfig)
	return imgDocker, nil
}
