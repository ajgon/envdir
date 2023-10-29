package main

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	parentSysEnvRegex    = regexp.MustCompile(`"msg":"read value from parent process",.*"name":"HOME"`)
	parentCustomEnvRegex = regexp.MustCompile(`"msg":"read value from parent process",.*"name":"VAR_FROM_PARENT"`)
	dirEnvRegex          = regexp.MustCompile(`"msg":"read value from directory",.*"name":"VAR_FROM_DIR"`)
)

func TestCmd_Success(t *testing.T) {
	var (
		cmdStdin  bytes.Buffer
		cmdStdout bytes.Buffer
		cmdStderr bytes.Buffer
	)

	t.Setenv("VAR_FROM_PARENT", "value-from-parent")

	envDir, err := os.MkdirTemp("", "env")
	if err != nil {
		t.Fatalf("error creating temporary dir: %v", err)
	}
	defer os.RemoveAll(envDir)

	envFile := filepath.Join(envDir, "VAR_FROM_DIR")
	if err := os.WriteFile(envFile, []byte("value-from-dir"), 0644); err != nil {
		t.Fatalf("error creative temporary env var file: %v", err)
	}

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"envdir", "-d", envDir, "-f", "-lf", "json", "-ll", "debug", "sh", "-c", "env"}

	cmd := NewCmd(&cmdStdin, &cmdStdout, &cmdStderr)
	exitCode := cmd.Execute()
	output := cmdStdout.String()

	if exitCode != 0 {
		t.Errorf("expected success exit code, got %d", exitCode)
	}

	if !parentSysEnvRegex.MatchString(output) {
		t.Errorf("output is missing debug information about setting system env variable from parent process, got output:\n%s", output)
	}

	if !parentCustomEnvRegex.MatchString(output) {
		t.Errorf("output is missing debug information about setting custom env variable from parent process, got output:\n%s", output)
	}

	if !dirEnvRegex.MatchString(output) {
		t.Errorf("output is missing debug information about setting env variable from directory, got output:\n%s", output)
	}

	pathEnvValue := os.Getenv("PATH")
	if !strings.Contains(output, "\nPATH="+pathEnvValue+"\n") {
		t.Errorf("spawned process is missing env variable from parent process, got output:\n%s", output)
	}

	if !strings.Contains(output, "\nVAR_FROM_PARENT=value-from-parent\n") {
		t.Errorf("spawned process is missing env variable from parent process, got output:\n%s", output)
	}

	if !strings.Contains(output, "\nVAR_FROM_DIR=value-from-dir\n") {
		t.Errorf("spawned process is missing env variable from directory, got output:\n%s", output)
	}
}

func TestCmd_Failure(t *testing.T) {
	t.Run("it fails when subcommand is not provided", func(t *testing.T) {
		var (
			cmdStdin  bytes.Buffer
			cmdStdout bytes.Buffer
			cmdStderr bytes.Buffer
		)

		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{"envdir"}

		cmd := NewCmd(&cmdStdin, &cmdStdout, &cmdStderr)
		exitCode := cmd.Execute()
		output := cmdStdout.String()

		if exitCode != 2 {
			t.Errorf("expected command error exit code, got %d", exitCode)
		}

		if !strings.Contains(output, `level=ERROR msg="missing command"`) {
			t.Errorf("expected output to return error about missing command, output:\n%s", output)
		}
	})

	t.Run("it fails when subcommand is missing", func(t *testing.T) {
		var (
			cmdStdin  bytes.Buffer
			cmdStdout bytes.Buffer
			cmdStderr bytes.Buffer
		)

		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{"envdir", "not-existing-subcommand"}

		cmd := NewCmd(&cmdStdin, &cmdStdout, &cmdStderr)
		exitCode := cmd.Execute()
		output := cmdStdout.String()

		if exitCode != 1 {
			t.Errorf("expected command error exit code, got %d", exitCode)
		}

		if !strings.Contains(output, `level=ERROR msg="error running subprocess" err="exec: \"not-existing-subcommand\": executable file not found in $PATH"`) {
			t.Errorf("expected output to return error about missing command, output:\n%s", output)
		}
	})

	t.Run("it fails when subcommand is not executable", func(t *testing.T) {
		var (
			cmdStdin  bytes.Buffer
			cmdStdout bytes.Buffer
			cmdStderr bytes.Buffer
		)

		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{"envdir", "./main.go"}

		cmd := NewCmd(&cmdStdin, &cmdStdout, &cmdStderr)
		exitCode := cmd.Execute()
		output := cmdStdout.String()

		if exitCode != 1 {
			t.Errorf("expected command error exit code, got %d", exitCode)
		}

		if !strings.Contains(output, `level=ERROR msg="error running subprocess" err="exec: \"./main.go\": permission denied"`) {
			t.Errorf("expected output to return error about missing command, output:\n%s", output)
		}
	})

	t.Run("it passes exit code of subcommand if it fails", func(t *testing.T) {
		var (
			cmdStdin  bytes.Buffer
			cmdStdout bytes.Buffer
			cmdStderr bytes.Buffer
		)

		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{"envdir", "-ll", "info", "false"}

		cmd := NewCmd(&cmdStdin, &cmdStdout, &cmdStderr)
		exitCode := cmd.Execute()
		output := cmdStdout.String()

		if exitCode != 1 {
			t.Errorf("expected command error exit code, got %d", exitCode)
		}

		if !strings.Contains(output, `level=INFO msg="subcommand exited with error" err="exit status 1"`) {
			t.Errorf("expected output to return error about missing command, output:\n%s", output)
		}
	})

	t.Run("it returns error if there was problem with env builder", func(t *testing.T) {
		var (
			cmdStdin  bytes.Buffer
			cmdStdout bytes.Buffer
			cmdStderr bytes.Buffer
		)

		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{"envdir", "-f", "-d", "/non-existing-directory", "true"}

		cmd := NewCmd(&cmdStdin, &cmdStdout, &cmdStderr)
		exitCode := cmd.Execute()
		output := cmdStdout.String()

		if exitCode != 3 {
			t.Errorf("expected command error exit code, got %d", exitCode)
		}

		if !strings.Contains(
			output,
			`level=ERROR msg="error parsing environment variables" err="error reading variables from directory: `+
				`open /non-existing-directory: no such file or directory"`,
		) {
			t.Errorf("expected output to return error about missing command, output:\n%s", output)
		}
	})
}
