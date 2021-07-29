package room

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type Logger struct{}

func (*Logger) Error(msg string) {
	log.Error().Msg(msg)
}

func (*Logger) Errorf(format string, v ...interface{}) {
	log.Error().Msg(fmt.Sprintf(format, v...))
}

func (*Logger) Warn(msg string) {
	log.Warn().Msg(msg)
}

func (*Logger) Warnf(format string, v ...interface{}) {
	log.Warn().Msg(fmt.Sprintf(format, v...))
}

func (*Logger) Info(msg string) {
	log.Info().Msg(msg)
}

func (*Logger) Infof(format string, v ...interface{}) {
	log.Info().Msg(fmt.Sprintf(format, v...))
}

func (*Logger) Debug(msg string) {
	log.Debug().Msg(msg)
}

func (*Logger) Debugf(format string, v ...interface{}) {
	log.Debug().Msg(fmt.Sprintf(format, v...))
}
