package main

import (
	"flag"
	"github.com/wq1019/cloud_disk/queue/subscribe"
	"github.com/wq1019/cloud_disk/server"
	"go.uber.org/zap"
)

var (
	h bool
	c string
)

func init() {
	flag.BoolVar(&h, "h", false, "the help")
	flag.StringVar(&c, "c", "config/config.yml", "set configuration `file`")
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return

	}
	svr := server.SetupServer(c)
	svr.Logger.Info("start queue", zap.Int("queue goroutine num", svr.Conf.QueueNum))
	subscribe.StartSubQueue(svr)
}
