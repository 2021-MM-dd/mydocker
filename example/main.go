package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
)

var Usage = "mydocker is a simple container implementation"

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = Usage
	app.Commands = []cli.Command{
		runCommand,
		initCommand,
	}
	app.Before = func(context *cli.Context) error { return nil }
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
