package main

import (
	"github.com/urfave/cli"
	"github.com/wq1019/cloud_disk/cmd/cli/command"
	"github.com/wq1019/cloud_disk/server"
	"log"
	"os"
)

var (
	defaultPath = "config/config.yml"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config,c",
			Value: defaultPath,
			Usage: "set configuration `file`",
		},
	}
	svr := server.SetupServer(getConfigPathFromArgs())
	app.Name = "命令行工具"
	app.Usage = "haha"
	app.Version = "1.0.1"
	app.Commands = append(app.Commands, command.RegisterCommand(svr)...)
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func getConfigPathFromArgs() string {
	var (
		exist      = false
		configPath = defaultPath
	)
	for _, v := range os.Args {
		if exist {
			configPath = v
			break
		} else if v == "-c" || v == "--config" {
			exist = true
		}
	}
	return configPath
}
