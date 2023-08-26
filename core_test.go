package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
	"time"
)

type (
	SampleConfig struct {
		App      AppConfig `config:"app"`
		Database struct {
			Host string `config:"host"`
			Port int    `config:"port"`
		} `config:"database"`
	}

	AppConfig struct {
		Server  string        `config:"server"`
		Timeout time.Duration `config:"timeout"`
		URL     string        `config:"url"`
		Log     bool          `config:"log"`
	}
)

func TestParseConfigSuccessfully(t *testing.T) {
	var configPath = path.Join("samples", "config.yaml")

	const envToOverride = "CFG_APP_TIMEOUT"
	const expectedTimeout = time.Second * 90 // 90s
	assert.NoError(t, os.Setenv(envToOverride, expectedTimeout.String()))
	assert.NoError(t, os.Setenv("CFG_SOME_UNKNOWN", "undefined"))

	data, err := Parse[SampleConfig](configPath, UseTag("config"), UseEnvPrefix("CFG_"))
	assert.NoError(t, err)
	assert.IsType(t, new(SampleConfig), data)
	assert.Equal(t, "/api/v1", data.App.URL)
	assert.Equal(t, expectedTimeout, data.App.Timeout)
}

func TestParseErrorDueToUnsupportedFile(t *testing.T) {
	var configPath = path.Join("samples", "config.invalidExt")
	data, err := Parse[SampleConfig](configPath, UseSquash(), UseUntagOmit())
	assert.ErrorContains(t, err, "invalid config file extension")
	assert.Nil(t, data)
}

func TestParseErrorDueToInvalidFilePath(t *testing.T) {
	var configPath = path.Join("samples", "not-exist", "config.yaml")
	data, err := Parse[SampleConfig](configPath, UseCaseSensitiveMode(false))
	assert.ErrorContains(t, err, "failed to load config file")
	assert.Nil(t, data)
}

type InvalidStruct struct {
	App int `config:"app"`
}

func TestParseErrorInvalidStruct(t *testing.T) {
	var configPath = path.Join("samples", "config.yaml")

	data, err := Parse[InvalidStruct](configPath)
	assert.ErrorContains(t, err, "failed to unmarshal config into object")
	assert.Nil(t, data)
}
