package travis

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/icza/dyno"
	yaml "gopkg.in/yaml.v2"
)

// Config represents values from a '.travis.yml' file. Config options are represented
// with a slice of strings. If the slice is nil, then the config option was not specified
// in the yml file. If there was only a single value, it will be the only member in the slice.
type Config struct {
	Language         string
	Version          []string
	BeforeInstall    []string
	Install          []string
	BeforeScript     []string
	Script           []string
	Env              []string
	DockerfileConfig string
}

// ConfigFromFile creates a config from a file path to '.travis.yml'
func ConfigFromFile(filepath string) (*Config, error) {
	// TODO: Verify yml extension

	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("Could not read the file: %v", err)
	}

	return configFromContent(b)
}

// configFromContent creates a Config from the bytes content of the yml file. For each field, the values
// are serialized into string slices, with checks at each step to verify that the yml data is valid.
func configFromContent(contents []byte) (*Config, error) {
	conf := make(map[string]interface{})
	err := yaml.Unmarshal(contents, &conf)
	if err != nil {
		return nil, fmt.Errorf("Could not parse yml file: %v", err)
	}

	// Parse all values from the map for usage in the config struct

	beforeInstall, err := getSlice(conf, "before_install")
	if err != nil {
		return nil, fmt.Errorf("Could not parse 'before_install' from the travis config: %v", err)
	}

	install, err := getSlice(conf, "install")
	if err != nil {
		return nil, fmt.Errorf("Could not parse 'install' from the travis config: %v", err)
	}

	beforeScript, err := getSlice(conf, "before_script")
	if err != nil {
		return nil, fmt.Errorf("Could not parse 'before_script' from the travis config: %v", err)
	}

	script, err := getSlice(conf, "script")
	if err != nil {
		return nil, fmt.Errorf("Could not parse 'script' from the travis config: %v", err)
	}

	env, err := getSlice(conf, "env")
	if err != nil {
		return nil, fmt.Errorf("Could not parse 'env' from the travis config: %v", err)
	}

	return &Config{
		BeforeInstall: beforeInstall,
		Install:       install,
		BeforeScript:  beforeScript,
		Script:        script,
		Env:           env, // TODO: Parse this into a map?
	}, nil
}

// getSlice dynamically looks up the "path" from the "from" interface and tries its best to marshal it into a slice of string.
func getSlice(from interface{}, path ...interface{}) ([]string, error) {
	// Check if it is a string value
	if v, err := dyno.GetString(from, path...); err == nil {
		return []string{v}, nil
	}

	// Check if it is a floating value. If it is, convert to string and return.
	if v, err := dyno.GetFloating(from, path...); err == nil {
		strv := strconv.FormatFloat(v, 'f', -1, 64)
		return []string{strv}, nil
	}

	// Check if it is a slice value. If not, we are done and the path doesn't have a value.
	v, err := dyno.GetSlice(from, path...)
	if err != nil {
		return nil, nil
	}

	// For each item in the slice, check if it is a float or a string and handle it accordingly.
	var values []string
	for _, item := range v {
		strItem, ok := item.(string)
		if ok {
			values = append(values, strItem)
			continue
		}

		floatItem, ok := item.(float64)
		if ok {
			stringItem := strconv.FormatFloat(floatItem, 'f', -1, 64)
			values = append(values, stringItem)
			continue
		}

		// We should never get to this point when we have valid string or float data
		return nil, fmt.Errorf("Invalid value: %v found in the config", item)
	}

	return values, nil
}
