package config

type (
	settings struct {
		tag                string
		env                envSettings
		useSquash          bool
		omitUntaggedFields bool
	}

	envSettings struct {
		enabled       bool
		prefix        string
		sensitiveCase bool
	}

	Option interface {
		use(s *settings)
	}

	useOption func(s *settings)
)

func (fn useOption) use(s *settings) {
	fn(s)
}

func UseTag(tag string) Option {
	return useOption(func(s *settings) {
		s.tag = tag
	})
}

func UseCaseSensitiveMode(enabled bool) Option {
	return useOption(func(s *settings) {
		s.env.enabled = true
		s.env.sensitiveCase = enabled
	})
}

func UseEnvPrefix(envPrefix string) Option {
	return useOption(func(s *settings) {
		s.env.prefix = envPrefix
		s.env.enabled = true
	})
}

func UseSquash() Option {
	return useOption(func(s *settings) {
		s.useSquash = true
	})
}

func UseUntagOmit() Option {
	return useOption(func(s *settings) {
		s.omitUntaggedFields = true
	})
}
