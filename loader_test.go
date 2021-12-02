package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
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
	configFiles = []string{
		"./mock/app_config.yml",
		"./mock/app_config.json",
	}
)

func TestLoadConfigSuccessCase(t *testing.T) {
	asst := assert.New(t)
	var sampleConfig = AppConfig{}
	err := LoadConfig(&sampleConfig, "APP_", configFiles...)
	asst.NoError(err, "should not have err")
	asst.EqualValues("test-json", sampleConfig.Env, "value of env should be %s", "test-json")
}

func TestLoadConfigFailedOfConfigObject(t *testing.T) {
	asst := assert.New(t)
	var sampleConfig = AppConfig{}
	expectedErr := fmt.Errorf("invalid input config object, should be pointer instead of %v", reflect.TypeOf(sampleConfig).Kind())
	err := LoadConfig(sampleConfig, "APP_", configFiles...)
	asst.EqualError(err, expectedErr.Error(), "error because of not pointer object")
}
