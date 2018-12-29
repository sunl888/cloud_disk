package command

import (
	"github.com/urfave/cli"
	"github.com/wq1019/cloud_disk/server"
	"log"
)

func NewExampleCommand(svr *server.Server) cli.Command {
	return cli.Command{
		Name:  "example",
		Usage: "命令行测试",
		Action: func(c *cli.Context) error {
			log.Println("example command is ok!")
			return nil
		},
	}
}
