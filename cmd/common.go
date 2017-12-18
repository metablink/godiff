package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

// GetPrintTextAction creates an action that displays static text
func GetPrintTextAction(message string) func(c *cli.Context) error {
	textAction := func(c *cli.Context) error {
		fmt.Println(message)
		return nil
	}

	return textAction
}
