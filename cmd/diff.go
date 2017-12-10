package cmd

import (
	"log"
	"os"

	"github.com/metablink/godiff/lib"
	"gopkg.in/urfave/cli.v1"
)

// DiffCmd generates a detailed diff of the given files
func DiffCmd() cli.Command {
	var (
		from string
		to   string
	)

	diffCmd := cli.Command{
		Name:  "diff",
		Usage: "Create a complete set of file differences.",
	}

	diffCmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "from, f",
			Usage:       "1st file for comparison",
			Destination: &from,
		},
		cli.StringFlag{
			Name:        "to, t",
			Usage:       "2nd file for comparison",
			Destination: &to,
		},
	}

	diffCmd.Action = func(c *cli.Context) error {

		fromFile, err := os.Open(from)
		if err != nil {
			log.Fatal(err)
		}

		toFile, err := os.Open(to)
		if err != nil {
			log.Fatal(err)
		}

		lib.DiffFile(fromFile, toFile)

		return nil
	}

	return diffCmd
}
