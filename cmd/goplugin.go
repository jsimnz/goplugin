package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/jsimnz/goplugin/cmd/subcommands"
)

func main() {
	app := cli.NewApp()
	app.Name = "goplugin"
	app.Usage = "Goplugin CLI program to help developers build great plugins for Go"
	app.Version = VERSION

	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate a plugin cgo interface",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path, p",
					Value: "./",
					Usage: "path to folder of plugin",
				},
			},
			Action: subcommands.GenerateCmd,
		},
		// {
		// 	Name:    "build",
		// 	Aliases: []string{"b"},
		// 	Usage:   "Build a shared library for a plugin",
		// 	Action:  subcommands.BuildCmd,
		// },
	}

	app.Run(os.Args)
}
