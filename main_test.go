package main

import (
	"os"
	"testing"
)

func TestApp(t *testing.T) {
	oldOsExit := osExit
	defer func() { osExit = oldOsExit }()

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"envdir"}

	var exitCode int
	testExit := func(code int) {
		exitCode = code
	}

	osExit = testExit
	main()

	if exitCode != 2 {
		t.Errorf("expected main function to return proper exit code, got %d", exitCode)
	}
}
