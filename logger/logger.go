package logger

import (
	"io"
	"os"
	"sync"
	"sync/atomic"
)

const (
	None = iota
	Err
	Warn
	Info
	All
)

type Logger struct {
	writer    io.Writer
	mutex     sync.Mutex
	levels    map[string]Level
	verbosity atomic.Int32
}

func New(writer io.Writer, levels map[string]Level, verbosity int) *Logger {
	logger := new(Logger)
	logger.SetWriter(writer).SetLevels(levels).SetVerbosity(verbosity)

	return logger
}

func Default() *Logger {
	return New(os.Stderr, nil, Err)
}

func (logger *Logger) SetWriter(writer io.Writer) *Logger {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	logger.writer = writer
	return logger
}

func (logger *Logger) GetVerbosity() int {
	return int(logger.verbosity.Load())
}

func (logger *Logger) SetVerbosity(verbosity int) *Logger {
	logger.verbosity.Store(int32(verbosity))
	return logger
}

func (logger *Logger) SetLevels(levels map[string]Level) *Logger {

	if logger.levels == nil {
		logger.levels = make(map[string]Level)
	}

	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	for name, level := range levels {
		logger.levels[name] = level
	}
	return logger
}

func (logger *Logger) GetLevel(name string) Level {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	level, ok := logger.levels[name]
	if !ok {
		level = DefaultLevel(name)
		logger.levels[name] = level
	}
	return level
}
