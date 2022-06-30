package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	tmint "github.com/gschnall/tmint/tmux_interface"
	twiz "github.com/gschnall/tmint/tmux_wizard"
)

var (
	tmintSessionName   = "|_ Tmint | a Tmux session manager _|"
	currentSessionName = ":"
)

var app = cli.App{
	Action: func(c *cli.Context) error {
		return nil
	},
}

func setupCliApp() {
	app.Name = "Tmux Session Interface"
	app.Usage = "A feature packed tmux tree for managing sessions, windows, and panes."
	app.Author = "Gabe Schnall"

	app.Commands = []cli.Command{}
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name: "p", Usage: "Prevents tmint from zooming the current pane.",
			Required: false,
		},
		&cli.BoolFlag{
			Name: "r", Usage: "Activates tmint resize pane interface",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "t",
			Usage:    "Used for tmux-keybindings",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "u",
			Usage:    "Used for tmux-keybindings",
			Required: false,
		},
		&cli.StringFlag{
			Name:        "s",
			Value:       tmintSessionName,
			Usage:       "The name of the session created with -t flag.",
			Required:    false,
			Destination: &tmintSessionName,
		},
		&cli.StringFlag{
			Name:        "current-tmint-session",
			Value:       currentSessionName,
			Usage:       "The name of the current session. Used for tmux keybindings workflow.",
			Required:    false,
			Destination: &currentSessionName,
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.Bool("t") == true {
			twiz.InitTmintForTmuxKeybindings(tmintSessionName)
		} else if c.Bool("r") == true {
			tmint.StartResizeInterface()
		} else {
			tmint.Start(c.Bool("p"), currentSessionName, tmintSessionName, c.Bool("r"))
		}
		return nil
	}
}

func main() {
	setupCliApp()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
