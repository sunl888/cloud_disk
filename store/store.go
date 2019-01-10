package store

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/store/db_store"
	"github.com/wq1019/cloud_disk/store/redis_store"
)

type Store interface {
	model.TicketStore
	model.UserStore
	model.CertificateStore
	model.FileStore
	model.ShareStore
	model.FolderStore
	model.GroupStore
	model.UserInfoStore
	model.FolderFileStore
}

type store struct {
	model.TicketStore
	model.UserStore
	model.CertificateStore
	model.FileStore
	model.ShareStore
	model.FolderStore
	model.GroupStore
	model.UserInfoStore
	model.FolderFileStore
}

func NewStore(db *gorm.DB, redisClient *redis.Client) Store {
	return &store{
		redis_store.NewRedisTicket(redisClient),
		db_store.NewDBUser(db),
		db_store.NewDBCertificate(db),
		db_store.NewDBFile(db),
		db_store.NewDBShare(db),
		db_store.NewDBFolder(db),
		db_store.NewDBGroup(db),
		db_store.NewDBUserInfo(db),
		db_store.NewDBFolderFile(db),
	}
}
