package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
)

type dbFile struct {
	db *gorm.DB
}

func (f *dbFile) DeleteFile(ids []int64, folderId int64) (err error) {
	err = f.db.Exec("DELETE FROM `folder_files` WHERE folder_id = ? AND file_id IN (?)", folderId, ids).Error
	return
}

func (f *dbFile) SaveFileToFolder(file *model.File, folder *model.Folder) (err error) {
	folders := make([]*model.Folder, 0, 1)
	folders = append(folders, folder)
	err = f.db.First(&file, "hash = ?", file.Hash).Error
	if gorm.IsRecordNotFoundError(err) {
		err = errors.RecordNotFound("文件不存在")
	}
	err = f.db.Model(&file).Association("Folders").Append(folder).Error
	return
}

func NewDBFile(db *gorm.DB) model.FileStore {
	return &dbFile{db}
}
