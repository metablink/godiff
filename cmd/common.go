package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

// GetUnimplementedAction is a placeholder action for planned, but unimplemented commands
func GetUnimplementedAction() func(c *cli.Context) error {
	return GetStaticTextAction("Unimplemented")
}

// GetStaticTextAction creates an action that displays static text
func GetStaticTextAction(message string) func(c *cli.Context) error {
	textAction := func(c *cli.Context) error {
		fmt.Println(message)
		return nil
	}

	return textAction
}
