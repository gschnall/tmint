package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"

	tmint "github.com/gschnall/tmint/tmux_interface"
	// twiz "github.com/gschnall/tmint/tmux_wizard"
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
			Name:     "p",
			Usage:    "Prevents tmint from zooming the current pane.",
			Required: false,
		},
	}

	app.Action = func(c *cli.Context) error {
		tmint.Start(c.Bool("p"))
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

	// On data refresh this will always be true - have to avoid setting it again
	// if !isUserPaneZoomed {
	// 	twiz.TmuxToggleFullscreen()
	// }
}
