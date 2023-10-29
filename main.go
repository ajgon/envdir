package main

import "os"

var osExit = os.Exit

func main() {
	c := NewCmd(os.Stdin, os.Stdout, os.Stderr)
	osExit(c.Execute())
}
