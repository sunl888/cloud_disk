package server

import (
	"github.com/NetEase-Object-Storage/nos-golang-sdk/nosclient"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/config"
	"github.com/wq1019/cloud_disk/pkg/pubsub"
	"github.com/wq1019/cloud_disk/service"
	"github.com/wq1019/go-file-uploader"
	"github.com/wq1019/go-image_uploader"
	"github.com/wq1019/go-image_uploader/image_url"
	"go.uber.org/zap"
)

type Server struct {
	Debug         bool
	BucketName    string
	AppEnv        string
	DB            *gorm.DB
	Logger        *zap.Logger
	ImageUrl      image_url.URL
	RedisClient   *redis.Client
	Conf          *config.Config
	Service       service.Service
	Pub           pubsub.PubQueue
	NosClient     *nosclient.NosClient
	ImageUploader image_uploader.Uploader
	FileUploader  go_file_uploader.Uploader
	//BaseFs        afero.Fs
}
