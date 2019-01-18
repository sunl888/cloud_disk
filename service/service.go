package service

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/config"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/pkg/hasher"
	"github.com/wq1019/cloud_disk/pkg/pubsub"
	"github.com/wq1019/cloud_disk/store"
	"runtime"
	"time"
)

type Service interface {
	model.TicketService
	model.UserService
	model.CertificateService
	model.FileService
	model.GroupService
	model.ShareService
	model.FolderService
	model.FolderFileService
}

type service struct {
	model.TicketService
	model.UserService
	model.CertificateService
	model.FileService
	model.GroupService
	model.ShareService
	model.FolderService
	model.FolderFileService
}

func NewService(db *gorm.DB, redisClient *redis.Client /*baseFs afero.Fs,*/, conf *config.Config, pub pubsub.PubQueue) Service {
	s := store.NewStore(db, redisClient)
	tSvc := NewTicketService(s, time.Duration(conf.Ticket.TTL)*time.Second)
	h := hasher.NewArgon2Hasher(
		[]byte(conf.AppSalt),
		3,
		32<<10,
		uint8(runtime.NumCPU()),
		32,
	)
	return &service{
		tSvc,
		NewUserService(s, s, tSvc, h),
		NewCertificateService(s),
		NewFileService(s),
		NewGroupService(s),
		NewShareService(s),
		NewFolderService(s),
		NewFolderFileService(s),
	}
}
