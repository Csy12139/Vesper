package sender

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

func handler(c *cli.Context) error {
	fileToSend := c.String("file")
	if fileToSend == "" {
		return fmt.Errorf("file parameter missing")
	}
	f, err := os.Open(fileToSend)
	if err != nil {
		return err
	}
	defer f.Close()
	//
	//
	return nil // 临时
}

// New creates the command
func New() *cli.Command {
	return &cli.Command{
		Name:    "send",
		Aliases: []string{"s"},
		Usage:   "Send a file",
		Action:  handler,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file, f",
				Usage: "Send content of file `FILE`",
			},
		},
	}
}
