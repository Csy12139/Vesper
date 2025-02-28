package receiver

import (
	"github.com/urfave/cli/v2"
)

func handler(c *cli.Context) error {
	//
	//
	return nil // 临时
}

// New creates the command
func New() *cli.Command {
	return &cli.Command{
		Name:    "receive",
		Aliases: []string{"r"},
		Usage:   "Receive a file",
		Action:  handler,
	}
}
