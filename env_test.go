package main

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	pathError *fs.PathError
	envOutput bytes.Buffer
)

func Test_BuildFailFlag(t *testing.T) {
	logger := NewLogger(&Flags{}, &envOutput)

	t.Run("it fails when directory does not exist and fail flag is set", func(t *testing.T) {
		t.Parallel()

		flags := &Flags{Fail: true}
		envBuilder := NewEnvBuilder(flags, logger)
		result, err := envBuilder.Build()
		if result != nil {
			t.Errorf("expected nil result, got %v", result)
		}

		if !errors.As(err, &pathError) {
			t.Errorf("expected *fs.PathError, got %v", err)
		}
	})

	t.Run("it returns parent process results, when directory does not exist and fail flag is not set", func(t *testing.T) {
		t.Parallel()

		flags := &Flags{Fail: false}
		envBuilder := NewEnvBuilder(flags, logger)
		result, err := envBuilder.Build()
		if len(result) < 1 {
			t.Errorf("expected result from parent process, got %v", result)
		}

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}

func Test_BuildParanoidFlag(t *testing.T) {
	logger := NewLogger(&Flags{}, &envOutput)

	t.Run("it passes only basic environment variables if flag is set", func(t *testing.T) {
		t.Parallel()

		flags := &Flags{Fail: false, Paranoid: true}
		os.Setenv("TEST_ENV", "lorem-ipsum")

		envBuilder := NewEnvBuilder(flags, logger)
		result, err := envBuilder.Build()

	LOOP:
		for _, env := range result {
			envName, _, _ := strings.Cut(env, `=`)
			for _, expectedEnvName := range []string{"HOME", "HOSTNAME", "PATH", "PWD", "TERM", "TZ", "UMASK"} {
				if envName == expectedEnvName {
					continue LOOP
				}
			}

			t.Errorf("expected only default environment variables, got %q", envName)
		}

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("it passes all environment variables if flag is not set", func(t *testing.T) {
		t.Parallel()

		flags := &Flags{Fail: false, Paranoid: false}
		os.Setenv("TEST_ENV", "lorem-ipsum")

		envBuilder := NewEnvBuilder(flags, logger)
		result, err := envBuilder.Build()

		foundEnv := false
		for _, env := range result {
			envName, _, _ := strings.Cut(env, `=`)
			if envName == "TEST_ENV" {
				foundEnv = true
			}
		}

		if !foundEnv {
			t.Error("expected custom env variable to be passed, but it was not found")
		}

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}

func Test_Build(t *testing.T) {
	t.Parallel()

	logger := NewLogger(&Flags{}, &envOutput)

	envDir, err := os.MkdirTemp("", "env")
	if err != nil {
		t.Fatalf("error creating temporary dir: %v", err)
	}
	defer os.RemoveAll(envDir)

	envFile := filepath.Join(envDir, "VAR_FROM_DIR")
	if err := os.WriteFile(envFile, []byte("value-from-dir"), 0644); err != nil {
		t.Fatalf("error creating temporary env var file: %v", err)
	}

	dummyDir := filepath.Join(envDir, "directory")
	if err := os.MkdirAll(dummyDir, 0755); err != nil {
		t.Fatalf("error creating temporary subdir: %v", err)
	}

	if err = os.Symlink(envFile, filepath.Join(envDir, "VAR_FROM_DIR_SYMLINK")); err != nil {
		t.Fatalf("error creative temporary env var file symlink: %v", err)
	}

	if err = os.Symlink(dummyDir, filepath.Join(envDir, "symlink-dir")); err != nil {
		t.Fatalf("error creating temporary subdir symlink: %v", err)
	}

	t.Run("it properly parses variables from existing directory", func(t *testing.T) {
		flags := &Flags{Fail: false, Dir: envDir}

		envBuilder := NewEnvBuilder(flags, logger)
		result, err := envBuilder.Build()

		foundEnv := false
		foundSymlinkEnv := false
		for _, env := range result {
			envName, envValue, _ := strings.Cut(env, `=`)
			if envName == "VAR_FROM_DIR" && envValue == "value-from-dir" {
				foundEnv = true
			}
			if envName == "VAR_FROM_DIR_SYMLINK" && envValue == "value-from-dir" {
				foundSymlinkEnv = true
			}
		}

		if !foundEnv {
			t.Error("expected variable from file in directory to be found, but it was missing")
		}

		if !foundSymlinkEnv {
			t.Error("expected variable from symlink to file in directory to be found, but it was missing")
		}

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("it fails if variable from dir cannot be read", func(t *testing.T) {
		err := os.Chmod(envFile, 0000)
		if err != nil {
			t.Fatal("error setting permissions to test file")
		}

		defer func() {
			err := os.Chmod(envFile, 0666)
			if err != nil {
				t.Fatal("error setting permissions to test file")
			}
		}()
		flags := &Flags{Fail: false, Dir: envDir}

		envBuilder := NewEnvBuilder(flags, logger)
		result, err := envBuilder.Build()

		if result != nil {
			t.Errorf("expected nil result, got %v", result)
		}

		if !errors.As(err, &pathError) {
			t.Errorf("expected *fs.PathError, got %v", err)
		}
	})

	t.Run("it ensures that variable from directory takes precedence over the one from parent", func(t *testing.T) {
		flags := &Flags{Fail: false, Dir: envDir}
		os.Setenv("VAR_FROM_DIR", "value-from-parent")

		envBuilder := NewEnvBuilder(flags, logger)
		result, err := envBuilder.Build()

		foundEnv := false
		for _, env := range result {
			envName, envValue, _ := strings.Cut(env, `=`)
			if envName == "VAR_FROM_DIR" && envValue == "value-from-dir" {
				foundEnv = true
			}
		}

		if !foundEnv {
			t.Error("expected variable from directory with proper value to be found, but it was missing")
		}

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}
