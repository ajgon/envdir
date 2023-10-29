package main

import (
	"bytes"
	"regexp"
	"testing"
)

var (
	debugJSONRegex = regexp.MustCompile(
		`(?m)[^\n]*"level":"DEBUG","msg":"debug","level":"debug".*\n.*"level":"INFO","msg":"info","level":"info".*\n` +
			`.*"level":"WARN","msg":"warn","level":"warn".*\n.*"level":"ERROR","msg":"error","level":"error"`,
	)
	debugRegex = regexp.MustCompile(
		`(?m)[^\n]*level=DEBUG msg=debug level=debug.*\n.*level=INFO msg=info level=info.*\n.*level=WARN msg=warn level=warn.*\n.*level=ERROR msg=error level=error`,
	)
	infoRegex = regexp.MustCompile(
		`(?m)[^\n]*level=INFO msg=info level=info.*\n.*level=WARN msg=warn level=warn.*\n.*level=ERROR msg=error level=error`,
	)
	warnRegex = regexp.MustCompile(
		`(?m)[^\n]*level=WARN msg=warn level=warn.*\n.*level=ERROR msg=error level=error`,
	)
	errorRegex = regexp.MustCompile(
		`(?m)[^\n]*level=ERROR msg=error level=error`,
	)
)

func TestLogger_Debug(t *testing.T) {
	var loggerOutput bytes.Buffer
	logger := NewLogger(&Flags{LogLevel: "debug", LogFormat: "text"}, &loggerOutput)

	logger.Debug("debug", LogFields{"level": "debug"})
	logger.Info("info", LogFields{"level": "info"})
	logger.Warn("warn", LogFields{"level": "warn"})
	logger.Error("error", LogFields{"level": "error"})

	output := loggerOutput.String()

	if !debugRegex.MatchString(output) {
		t.Errorf("debug level logger is missing some log lines:\n%s", output)
	}
}

func TestLogger_Info(t *testing.T) {
	var loggerOutput bytes.Buffer
	logger := NewLogger(&Flags{LogLevel: "info", LogFormat: "text"}, &loggerOutput)

	logger.Debug("debug", LogFields{"level": "debug"})
	logger.Info("info", LogFields{"level": "info"})
	logger.Warn("warn", LogFields{"level": "warn"})
	logger.Error("error", LogFields{"level": "error"})

	output := loggerOutput.String()

	if !infoRegex.MatchString(output) {
		t.Errorf("info level logger is missing some log lines:\n%s", output)
	}
}

func TestLogger_Warn(t *testing.T) {
	var loggerOutput bytes.Buffer
	logger := NewLogger(&Flags{LogLevel: "warn", LogFormat: "text"}, &loggerOutput)

	logger.Debug("debug", LogFields{"level": "debug"})
	logger.Info("info", LogFields{"level": "info"})
	logger.Warn("warn", LogFields{"level": "warn"})
	logger.Error("error", LogFields{"level": "error"})

	output := loggerOutput.String()

	if !warnRegex.MatchString(output) {
		t.Errorf("warn level logger is missing some log lines:\n%s", output)
	}
}

func TestLogger_Error(t *testing.T) {
	var loggerOutput bytes.Buffer
	logger := NewLogger(&Flags{LogLevel: "error", LogFormat: "text"}, &loggerOutput)

	logger.Debug("debug", LogFields{"level": "debug"})
	logger.Info("info", LogFields{"level": "info"})
	logger.Warn("warn", LogFields{"level": "warn"})
	logger.Error("error", LogFields{"level": "error"})

	output := loggerOutput.String()

	if !errorRegex.MatchString(output) {
		t.Errorf("error level logger is missing some log lines:\n%s", output)
	}
}

func TestLogger_JSON(t *testing.T) {
	var loggerOutput bytes.Buffer
	logger := NewLogger(&Flags{LogLevel: "debug", LogFormat: "json"}, &loggerOutput)

	logger.Debug("debug", LogFields{"level": "debug"})
	logger.Info("info", LogFields{"level": "info"})
	logger.Warn("warn", LogFields{"level": "warn"})
	logger.Error("error", LogFields{"level": "error"})

	output := loggerOutput.String()

	if !debugJSONRegex.MatchString(output) {
		t.Errorf("invalid format of json logger:\n%s", output)
	}
}
