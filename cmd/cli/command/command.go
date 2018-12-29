package command

import (
	"github.com/urfave/cli"
	"github.com/wq1019/cloud_disk/server"
)

func RegisterCommand(svr *server.Server) []cli.Command {
	return []cli.Command{
		NewExampleCommand(svr),
	}
}
