package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/model"
)

type dbFolder struct {
	db *gorm.DB
}

func NewDBFolder(db *gorm.DB) model.FolderStore {
	return &dbFolder{db}
}
