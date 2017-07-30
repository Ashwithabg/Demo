package config

import (
	"testing"
	"path"
	"os"
	"io/ioutil"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

const testConfig = `{
"baseURL": "/github/api"
}
`

func TestLoadConfig(t *testing.T) {
	filePath := path.Join(os.TempDir(), "ConfigurationForTest.json")
	err := ioutil.WriteFile(filePath, []byte(testConfig), 0544)
	require.NoError(t, err)
	defer os.Remove(filePath)

	config := LoadConfig(filePath)

	assert.Equal(t, "", config.Errors(), "Expected config to have no errors")

	assert.Equal(t, "/github/api", config.BaseURL())
}

func TestLoadConfigFails(t *testing.T) {
	filePath := path.Join(os.TempDir(), "ConfigurationForTest.json")
	err := ioutil.WriteFile(filePath, []byte("invalid"), 0544)
	require.NoError(t, err)
	defer os.Remove(filePath)

	config := LoadConfig(filePath)

	assert.Contains(t, config.Errors(), "invalid configuration: ", "Expected config to have no errors")
}

func TestLoadConfigWithNonExistingFilePathHasError(t *testing.T) {
	config := LoadConfig("InvalidFilePath")

	assert.True(t, config.HasErrors(), "Expected config to have errors")
	assert.Contains(t, config.Errors(), "invalid configuration: failed to open configuration file 'InvalidFilePath': open InvalidFilePath: ")
}
