package main

import (
	"github.com/Csy12139/Vesper/p2p/cmd"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func run(args []string) error {
	app := cli.NewApp()
	app.Name = "vesper"
	cmd.Install(app)

	return app.Run(args)
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
