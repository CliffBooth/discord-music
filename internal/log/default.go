package log

import (
	"log"
	"sync"

	"go.uber.org/zap"
)

var (
	defaultLogger Logger
	once          = sync.Once{}
)

func Default() Logger {
	once.Do(func() {
		l, err := zap.NewDevelopment()
		if err != nil {
			log.Fatalf("error creating logger: %v\n", err)
		}
		defaultLogger = &ZapLogger{
			logger: l,
		}
	})
	return defaultLogger
}
