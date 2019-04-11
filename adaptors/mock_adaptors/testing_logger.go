package mock_adaptors

import (
	"github.com/int128/kubelogin/adaptors/interfaces"
)

func NewLogger(l testingLogger) adaptors.Logger {
	return &logger{l}
}

type testingLogger interface {
	Logf(format string, v ...interface{})
}

type logger struct {
	testingLogger
}

func (l *logger) Debugf(level adaptors.DebugLevel, format string, v ...interface{}) {
	l.Logf(format, v...)
}

func (*logger) SetDebugLevel(level adaptors.DebugLevel) {
}

func (*logger) GetDebugLevel() adaptors.DebugLevel {
	return 1
}
