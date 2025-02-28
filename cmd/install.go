package cmd

import (
	"github.com/Csy12139/Vesper/cmd/receiver"
	"github.com/Csy12139/Vesper/cmd/sender"
	"github.com/urfave/cli/v2"
	"sort"
)

// Install all the commands
func Install(app *cli.App) {
	app.Commands = []*cli.Command{
		sender.New(),
		receiver.New(),
	}
	sort.Sort(cli.CommandsByName(app.Commands))
}
