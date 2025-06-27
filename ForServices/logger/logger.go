package logger

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger(options ...Options) *Logger {
	var opt Options = defaultOptions

	if len(options) == 1 {
		opt = options[0]
	}

	logger := logrus.New()
	logger.SetLevel(opt.LogLevel)

	formatter := &prefixed.TextFormatter{
		TimestampFormat: opt.TimestampFormat,
		FullTimestamp:   opt.FullTimestamp,
		DisableSorting:  opt.DisableSorting,
	}

	formatter.SetColorScheme(opt.ColorScheme)

	logger.SetFormatter(formatter)
	return &Logger{logger}
}
