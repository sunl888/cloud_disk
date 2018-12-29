package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
)

type dbFolder struct {
	db *gorm.DB
}

func (f *dbFolder) LoadFolder(id int64) (folder *model.Folder, err error) {
	folder = &model.Folder{}
	err = f.db.First(&folder, "id = ?", id).Error
	if gorm.IsRecordNotFoundError(err) {
		err = errors.RecordNotFound("目录不存在")
	}
	return
}

func NewDBFolder(db *gorm.DB) model.FolderStore {
	return &dbFolder{db}
}
