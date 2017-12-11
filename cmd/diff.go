package cmd

import (
	"errors"
	"os"

	"github.com/metablink/godiff/lib"
	"gopkg.in/urfave/cli.v1"
)

var (
	fromPath string
	toPath   string
)

// DiffCmd generates a detailed diff of the given files
func DiffCmd() cli.Command {

	diffCmd := cli.Command{
		Name:   "diff",
		Usage:  "Create a complete set of file differences.",
		Flags:  diffFlags(),
		Action: diffAction(),
	}

	return diffCmd
}

func diffFlags() []cli.Flag {

	diffFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "from, f",
			Usage:       "1st file for comparison",
			Destination: &fromPath,
		},
		cli.StringFlag{
			Name:        "to, t",
			Usage:       "2nd file for comparison",
			Destination: &toPath,
		},
	}

	return diffFlags
}

func diffAction() func(c *cli.Context) error {
	diffAction := func(c *cli.Context) error {
		paths := []string{fromPath, toPath}
		files, err := checkPaths(paths)

		// Bail out if there was a file path problem
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if err := lib.DiffFile(files[0], files[1]); err != nil {
			return cli.NewExitError(err, 1)
		}

		return nil
	}

	return diffAction
}

func checkPaths(paths []string) (files []*os.File, err error) {

	files = make([]*os.File, len(paths))

	for idx, path := range paths {

		// If the path wasn't set by it's flag, it can still be an argument
		if path == "" {
			// No path provided. Error out.
			return nil, errors.New("diff requires two file paths to run (see the --to and --from flags)")
		}

		var err error
		files[idx], err = os.Open(path)

		if err != nil {
			return nil, err
		}
	}

	return files, nil
}
