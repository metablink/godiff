package main

import (
	"fmt"
	"os"

	"github.com/metablink/godiff/cmd"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	var app = cli.NewApp()
	app.Name = "godiff"
	app.Usage = "a powerful csv differ"

	app.Action = func(c *cli.Context) error {
		fmt.Println("Please specify a command.")
		return nil
	}

	app.Commands = []cli.Command{
		cmd.DiffCmd(),
		cmd.SummaryCmd(),
	}

	app.Run(os.Args)
}
