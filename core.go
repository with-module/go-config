package config

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/mitchellh/mapstructure"
	"path/filepath"
	"slices"
	"strings"
)

func Parse[T any](filename string, opts ...Option) (*T, error) {
	sts := &settings{tag: "config", env: envSettings{enabled: false, sensitiveCase: false}}
	for _, h := range opts {
		h.use(sts)
	}

	inst := koanf.New(".")
	parser := yaml.Parser()
	if ext := strings.ToLower(filepath.Ext(filename)); ext != ".yml" && ext != ".yaml" {
		return nil, fmt.Errorf("invalid config file extension, only YAML supported: %s", filename)
	}

	if err := inst.Load(file.Provider(filename), parser); err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", filename, err)
	}

	if conf := sts.env; conf.enabled {
		callbackFunc := caseSensitiveEnvCallbackFunc(conf.prefix)
		if !conf.sensitiveCase {
			callbackFunc = caseInSensitiveEnvCallbackFunc(inst.Keys(), sts.env.prefix)
		}
		if err := inst.Load(env.Provider(sts.env.prefix, ".", callbackFunc), nil); err != nil {
			return nil, fmt.Errorf("failed to load config from environment: %w", err)
		}
	}
	objectData := new(T)
	decoder := &mapstructure.DecoderConfig{
		DecodeHook:           mapstructure.ComposeDecodeHookFunc(mapstructure.StringToTimeDurationHookFunc(), mapstructure.TextUnmarshallerHookFunc()),
		WeaklyTypedInput:     true,
		Squash:               sts.useSquash,
		Metadata:             nil,
		Result:               objectData,
		TagName:              sts.tag,
		IgnoreUntaggedFields: sts.omitUntaggedFields,
		MatchName:            strings.EqualFold,
	}
	if err := inst.UnmarshalWithConf("", objectData, koanf.UnmarshalConf{Tag: sts.tag, DecoderConfig: decoder}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into object: %v", err)
	}
	return objectData, nil
}

func caseSensitiveEnvCallbackFunc(prefix string) func(s string) string {
	return func(s string) string {
		return strings.Replace(strings.TrimPrefix(s, prefix), "_", ".", -1)
	}
}

func caseInSensitiveEnvCallbackFunc(keys []string, prefix string) func(s string) string {
	return func(s string) string {
		target := caseSensitiveEnvCallbackFunc(prefix)(s)
		idx := slices.IndexFunc(keys, func(it string) bool {
			return strings.EqualFold(it, target)
		})
		if idx >= 0 {
			return keys[idx]
		}

		return target
	}
}
