package cmd

import (
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

		paths := [2]string{from, to}
		files := [2]*os.File{}

		for idx, path := range paths {

			// If the path wasn't set by it's flag, it can still be an argument
			if path == "" {
				// No path provided. Error out.
				return cli.NewExitError("diff requires two file paths to run (see the --to and --from flags)", 1)
			}

			var err error
			files[idx], err = os.Open(path)

			if err != nil {
				return cli.NewExitError(err, 1)
			}
		}

		if err := lib.DiffFile(files[0], files[1]); err != nil {
			return cli.NewExitError(err, 1)
		}

		return nil
	}

	return diffCmd
}
