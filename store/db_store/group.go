package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/model"
)

type dbGroup struct {
	db *gorm.DB
}

func NewDBGroup(db *gorm.DB) model.GroupStore {
	return &dbGroup{db}
}
