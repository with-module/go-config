package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	AppConfig struct {
		Server struct {
			Host string `config:"host"`
			Port int    `config:"port"`
		} `config:"server"`
		Env string `config:"env"`
	}
)

var (
	configFiles = []string{"./mock/app_config.yml"}
)

func TestLoadConfig(t *testing.T) {
	asst := assert.New(t)
	var sampleConfig = AppConfig{}
	err := LoadConfig(&sampleConfig, "APP_", configFiles...)
	asst.NoError(err, "should not have err")
	asst.EqualValues("test-yml", sampleConfig.Env, "value of env should be %s", "test-yml")
}
