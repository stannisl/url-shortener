package logger

import "github.com/wb-go/wbf/zlog"

type ZLogger struct {
	l zlog.Zerolog
}

func NewZLogger(l zlog.Zerolog) *ZLogger {
	return &ZLogger{l: l}
}

func (z *ZLogger) Info(msg string) {
	z.l.Info().Msg(msg)
}

func (z *ZLogger) Error(msg string) {
	z.l.Error().Msg(msg)
}

func (z *ZLogger) ErrorErr(err error, msg string) {
	z.l.Error().Err(err).Msg(msg)
}

func (z *ZLogger) FatalErr(err error, msg string) {
	z.l.Fatal().Err(err).Msg(msg)
}
