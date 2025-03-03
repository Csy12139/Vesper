package sender

import (
	"fmt"
	"github.com/Csy12139/Vesper/p2p/pkg/session/sender"
	"github.com/urfave/cli/v2"
	"os"
)

func handler(c *cli.Context) error {
	fileName := c.String("file")
	if fileName == "" {
		return fmt.Errorf("file parameter missing")
	}
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	sess := sender.New(f)
	return sess.Start()
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
