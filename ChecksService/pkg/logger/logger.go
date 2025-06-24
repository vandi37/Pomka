package logger

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type Options struct {
	Level           logrus.Level
	TimestampFormat string
	FullTimestamp   bool
	DisableSorting  bool
}

func NewLogger(options ...*Options) *logrus.Logger {
	var opt *Options

	if len(options) == 0 || len(options) > 1 {
		opt = &Options{
			Level:           logrus.DebugLevel,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			DisableSorting:  true,
		}
	} else {
		opt = options[0]
	}

	logger := logrus.New()
	logger.SetLevel(opt.Level)

	formatter := &prefixed.TextFormatter{
		TimestampFormat: opt.TimestampFormat,
		FullTimestamp:   opt.FullTimestamp,
		DisableSorting:  opt.DisableSorting,
	}

	formatter.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle:  "green+b",
		WarnLevelStyle:  "yellow+b",
		ErrorLevelStyle: "red+b",
		FatalLevelStyle: "red+b",
		PanicLevelStyle: "red+b",
		DebugLevelStyle: "blue+b",
		PrefixStyle:     "cyan+b",
		TimestampStyle:  "black+h",
	})

	logger.SetFormatter(formatter)
	return logger
}
