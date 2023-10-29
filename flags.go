package main

import (
	"flag"
	"io"
	"os"
)

type Flags struct {
	Dir       string
	Fail      bool
	Paranoid  bool
	LogFormat string
	LogLevel  string

	Cmd  string
	Args []string
}

func (f *Flags) Getenv(envName, envDefault string) string {
	env := os.Getenv(envName)
	if env == "" {
		env = envDefault
	}

	return env
}

func NewFlags(outputBuffer io.Writer) *Flags {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(outputBuffer)
	flags := Flags{}

	flagSet.StringVar(&flags.Dir, "d", flags.Getenv("ENVDIR_DIRECTORY", "/secrets"), "Directory to read files from")
	flagSet.BoolVar(&flags.Fail, "f", flags.Getenv("ENVDIR_FAIL", "false") == "true", "Fail if missing directory")
	flagSet.BoolVar(&flags.Paranoid, "p", flags.Getenv("ENVDIR_PARANOID", "false") == "true", "Don't pass any env vars except default system ones")
	flagSet.StringVar(&flags.LogFormat, "lf", flags.Getenv("ENVDIR_LOG_FORMAT", "text"), "Log format (text/json)")
	flagSet.StringVar(&flags.LogLevel, "ll", flags.Getenv("ENVDIR_LOG_LEVEL", "warn"), "Log level (error/warn/info/debug)")

	_ = flagSet.Parse(os.Args[1:])

	args := flagSet.Args()
	if len(args) == 0 {
		flags.Cmd = ""
		flags.Args = make([]string, 0)
		return &flags
	}

	flags.Cmd = args[0]
	flags.Args = args[1:]

	return &flags
}
