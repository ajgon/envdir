package main

import (
	"io"
	"log/slog"
)

var logLevelMapper = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

type LogFields map[string]any

type Logger struct {
	Slog     *slog.Logger
	LogLevel slog.Level
}

func (l Logger) Level() slog.Level {
	return l.LogLevel
}

func (l *Logger) buildArgs(fields LogFields) []any {
	args := make([]any, len(fields)*2)
	i := 0

	for key, value := range fields {
		args[i] = key
		args[i+1] = value
		i += 2
	}

	return args
}

func (l *Logger) Debug(msg string, fields LogFields) {
	l.Slog.Debug(msg, l.buildArgs(fields)...)
}

func (l *Logger) Info(msg string, fields LogFields) {
	l.Slog.Info(msg, l.buildArgs(fields)...)
}

func (l *Logger) Warn(msg string, fields LogFields) {
	l.Slog.Warn(msg, l.buildArgs(fields)...)
}

func (l *Logger) Error(msg string, fields LogFields) {
	l.Slog.Error(msg, l.buildArgs(fields)...)
}

func NewLogger(flags *Flags, output io.Writer) *Logger {
	logger := Logger{}
	logger.LogLevel = logLevelMapper[flags.LogLevel]

	if flags.LogFormat == "json" {
		logger.Slog = slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{Level: logger}))
	} else {
		logger.Slog = slog.New(slog.NewTextHandler(output, &slog.HandlerOptions{Level: logger}))
	}

	return &logger
}
