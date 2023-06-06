package config

import (
	"os"
	"path"
	"testing"
	"time"
)

type (
	SampleConfig struct {
		App      AppConfig `config:"App"`
		Database struct {
			Host string `config:"Host"`
			Port int    `config:"Port"`
		} `config:"Database"`
	}

	AppConfig struct {
		Server     string        `config:"Server"`
		Timeout    time.Duration `config:"Timeout"`
		BaseURL    string        `config:"BaseUrl"`
		LogEnabled bool          `config:"LogEnabled"`
	}
)

func TestLoadConfigSuccessfully(t *testing.T) {

	var config = new(SampleConfig)
	var configPath = path.Join("samples", "config.yaml")

	const envToOverride = "CFG_App_Timeout"
	const expectedTimeout = time.Second * 90 // 90s
	if err := os.Setenv(envToOverride, expectedTimeout.String()); err != nil {
		t.Errorf("failed to set env variable %s, expect %s, got %s", envToOverride, expectedTimeout, os.Getenv(envToOverride))
	}

	if err := Load(config, "CFG_", configPath); err != nil {
		t.Errorf("expect no error when loading config object: %v", err)
	}

	const expectedBaseURL = "/api/v1"
	if config.App.BaseURL != expectedBaseURL {
		t.Errorf("expect baseURL value to be %s, got %s", expectedBaseURL, config.App.BaseURL)
	}

	if config.App.Timeout != expectedTimeout {
		t.Errorf("expect timeout value to be %s, got %s, env %s", expectedTimeout, config.App.Timeout, os.Getenv(envToOverride))
	}
}
