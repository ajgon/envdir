package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type EnvBuilder struct {
	Flags  *Flags
	Logger *Logger
}

func (eb *EnvBuilder) parentEnvs() []string {
	if !eb.Flags.Paranoid {
		parentEnvVars := os.Environ()

		if eb.Logger.LogLevel == slog.LevelDebug {
			for _, envLine := range parentEnvVars {
				envName, envValue, _ := strings.Cut(envLine, `=`)
				eb.Logger.Debug("read value from parent process", LogFields{"name": envName, "value": envValue})
			}
		}

		return parentEnvVars
	}

	parentEnvs := make([]string, 0)

	for _, envName := range []string{"HOME", "HOSTNAME", "PATH", "PWD", "TERM", "TZ", "UMASK"} {
		envValue := os.Getenv(envName)
		parentEnvs = append(parentEnvs, envName+`=`+envValue)

		eb.Logger.Debug("read value from parent process", LogFields{"name": envName, "value": envValue})
	}

	return parentEnvs
}

func (eb *EnvBuilder) flagError(err error) ([]string, error) {
	if eb.Flags.Fail {
		return nil, err
	}

	return make([]string, 0), nil
}

func (eb *EnvBuilder) directoryEnvs() ([]string, error) {
	envFiles, err := os.ReadDir(eb.Flags.Dir)
	if err != nil {
		return eb.flagError(err)
	}

	dirEnvs := make([]string, 0)

	for _, envFile := range envFiles {
		envPath := filepath.Join(eb.Flags.Dir, envFile.Name())
		envData, err := os.ReadFile(envPath)
		if err != nil {
			return nil, fmt.Errorf("reading env file `%s`: %w", envPath, err)
		}

		envValue := strings.TrimSuffix(string(envData), "\n")

		eb.Logger.Debug("read value from directory", LogFields{"name": envFile.Name(), "value": envValue})

		dirEnvs = append(dirEnvs, envFile.Name()+`=`+envValue)
	}

	return dirEnvs, nil
}

func (eb *EnvBuilder) Build() ([]string, error) {
	dirEnvs, err := eb.directoryEnvs()
	if err != nil {
		return nil, fmt.Errorf("error reading variables from directory: %w", err)
	}

	return append(eb.parentEnvs(), dirEnvs...), nil
}

func NewEnvBuilder(flags *Flags, logger *Logger) *EnvBuilder {
	return &EnvBuilder{
		Flags:  flags,
		Logger: logger,
	}
}
