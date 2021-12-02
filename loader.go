package config

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"path/filepath"
	"reflect"
	"strings"
)

func LoadConfig(configObject interface{}, envPrefix string, configFiles ...string) error {
	if pt := reflect.TypeOf(configObject).Kind(); pt != reflect.Ptr {
		return fmt.Errorf("invalid input config object, should be pointer instead of %v", pt)
	}

	loader := koanf.New(".")
	for _, filename := range configFiles {
		var parser koanf.Parser
		switch ext := strings.ToLower(filepath.Ext(filename)); ext {
		case ".yml", ".yaml":
			parser = yaml.Parser()
		case ".json":
			parser = json.Parser()
		default:
			return fmt.Errorf("invalid config file: %s", filename)
		}

		if err := loader.Load(file.Provider(filename), parser); err != nil {
			return fmt.Errorf("failed to load config file %s: %v", filename, err)
		}
	}

	if err := loader.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil); err != nil {
		return fmt.Errorf("failed to load config from environment: %v", err)
	}
	if err := loader.UnmarshalWithConf("", configObject, koanf.UnmarshalConf{Tag: "config"}); err != nil {
		return fmt.Errorf("failed to unmarshal config into object: %v", err)
	}
	return nil
}
