package main

import (
	"bytes"
	"os"
	"testing"
)

var flagsOutput bytes.Buffer

func Test_FlagsDefaults(t *testing.T) {
	t.Setenv("ENVDIR_DIRECTORY", "")
	t.Setenv("ENVDIR_FAIL", "")
	t.Setenv("ENVDIR_PARANOID", "")
	t.Setenv("ENVDIR_LOG_FORMAT", "")
	t.Setenv("ENVDIR_LOG_LEVEL", "")
	flags := NewFlags(&flagsOutput)

	var tests = []struct {
		flagName     string
		flagValue    any
		defaultValue any
	}{
		{"d", flags.Dir, "/secrets"},
		{"f", flags.Fail, false},
		{"p", flags.Paranoid, false},
		{"lf", flags.LogFormat, "text"},
		{"ll", flags.LogLevel, "warn"},
	}

	for _, tt := range tests {
		switch defaultValue := tt.defaultValue.(type) {
		case string:
			flagValue := tt.flagValue.(string)
			if flagValue != defaultValue {
				t.Errorf("invalid default value of flag %q: expected %q, got %q", tt.flagName, flagValue, defaultValue)
			}
		case bool:
			flagValue := tt.flagValue.(bool)
			if flagValue != defaultValue {
				t.Errorf("invalid default value of flag %q: expected %t, got %t", tt.flagName, flagValue, defaultValue)
			}
		default:
			t.Fatal("broken flags default test")
		}
	}
}

func Test_FlagsFromEnv(t *testing.T) {
	t.Setenv("ENVDIR_DIRECTORY", "/test")
	t.Setenv("ENVDIR_FAIL", "true")
	t.Setenv("ENVDIR_PARANOID", "true")
	t.Setenv("ENVDIR_LOG_FORMAT", "json")
	t.Setenv("ENVDIR_LOG_LEVEL", "debug")
	flags := NewFlags(&flagsOutput)

	var tests = []struct {
		flagName  string
		flagEnv   string
		flagValue any
		envValue  any
	}{
		{"d", "ENVDIR_DIRECTORY", flags.Dir, "/test"},
		{"f", "ENVDIR_FAIL", flags.Fail, true},
		{"p", "ENVDIR_PARANOID", flags.Paranoid, true},
		{"lf", "ENVDIR_LOG_FORMAT", flags.LogFormat, "json"},
		{"ll", "ENVDIR_LOG_LEVEL", flags.LogLevel, "debug"},
	}

	for _, tt := range tests {
		switch envValue := tt.envValue.(type) {
		case string:
			flagValue := tt.flagValue.(string)
			if flagValue != envValue {
				t.Errorf("invalid env value of flag %q: expected %q, got %q", tt.flagName, flagValue, envValue)
			}
		case bool:
			flagValue := tt.flagValue.(bool)

			if flagValue != envValue {
				t.Errorf("invalid env value of flag %q: expected %t, got %t", tt.flagName, flagValue, envValue)
			}
		default:
			t.Fatal("broken flags default test")
		}
	}
}

func Test_FlagsFromArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"envdir", "-d", "/dir", "-f", "-p", "-lf", "json", "-ll", "error", "sh", "-c", "ls -l"}
	flags := NewFlags(&flagsOutput)

	var tests = []struct {
		flagName     string
		flagValue    any
		defaultValue any
	}{
		{"args0", flags.Cmd, "sh"},
		{"cmd1", flags.Args[0], "-c"},
		{"cmd2", flags.Args[1], "ls -l"},
		{"d", flags.Dir, "/dir"},
		{"f", flags.Fail, true},
		{"p", flags.Paranoid, true},
		{"lf", flags.LogFormat, "json"},
		{"ll", flags.LogLevel, "error"},
	}

	for _, tt := range tests {
		switch defaultValue := tt.defaultValue.(type) {
		case string:
			flagValue := tt.flagValue.(string)
			if flagValue != defaultValue {
				t.Errorf("invalid default value of flag %q: expected %q, got %q", tt.flagName, flagValue, defaultValue)
			}
		case bool:
			flagValue := tt.flagValue.(bool)
			if flagValue != defaultValue {
				t.Errorf("invalid default value of flag %q: expected %t, got %t", tt.flagName, flagValue, defaultValue)
			}
		default:
			t.Fatal("broken flags default test")
		}
	}
}

func Test_FlagsEmptyArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"envdir"}
	flags := NewFlags(&flagsOutput)

	var tests = []struct {
		flagName     string
		flagValue    any
		defaultValue any
	}{
		{"args0", flags.Cmd, ""},
		{"d", flags.Dir, "/secrets"},
		{"f", flags.Fail, false},
		{"p", flags.Paranoid, false},
		{"lf", flags.LogFormat, "text"},
		{"ll", flags.LogLevel, "warn"},
	}

	for _, tt := range tests {
		switch defaultValue := tt.defaultValue.(type) {
		case string:
			flagValue := tt.flagValue.(string)
			if flagValue != defaultValue {
				t.Errorf("invalid default value of flag %q: expected %q, got %q", tt.flagName, flagValue, defaultValue)
			}
		case bool:
			flagValue := tt.flagValue.(bool)
			if flagValue != defaultValue {
				t.Errorf("invalid default value of flag %q: expected %t, got %t", tt.flagName, flagValue, defaultValue)
			}
		default:
			t.Fatal("broken flags default test")
		}
	}

}
