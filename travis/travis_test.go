package travis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigFromFile(t *testing.T) {
	c, err := ConfigFromFile("./testdata/.travis.yml")

	assert.NoError(t, err)
	assert.Equal(t, c.Language, "go")
	assert.ElementsMatch(t, c.Version, []string{"1.9", "1.8"})
	assert.ElementsMatch(t, c.Install, []string{})
	assert.Equal(t, c.Env, []string{"DB=1234"})
}
