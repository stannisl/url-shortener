package logger

type Logger interface {
	Info(msg string)
	Error(msg string)
	ErrorErr(err error, msg string)
	FatalErr(err error, msg string)
}
