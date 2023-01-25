package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Println(v ...interface{})
	Debug(msg string, v ...interface{})
	Info(msg string, v ...interface{})
	Warn(msg string, v ...interface{})
	Error(msg string, v ...interface{})
}

func NewLogger() Logger {
	fancyLogger := logrus.New()
	fancyLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return appLogger{logger: fancyLogger}
}

type appLogger struct {
	logger *logrus.Logger
}

func (l appLogger) Println(v ...interface{}) {
	l.logger.WithFields(l.toFields(v...))
}

func (l appLogger) Debug(msg string, v ...interface{}) {
	l.logger.WithFields(l.toFields(v...)).Debug(msg)
}

func (l appLogger) Info(msg string, v ...interface{}) {
	l.logger.WithFields(l.toFields(v...)).Info(msg)
}

func (l appLogger) Warn(msg string, v ...interface{}) {
	l.logger.WithFields(l.toFields(v...)).Warn(msg)
}

func (l appLogger) Error(msg string, v ...interface{}) {
	l.logger.WithFields(l.toFields(v...)).Error(msg)
}

func (l appLogger) toFields(v ...interface{}) logrus.Fields {
	fields := make(logrus.Fields)
	keys := make([]string, 0)
	values := make([]string, 0)
	//sort.Slice(v, func(i int, j int) bool {
	//	a, _ := v[i].(string)
	//	b, _ := v[j].(string)
	//	return strings.Compare(a, b) == -1
	//})

	for i, val := range v {
		if i%2 == 0 {
			keys = append(keys, fmt.Sprintf("%v", val))
		} else {
			values = append(values, fmt.Sprintf("%v", val))
		}
	}

	if len(values) < len(keys) {
		for i := len(values); i < len(keys); i++ {
			values = append(values, "")
		}
	}

	for i, k := range keys {
		fields[k] = values[i]
	}

	return fields
}
