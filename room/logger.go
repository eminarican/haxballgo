package room

import "github.com/rs/zerolog/log"

type Logger struct{}

func (*Logger) Error(msg string) {
	log.Error().Msg(msg)
}

func (*Logger) Errorf(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}

func (*Logger) Warn(msg string) {
	log.Warn().Msg(msg)
}

func (*Logger) Warnf(format string, v ...interface{}) {
	log.Warn().Msgf(format, v)
}

func (*Logger) Info(msg string) {
	log.Info().Msg(msg)
}

func (*Logger) Infof(format string, v ...interface{}) {
	log.Info().Msgf(format, v)
}

func (*Logger) Debug(msg string) {
	log.Debug().Msg(msg)
}

func (*Logger) Debugf(format string, v ...interface{}) {
	log.Debug().Msgf(format, v)
}
