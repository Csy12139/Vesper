package main

import (
	"github.com/Csy12139/Vesper/cmd"
	"github.com/Csy12139/Vesper/log"
	"github.com/urfave/cli/v2"
	"os"
)

func run(args []string) error {
	app := cli.NewApp()
	app.Name = "vesper"

	cmd.Install(app)

	return app.Run(args)
}

func main() {
	if err := log.InitLog("./", 1024, 10, "info"); err != nil {
		panic(err)
	}
	log.Debug("this is debug log")
	log.Info("this is info log")
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
