package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/model"
)

type dbFile struct {
	db *gorm.DB
}

func (f *dbFile) DeleteFile(id int64) error {
	panic("implement me")
}

func (f *dbFile) DeletedFileList(limit, offset int64) (files []*model.File, count int64, err error) {
	panic("implement me")
}

func (f *dbFile) RecoverFile(id int64) error {
	panic("implement me")
}

func (f *dbFile) ListFile(limit, offset int64) (files []*model.File, count int64, err error) {
	panic("implement me")
}

func (f *dbFile) IsExistFile(id, userId int64) (isExist bool, err error) {
	panic("implement me")
}

func (f *dbFile) UpdateFile(id int64, file *model.File) (err error) {
	panic("implement me")
}

func NewDBFile(db *gorm.DB) model.FileStore {
	return &dbFile{db}
}
