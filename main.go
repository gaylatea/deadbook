package main

import (
	"fmt"
	"github.com/mitchellh/cli"
	"os"
)

// main() serves as the console entry point for the application and is
// also a stub, allowing us to return meaningful exit codes to the OS
// after execution is complete.
func main() {
	os.Exit(realMain())
}

// realMain() is used to return exit codes to the OS from commands,
// as main() does not have this capability.
func realMain() int {
	// Setup the CLI commands and formatting for this tool.
	// TODO(silversupreme): Make --version work like Serf's does.
	cli := cli.CLI{
		Args:     os.Args[1:],
		Commands: Commands,
	}

	exitCode, err := cli.Run()
	if err != nil {
		// TODO(silversupreme): Allow this to output red text on error.
		// Much more interesting and visually distinctive.
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
