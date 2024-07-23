package logger

import (
	"fmt"
)

func (logger *Logger) Errorf(format string, args ...any) (int, error) {
	return logger.printf("error", format, args...)
}

func (logger *Logger) Errorln(args ...any) (int, error) {
	return logger.println("error", args...)
}

func (logger *Logger) Fatalf(format string, args ...any) (int, error) {
	return logger.printf("fatal", format, args...)
}

func (logger *Logger) Fatalln(args ...any) (int, error) {
	return logger.println("fatal", args...)
}

func (logger *Logger) Infof(format string, args ...any) (int, error) {
	return logger.printf("info", format, args...)
}

func (logger *Logger) Infoln(args ...any) (int, error) {
	return logger.println("info", args...)
}

func (logger *Logger) Panicf(format string, args ...any) (int, error) {
	return logger.printf("panic", format, args...)
}

func (logger *Logger) Panicln(args ...any) (int, error) {
	return logger.println("panic", args...)
}

func (logger *Logger) Warnf(format string, args ...any) (int, error) {
	return logger.printf("warn", format, args...)
}

func (logger *Logger) Warnln(args ...any) (int, error) {
	return logger.println("warn", args...)
}

func (logger *Logger) printf(levelName string, format string, args ...any) (int, error) {
	level := logger.GetLevel(levelName)
	if logger.GetVerbosity() < level.minVerbosity {
		return 0, nil
	}

	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	return level.clrs.Fprintf(logger.writer, "%s%s", level.prefix, fmt.Sprintf(format, args...))
}

func (logger *Logger) println(levelName string, args ...any) (int, error) {
	level := logger.GetLevel(levelName)
	if logger.GetVerbosity() < level.minVerbosity {
		return 0, nil
	}

	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	return level.clrs.Fprintln(logger.writer, fmt.Sprintf("%s%s", level.prefix, fmt.Sprintln(args...)))
}
