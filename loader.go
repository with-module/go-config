package config

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"path/filepath"
	"reflect"
	"strings"
)

func Load(configObject any, envPrefix string, configFiles ...string) error {
	if pt := reflect.TypeOf(configObject).Kind(); pt != reflect.Ptr {
		return fmt.Errorf("invalid input config object, should be pointer instead of %v", pt)
	}

	inst := koanf.New(".")
	parser := yaml.Parser()
	for _, filename := range configFiles {
		if ext := strings.ToLower(filepath.Ext(filename)); ext != ".yml" && ext != ".yaml" {
			return fmt.Errorf("invalid config file extension, only YAML supported: %s", filename)
		}

		if err := inst.Load(file.Provider(filename), parser); err != nil {
			return fmt.Errorf("failed to load config file %s: %v", filename, err)
		}
	}

	if err := inst.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil); err != nil {
		return fmt.Errorf("failed to load config from environment: %v", err)
	}
	if err := inst.UnmarshalWithConf("", configObject, koanf.UnmarshalConf{Tag: "config"}); err != nil {
		return fmt.Errorf("failed to unmarshal config into object: %v", err)
	}
	return nil
}
