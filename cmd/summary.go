package cmd

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

// SummaryCmd generates an aggregate summary of file differences
func SummaryCmd() cli.Command {

	summaryCmd := cli.Command{
		Name:  "summary",
		Usage: "Create a summary of file differences.",
		Action: func(c *cli.Context) error {
			fmt.Println("Not implimented")
			return nil
		},
	}
	return summaryCmd
}
