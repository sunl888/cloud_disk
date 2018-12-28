package server

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/afero"
	"github.com/wq1019/cloud_disk/config"
	"github.com/wq1019/cloud_disk/pkg/pubsub"
	"github.com/wq1019/cloud_disk/service"
	"go.uber.org/zap"
)

type Server struct {
	AppEnv      string
	Debug       bool
	BaseFs      afero.Fs
	RedisClient *redis.Client
	DB          *gorm.DB
	Conf        *config.Config
	Logger      *zap.Logger
	Service     service.Service
	Pub         pubsub.PubQueue
}
