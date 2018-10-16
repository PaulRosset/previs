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
	env              reflect.Value
	dockerfileConfig string
}

func createDockerFile() (*os.File, string, error) {
	imgDocker := "Previsfile"
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

func (c *Config) writerFrom() {
	var from string
	images := map[string]string{
		"node_js": "node",
		"go":      "golang",
	}
	if c.version.IsValid() {
		firstVersion := c.version.Index(0)
		if platform, ok := images[c.platform]; ok {
			from = fmt.Sprintf("FROM %+v:%+v\n", platform, firstVersion)
		} else {
			from = fmt.Sprintf("FROM %+v:%+v\n", c.platform, firstVersion)
		}
		c.dockerfileConfig = c.dockerfileConfig + from
	} else {
		fmt.Fprintf(os.Stderr, "The <from> and <version> directive is mandatory in your travis config")
		os.Exit(2)
	}
}

func (c *Config) writerAddConfig() {
	c.dockerfileConfig = c.dockerfileConfig + "RUN apt-get update\nWORKDIR /home/app/\nCOPY ./ /home/app\n"
}

func (c *Config) writerRunBeforeInstall() {
	if c.beforeInstall.IsValid() {
		var runBeforeInstall string
		for i := 0; i < c.beforeInstall.Len(); i++ {
			runBeforeInstall = runBeforeInstall + fmt.Sprintf("RUN %+v\n", c.beforeInstall.Index(i))
		}
		c.dockerfileConfig = c.dockerfileConfig + runBeforeInstall
	}
}

func (c *Config) writerRunInstall() {
	if c.install.IsValid() {
		var runInstall string
		for i := 0; i < c.install.Len(); i++ {
			runInstall = runInstall + fmt.Sprintf("RUN %+v\n", c.install.Index(i))
		}
		c.dockerfileConfig = c.dockerfileConfig + runInstall
	}
}

func (c *Config) writerRunBeforeScript() {
	if c.beforeScript.IsValid() {
		var runBeforeScript string
		for i := 0; i < c.beforeScript.Len(); i++ {
			runBeforeScript = runBeforeScript + fmt.Sprintf("RUN %+v\n", c.beforeScript.Index(i))
		}
		c.dockerfileConfig = c.dockerfileConfig + runBeforeScript
	}
}

func (c *Config) writerRunScript() {
	if c.script.IsValid() {
		entrypoint := c.script.Index(0)
		cmd := fmt.Sprintf("CMD %+v\n", entrypoint)
		c.dockerfileConfig = c.dockerfileConfig + cmd
	} else {
		fmt.Fprintln(os.Stderr, "The <script> directive is mandatory in your travis config")
		os.Exit(2)
	}
}

func (c *Config) getEnvsVariables() []string {
	var envs []string
	if c.env.IsValid() {
		for i := 0; i < c.env.Len(); i++ {
			envs = append(envs, fmt.Sprintf("%+v", c.env.Index(i)))
		}
	}
	return envs
}

// writer is writting the config from travis to a new a dockerfile
func writer(configFile string) (string, []string, error) {
	config, err := GetConfigFromTravis(configFile)
	if err != nil {
		return "", nil, err
	}
	exploitConfig := &Config{
		platform:         config["language"].(string),
		version:          reflectInterface(config[config["language"].(string)]),
		beforeInstall:    reflectInterface(config["before_install"]),
		install:          reflectInterface(config["install"]),
		beforeScript:     reflectInterface(config["before_script"]),
		script:           reflectInterface(config["script"]),
		env:              reflectInterface(config["env"]),
		dockerfileConfig: "",
	}
	file, imgDocker, err := createDockerFile()
	if err != nil {
		return "", nil, err
	}
	exploitConfig.writerFrom()
	exploitConfig.writerAddConfig()
	exploitConfig.writerRunBeforeInstall()
	exploitConfig.writerRunInstall()
	exploitConfig.writerRunBeforeScript()
	exploitConfig.writerRunScript()
	envs := exploitConfig.getEnvsVariables()
	file.WriteString(exploitConfig.dockerfileConfig)
	err = file.Close()
	if err != nil {
		return "", nil, err
	}
	return imgDocker, envs, nil
}
