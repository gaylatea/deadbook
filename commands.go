package main

import (
	"github.com/mitchellh/cli"
	"github.com/silversupreme/deadbook/commands"
	"os"
)

var Commands map[string]cli.CommandFactory

func init() {
	basicUi := &cli.BasicUi{Writer: os.Stdout}

	realUi := &cli.ColoredUi{
		Ui:          basicUi,
		InfoColor:   cli.UiColorCyan,
		ErrorColor:  cli.UiColorRed,
		OutputColor: cli.UiColorGreen,
	}

	Commands = map[string]cli.CommandFactory{
		"update": func() (cli.Command, error) {
			return &commands.UpdateCommand{Ui: realUi}, nil
		},
	}
}
