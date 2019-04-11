package adaptors

import (
	"log"
	"os"

	"github.com/int128/kubelogin/adaptors/interfaces"
)

func NewLogger() adaptors.Logger {
	return &Logger{
		logger: log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds),
	}
}

type stdLogger interface {
	Printf(format string, v ...interface{})
}

// Logger wraps the standard log.Logger and just provides debug level.
type Logger struct {
	logger stdLogger
	level  adaptors.DebugLevel
}

func (l *Logger) Logf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l *Logger) Debugf(level adaptors.DebugLevel, format string, v ...interface{}) {
	if level <= l.level {
		l.logger.Printf(format, v...)
	}
}

func (l *Logger) SetDebugLevel(level adaptors.DebugLevel) {
	l.level = level
}

func (l *Logger) GetDebugLevel() adaptors.DebugLevel {
	return l.level
}
