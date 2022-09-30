package clients

import "github.com/rs/zerolog"

func jsonToLog(l zerolog.Logger, msg map[string]interface{}) {
	switch msg["level"] {
	case "trace":
		l.Trace().Fields(msg).Msg("")
	case "debug":
		l.Debug().Fields(msg).Msg("")
	case "info":
		l.Debug().Fields(msg).Msg("")
	case "warn":
		l.Warn().Fields(msg).Msg("")
	case "error":
		l.Warn().Fields(msg).Msg("")
	default:
		l.Error().Fields(msg).Msg("unknown level")
	}
}
