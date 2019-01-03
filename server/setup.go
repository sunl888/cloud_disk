package server

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/minio/minio-go"
	"github.com/spf13/afero"
	"github.com/wq1019/cloud_disk/config"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/pkg/pubsub"
	"github.com/wq1019/cloud_disk/service"
	"github.com/zm-dev/go-file-uploader"
	fileUploaderMinio "github.com/zm-dev/go-file-uploader/minio"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

func setupGorm(debug bool, databaseConfig *config.DatabaseConfig) *gorm.DB {
	var dataSourceName string
	switch databaseConfig.Driver {
	case "sqlite3":
		dataSourceName = databaseConfig.DBName
	case "mysql":
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			databaseConfig.User,
			databaseConfig.Password,
			databaseConfig.Host+":"+databaseConfig.Port,
			databaseConfig.DBName,
		)
	}
	var (
		db  *gorm.DB
		err error
	)
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(databaseConfig.Driver, dataSourceName)
		if err == nil {
			db.LogMode(debug)
			// group by 问题
			db.Exec("set session sql_mode='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'")
			if debug {
				autoMigrate(db)
			}
			return db
		}
		log.Println(err)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("数据库链接失败！ error: %+v", err)
	return nil
}

func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.UserInfo{},
		&model.Certificate{},
		&model.File{},
		&model.Group{},
		&model.Share{},
		&model.Folder{},
		&model.FolderFile{},
	).Error
	if err != nil {
		log.Fatalf("AutoMigrate 失败！ error: %+v", err)
	}
}

func setupRedis(redisConfig *config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: redisConfig.Address + ":" + redisConfig.Port,
	})
}

func setupFilesystem(fsConfig *config.FilesystemConfig) afero.Fs {
	switch fsConfig.Driver {
	case "os":
		return afero.NewBasePathFs(afero.NewOsFs(), fsConfig.Root)
	case "memory":
		return afero.NewBasePathFs(afero.NewMemMapFs(), fsConfig.Root)
	default:
		return afero.NewBasePathFs(afero.NewOsFs(), fsConfig.Root)
	}
}

func setupFileUploader(s *Server) go_file_uploader.Uploader {
	return fileUploaderMinio.NewMinioUploader(
		go_file_uploader.HashFunc(go_file_uploader.MD5HashFunc),
		setupMinio(s),
		setupFileStore(s),
		s.Conf.Minio.BucketName,
		go_file_uploader.Hash2StorageNameFunc(go_file_uploader.TwoCharsPrefixHash2StorageNameFunc),
	)
}

func setupFileStore(s *Server) go_file_uploader.Store {
	return go_file_uploader.NewDBStore(s.DB)
}

func setupMinio(s *Server) *minio.Client {
	// TODO 这里不知道为什么不能设置为bool类型,只能用字符串表示,然后转成bool
	SslEnable := s.Conf.Minio.SSL == "true"
	minioClient, err := minio.New(
		s.Conf.Minio.Host,
		s.Conf.Minio.AccessKey,
		s.Conf.Minio.SecretKey,
		SslEnable,
	)
	if err != nil {
		log.Fatalf("minio client 创建失败! error: %+v", err)
	}
	return minioClient
}

func loadEnv(appEnv string) string {
	if appEnv == "" {
		appEnv = "production"
	}
	return appEnv
}

func setupLogger(serv *Server) *zap.Logger {
	var err error
	var logger *zap.Logger
	if serv.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatal(err)
	}
	return logger
}

func SetupServer(configPath string) *Server {
	s := &Server{}
	s.AppEnv = loadEnv(os.Getenv("APP_ENV"))
	s.Debug = os.Getenv("DEBUG") == "true"
	s.Logger = setupLogger(s)
	s.Logger.Debug("load config...")
	s.Conf = config.LoadConfig(configPath)
	s.Logger.Debug("load filesystem...")
	s.BaseFs = setupFilesystem(&s.Conf.Fs)
	s.Logger.Debug("load redis...")
	s.RedisClient = setupRedis(&s.Conf.Redis)
	s.Logger.Debug("load database...")
	s.DB = setupGorm(s.Debug, &s.Conf.DB)
	s.Logger.Debug("load service...")
	s.Pub = pubsub.NewPub(s.RedisClient, s.Logger)
	s.Service = service.NewService(s.DB, s.RedisClient, s.BaseFs, s.Conf, s.Pub)
	s.Logger.Debug("load uploader service...")
	s.FileUploader = setupFileUploader(s)
	return s
}
