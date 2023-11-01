package main

import (
	"io"
	"os/exec"
)

var (
	version = "dev"
	commit  = "unavailable"
	date    = "unknown"
)

type Cmd struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (c Cmd) ShowVersion() {
	if len(commit) > 8 {
		commit = commit[:7]
	}
	_, _ = c.Stdout.Write([]byte(`envdir version ` + version + `, build ` + commit + ` (` + date + ")\n"))
}

func (c Cmd) Execute() int {
	flags := NewFlags(c.Stdout)

	if flags.Help {
		return 0
	}

	if flags.ShowVersion {
		c.ShowVersion()
		return 0
	}

	logger := NewLogger(flags, c.Stdout)

	logger.Debug("using config", LogFields{"dir": flags.Dir, "fail": flags.Fail, "log-level": flags.LogLevel, "log-format": flags.LogFormat})

	if flags.Cmd == "" {
		logger.Error("missing command", LogFields{})

		return 2
	}

	arg0, err := exec.LookPath(flags.Cmd)
	if err != nil {
		logger.Error("error running subprocess", LogFields{"err": err.Error()})

		return 1
	}
	logger.Debug("using command", LogFields{"cmd": arg0, "args": flags.Args})

	cmd := exec.Command(arg0, flags.Args...)
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr

	envBuilder := NewEnvBuilder(flags, logger)

	cmd.Env, err = envBuilder.Build()
	if err != nil {
		logger.Error("error parsing environment variables", LogFields{"err": err.Error()})

		return 3
	}

	err = cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		logger.Info("subcommand exited with error", LogFields{"err": err.Error()})

		return exitError.ExitCode()
	}

	return 0
}

func NewCmd(stdin io.Reader, stdout, stderr io.Writer) *Cmd {
	return &Cmd{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}
}
