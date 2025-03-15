package log

type Logger interface {
	Debug(msg string, kv ...interface{})
	Info(msg string, kv ...interface{})
	Warn(msg string, kv ...interface{})
	Error(msg string, err error, kv ...interface{})
	Fatal(msg string, kv ...interface{})
}
