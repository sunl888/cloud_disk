package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/model"
)

type dbShare struct {
	db *gorm.DB
}

func NewDBShare(db *gorm.DB) model.ShareStore {
	return &dbShare{db}
}
