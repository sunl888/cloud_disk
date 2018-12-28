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
}

type store struct {
	model.TicketStore
	model.UserStore
	model.CertificateStore
}

func NewStore(db *gorm.DB, redisClient *redis.Client) Store {
	return &store{redis_store.NewRedisTicket(redisClient),
		db_store.NewDBUser(db),
		db_store.NewDBCertificate(db),
	}
}
