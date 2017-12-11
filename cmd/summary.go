package cmd

import (
	"gopkg.in/urfave/cli.v1"
)

// SummaryCmd generates an aggregate summary of file differences
func SummaryCmd() cli.Command {

	summaryCmd := cli.Command{
		Name:   "summary",
		Usage:  "Create a summary of file differences.",
		Action: GetUnimplementedAction(),
	}
	return summaryCmd
}
