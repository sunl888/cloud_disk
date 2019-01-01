package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
)

type dbFolder struct {
	db *gorm.DB
}

func (f *dbFolder) ExistFolder(userId int64, folderName string) (isExist bool) {
	var (
		count uint8
	)
	f.db.Model(model.Folder{}).
		Where("user_id = ? and folder_name = ?", userId, folderName).
		Limit(1).
		Count(&count)
	if count > 0 {
		isExist = true
	}
	return
}

func (f *dbFolder) CreateFolder(folder *model.Folder) (err error) {
	err = f.db.Create(&folder).Error
	return
}

func (f *dbFolder) LoadFolder(id, userId int64, isLoadRelated bool) (folder *model.Folder, err error) {
	folder = &model.Folder{}
	q := f.db.Model(model.Folder{})
	if isLoadRelated {
		q = q.Preload("Files").
			Preload("Folders", "user_id = ?", userId)
	}
	q = q.Where("user_id = ?", userId)
	// 如果没有传目录id表示加载根目录
	if id == 0 {
		q = q.Where("level = 1")
	} else {
		q = q.Where("id = ?", id)
	}
	err = q.First(&folder).Error
	if folder.Files == nil {
		folder.Files = make([]*model.File, 0, 1)
	}
	if gorm.IsRecordNotFoundError(err) {
		err = errors.RecordNotFound("目录不存在")
	}
	return
}

func NewDBFolder(db *gorm.DB) model.FolderStore {
	return &dbFolder{db}
}
