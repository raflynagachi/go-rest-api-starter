package logger

type Options struct {
	Level     Level
	AddSource bool
	Env       string
}

type Option func(*Options)

// WithLevel sets the log level. The default level is Info.
func WithLevel(level Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

// WithAddSource sets the add source option.
func WithAddSource(addSource bool) Option {
	return func(o *Options) {
		o.AddSource = addSource
	}
}

// WithEnv sets the system environment
func WithEnv(env string) Option {
	return func(o *Options) {
		o.Env = env
	}
}
