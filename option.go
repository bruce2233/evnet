package evnet

import (
	log "github.com/sirupsen/logrus"
)

type SetOption func(optArg *OptionArg)

type OptionArg struct {
	LogPath  string
	LogLevel log.Level
}

func WithLogPath(logPath string) SetOption {
	return func(optArg *OptionArg) {
		optArg.LogPath = logPath
	}
}

func WithLogLevel(logLevel log.Level) SetOption {
	return func(optArg *OptionArg) {
		optArg.LogLevel = logLevel
	}
}

func loadOptions(optList []SetOption) *OptionArg {
	optArg := new(OptionArg)
	for _, opt := range optList {
		opt(optArg)
	}
	return optArg
}
