package evnet

type Option func(opts *Options)

type Options struct {
	LogPath string
}

func WithLogPath(logPath string) Option {
	return func(opts *Options) {
		opts.LogPath = logPath
	}
}

func loadOptions(optList []Option) *Options {
	opts := new(Options)
	for _, opt := range optList {
		opt(opts)
	}
	return opts
}
