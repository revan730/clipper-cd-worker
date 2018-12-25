package log

type Logger interface {
	LogInfo(i string)
	LogError(msg string, err error)
	LogFatal(msg string, err error)
}
