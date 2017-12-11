package main

import (
	"os"

	"github.com/metablink/godiff/cmd"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	var app = cli.NewApp()
	app.Name = "godiff"
	app.Usage = "a powerful csv differ"
	app.Action = cmd.GetStaticTextAction("Please specify a command.")
	app.Commands = []cli.Command{
		cmd.DiffCmd(),
		cmd.SummaryCmd(),
	}

	app.Run(os.Args)
}
