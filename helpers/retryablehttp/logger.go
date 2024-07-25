package retryablehttp

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog"
)

type leveledLogger struct {
	zerolog.Logger
}

func (l *leveledLogger) Debug(msg string, keysAndValues ...any) {
	l.Logger.Debug().Fields(keysAndValues).Msg(msg)
}

func (l *leveledLogger) Error(msg string, keysAndValues ...any) {
	l.Logger.Error().Fields(keysAndValues).Msg(msg)
}

func (l *leveledLogger) Info(msg string, keysAndValues ...any) {
	l.Logger.Info().Fields(keysAndValues).Msg(msg)
}

func (l *leveledLogger) Warn(msg string, keysAndValues ...any) {
	l.Logger.Warn().Fields(keysAndValues).Msg(msg)
}

func NewLeveledLogger(logger zerolog.Logger) retryablehttp.LeveledLogger {
	return &leveledLogger{logger}
}
