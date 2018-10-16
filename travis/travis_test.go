package travis

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigFromFile(t *testing.T) {
	tt := []struct {
		name        string
		description string
		filepath    string
		err         error
		config      *Config
	}{
		{
			name:     "File is not present",
			filepath: "./testdata/foobar.yml",
			err:      fmt.Errorf("Could not read the file"),
		},
		{
			name:        "Valid file is present - case 1",
			filepath:    "./testdata/.travis.yml",
			description: "This covers the sane case. Values for go version are floats, and slices are used.",
			config: &Config{
				Language: "go",
				Version:  []string{"1.9", "1.8"},
				Env:      []string{"DB=1234"},
				Script:   []string{"make test"},
			},
		},
		{
			name:        "Valid file is present - case 2",
			description: "This covers some edge cases. 'env' key is present but doesn't have a value. Values for version is string.",
			filepath:    "./testdata/.travis.2.yml",
			config: &Config{
				Language: "go",
				Version:  []string{"1.9"},
				Script:   []string{"make test"},
			},
		},
		{
			name:     "File is present but invalid",
			filepath: "./testdata/invalidYml.yml",
			err:      fmt.Errorf("Could not parse yml file"),
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			c, err := ConfigFromFile(test.filepath)
			if err != nil {
				assert.Contains(t, err.Error(), test.err.Error())
			} else {
				assert.Equal(t, test.config, c)
			}
		})
	}
}
