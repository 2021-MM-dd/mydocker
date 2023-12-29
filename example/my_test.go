package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"testing"
)

func TestCommandLine(t *testing.T) {
	app := cli.NewApp()
	// 指定全局参数
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config, c", Usage: "Load configuration from `FILE`"},
	}
	// 指定支持的命令
	app.Commands = []cli.Command{{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "add a task to the list",
		Action: func(c *cli.Context) error {
			log.Panicln("run command add")
			return nil
		},
		Flags: []cli.Flag{cli.Int64Flag{
			Name:  "priority",
			Value: 1,
			Usage: "priority of this task",
		}},
	}}
	//运行
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
