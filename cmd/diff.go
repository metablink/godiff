package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"github.com/metablink/godiff/lib"
	"gopkg.in/urfave/cli.v1"
)

var (
	fromPath     string
	toPath       string
	keyColumn    string
	ignoreFields string
)

// DiffCmd generates a detailed diff of the given files
func DiffCmd() cli.Command {
	return cli.Command{
		Name:   "diff",
		Usage:  "Create a complete set of file differences.",
		Flags:  diffFlags(),
		Action: diffAction(),
	}
}

func diffFlags() []cli.Flag {
	return []cli.Flag{
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
		cli.StringFlag{
			Name:        "key, k",
			Usage:       "Key column for the comparison",
			Destination: &keyColumn,
		},
		cli.StringFlag{
			Name:        "ignore, i",
			Usage:       "Comma-separated list of columns to ignore",
			Destination: &ignoreFields,
		},
	}
}

func diffAction() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		paths := []string{fromPath, toPath}
		files, err := checkPaths(paths)

		// Bail out if there was a file path problem
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		fromProvider := &lib.CsvRowProvider{Reader: csv.NewReader(files[0])}
		toProvider := &lib.CsvRowProvider{Reader: csv.NewReader(files[1])}
		ignoreMap := lib.StringToSet(ignoreFields, ",")

		diffStats := lib.NewDiffStats(fromProvider, toProvider, keyColumn, ignoreMap)

		if err := diffStats.Diff(); err != nil {
			return cli.NewExitError(err, 1)
		}

		diffStats.Print(fmt.Printf)

		return nil
	}
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
