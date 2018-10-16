package api

import (
	"fmt"
	"os"
	"reflect"

	"github.com/srishanbhattarai/previs/travis"
)

// API handles all the Docker related tasks as defined by the config.
type API struct {
	config *travis.Config
}

// FromConfig creates a new instance of the API from the provided config file.
func FromConfig(configFile string) (*API, error) {
	config, err := travis.ConfigFromFile(configFile)
	if err != nil {
		return nil, err
	}

	return &API{config}, nil
}

// WriteDockerfile writes the config from travis to a new Dockerfile.
func (api *API) WriteDockerfile() (string, []string, error) {
	fmt.Println("Creating Dockerfile")

	file, imgDocker, err := createDockerFile()
	if err != nil {
		return "", nil, err
	}

	api.writerFrom()
	api.writerAddConfig()
	api.writerRunBeforeInstall()
	api.writerRunInstall()
	api.writerRunBeforeScript()
	api.writerRunScript()
	envs := api.getEnvsVariables()
	file.WriteString(api.config.DockerfileConfig)

	err = file.Close()
	if err != nil {
		return "", nil, err
	}

	return imgDocker, envs, nil
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

func (api *API) writerFrom() {
	var from string
	images := map[string]string{
		"node_js": "node",
		"go":      "golang",
	}

	if api.config.Version == nil {
		fmt.Fprintf(os.Stderr, "The <from> and <version> directive is mandatory in your travis config")
		os.Exit(2)
	}

	version := api.config.Version[0]

	if platform, ok := images[api.config.Language]; ok {
		from = fmt.Sprintf("FROM %+v:%+v\n", platform, version)
	} else {
		from = fmt.Sprintf("FROM %+v:%+v\n", api.config.Language, version)
	}

	api.config.DockerfileConfig = api.config.DockerfileConfig + from
}

func (api *API) writerAddConfig() {
	api.config.DockerfileConfig = api.config.DockerfileConfig + "RUN apt-get update\nWORKDIR /home/app/\nCOPY ./ /home/app\n"
}

func (api *API) writerRunBeforeInstall() {
	if api.config.BeforeInstall != nil {
		var runBeforeInstall string
		for i := 0; i < len(api.config.BeforeInstall); i++ {
			runBeforeInstall = runBeforeInstall + fmt.Sprintf("RUN %+v\n", api.config.BeforeInstall[i])
		}

		api.config.DockerfileConfig = api.config.DockerfileConfig + runBeforeInstall
	}
}

func (api *API) writerRunInstall() {
	if api.config.Install != nil {
		var runInstall string
		for i := 0; i < len(api.config.Install); i++ {
			runInstall = runInstall + fmt.Sprintf("RUN %+v\n", api.config.Install[i])
		}

		api.config.DockerfileConfig = api.config.DockerfileConfig + runInstall
	}
}

func (api *API) writerRunBeforeScript() {
	if api.config.BeforeScript != nil {
		var runBeforeScript string
		for i := 0; i < len(api.config.BeforeScript); i++ {
			runBeforeScript = runBeforeScript + fmt.Sprintf("RUN %+v\n", api.config.BeforeScript[i])
		}

		api.config.DockerfileConfig = api.config.DockerfileConfig + runBeforeScript
	}
}

func (api *API) writerRunScript() {
	if api.config.Script == nil {
		fmt.Fprintln(os.Stderr, "The <script> directive is mandatory in your travis config")
		os.Exit(2)
	}

	entrypoint := api.config.Script[0]
	cmd := fmt.Sprintf("CMD %+v\n", entrypoint)
	api.config.DockerfileConfig = api.config.DockerfileConfig + cmd
}

func (api *API) getEnvsVariables() []string {
	var envs []string
	if api.config.Env != nil {
		for i := 0; i < len(api.config.Env); i++ {
			envs = append(envs, fmt.Sprintf("%+v", api.config.Env[i]))
		}
	}
	return envs
}
