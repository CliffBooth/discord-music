package log

// если мы хотим иметь возможность менять либы для лога, нам нужен единый интерфейс логгера
// в этот пакет можно добавлять реализации этого интерфейса через разные либы

type Logger interface {
	Debug(args ...interface{})
	Debugf(pattern string, args ...interface{})

	Info(args ...interface{})
	Infof(pattern string, args ...interface{})

	Warn(args ...interface{})
	Warnf(pattern string, args ...interface{})

	Error(args ...interface{})
	Errorf(pattern string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(pattern string, args ...interface{})
}
