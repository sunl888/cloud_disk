package server

import (
	"github.com/NetEase-Object-Storage/nos-golang-sdk/nosclient"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/afero"
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
	AppEnv        string
	BaseFs        afero.Fs
	DB            *gorm.DB
	Logger        *zap.Logger
	RedisClient   *redis.Client
	Conf          *config.Config
	Pub           pubsub.PubQueue
	Service       service.Service
	FileUploader  go_file_uploader.Uploader
	ImageUploader image_uploader.Uploader
	ImageUrl      image_url.URL
	NosClient     *nosclient.NosClient
	BucketName    string
}
