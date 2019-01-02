package main

import (
	"flag"
	"github.com/rs/cors"
	_ "github.com/wq1019/cloud_disk/docs"
	"github.com/wq1019/cloud_disk/handler"
	"github.com/wq1019/cloud_disk/server"
	"go.uber.org/zap"
	"log"
	"net/http"
)

var (
	h bool
	c string
)

func init() {
	flag.BoolVar(&h, "h", false, "the help")
	flag.StringVar(&c, "c", "config/config.yml", "set the relative path of the configuration `file`.")
}

// @title 云盘 Api 服务
// @version 1.0
// @description 云盘的 Api 服务.
// @termsOfService https://github.com/zm-dev
// @contact.name API Support
// @contact.url https://github.com/wq1019
// @contact.email 2013855675@qq.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api
func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return

	}
	svr := server.SetupServer(c)
	svr.Logger.Info("listen", zap.String("addr", svr.Conf.ServerAddr))
	// cors 跨域用
	log.Fatal(http.ListenAndServe(svr.Conf.ServerAddr, cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "DELETE", "PUT", "HEAD"},
		AllowCredentials: true,
	}).Handler(handler.CreateHTTPHandler(svr))))
}
