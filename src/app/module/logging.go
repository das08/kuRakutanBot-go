package module

import (
	"github.com/rs/zerolog"
	"io"
	"os"
)

type Logger struct {
	logger *zerolog.Logger
}

func New(isDebug bool) *Logger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &Logger{logger: &logger}
}

func (l *Logger) Output(w io.Writer) zerolog.Logger {
	return l.logger.Output(w)
}

func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}
