package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/model"
)

type dbUserInfo struct {
	db *gorm.DB
}

func NewDBUserInfo(db *gorm.DB) model.UserInfoStore {
	return &dbUserInfo{db}
}
