package api

import (
	"fmt"
	"os"
	"reflect"

	"github.com/PaulRosset/previs/travis"
)

const pathGarnet = "/home/travis/.local/bin:/opt/pyenv/shims:/home/travis/.phpenv/shims:/home/travis/perl5/perlbrew/bin:/home/travis/.nvm/versions/node/v8.9.1/bin:/home/travis/.kiex/elixirs/elixir-1.4.5/bin:/home/travis/.kiex/bin:/home/travis/.rvm/gems/ruby-2.4.1/bin:/home/travis/.rvm/gems/ruby-2.4.1@global/bin:/home/travis/.rvm/rubies/ruby-2.4.1/bin:/home/travis/.phpenv/shims:/home/travis/gopath/bin:/home/travis/.gimme/versions/go1.7.4.linux.amd64/bin:/usr/local/phantomjs/bin:/usr/local/phantomjs:/usr/local/neo4j-3.2.7/bin:/usr/local/maven-3.5.2/bin:/usr/local/cmake-3.9.2/bin:/usr/local/clang-5.0.0/bin:/home/travis/.gimme/versions/go1.7.4.linux.amd64/bin:/usr/local/phantomjs/bin:/usr/local/phantomjs:/usr/local/neo4j-3.2.7/bin:/usr/local/maven-3.5.2/bin:/usr/local/cmake-3.9.2/bin:/usr/local/clang-5.0.0/bin:/usr/local/bin:/usr/bin:/bin:/usr/local/games:/usr/games:/home/travis/.rvm/bin:/home/travis/.phpenv/bin:/opt/pyenv/bin"

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
func (api *API) WriteDockerfile() (string, []string, []string, error) {
	fmt.Println("Creating Dockerfile")

	file, imgDocker, err := createDockerFile()
	if err != nil {
		return "", nil, nil, err
	}

	api.writerFrom()
	api.writerAddConfig()
	api.writerRunBeforeInstall()
	api.writerRunInstall()
	api.writerRunBeforeScript()
	api.writerRunScript()
	envs := api.getEnvsVariables()
	services := api.getServices()
	file.WriteString(api.config.DockerfileConfig)

	err = file.Close()
	if err != nil {
		return "", nil, nil, err
	}

	return imgDocker, envs, services, nil
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
	from := "FROM travisci/ci-garnet:packer-1515445631-7dfb2e1\n"
	api.config.DockerfileConfig = api.config.DockerfileConfig + from
}

func (api *API) writerAddConfig() {
	api.config.DockerfileConfig = api.config.DockerfileConfig + "USER travis\nENV PATH " + pathGarnet + "\nWORKDIR /home/travis/builds/\nCOPY ./ /home/travis/builds/\n"
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

func (api *API) getServices() []string {
	var services []string
	if api.config.Services != nil {
		for i := 0; i < len(api.config.Services); i++ {
			services = append(services, fmt.Sprintf("%+v", api.config.Services[i]))
		}
	}
	return services
}
