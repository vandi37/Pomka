package logger

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type Options struct {
	LogLevel        logrus.Level
	TimestampFormat string
	FullTimestamp   bool
	DisableSorting  bool
	ColorScheme     *prefixed.ColorScheme
}

var defaultOptions = Options{

	LogLevel:        logrus.DebugLevel,
	TimestampFormat: "2006-01-02 15:04:05",
	FullTimestamp:   true,
	ColorScheme: &prefixed.ColorScheme{
		InfoLevelStyle:  "green+b",
		WarnLevelStyle:  "yellow+b",
		ErrorLevelStyle: "red+b",
		FatalLevelStyle: "red+b",
		PanicLevelStyle: "red+b",
		DebugLevelStyle: "blue+b",
		PrefixStyle:     "cyan+b",
		TimestampStyle:  "black+h",
	},
}
