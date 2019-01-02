package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
)

type dbFile struct {
	db *gorm.DB
}

func (f *dbFile) RenameFile(folderId, fileId int64, newName string) (err error) {
	err = f.db.Model(model.FolderFile{}).
		Where("folder_id = ? AND file_id = ?", folderId, fileId).
		Update("filename", newName).
		Error
	return err
}

func (f *dbFile) CopyFile(toId int64, fileIds []int64) (err error) {
	//TODO 复制文件时需要提供folderId
	//for _, fileId := range fileIds {
	//	f.db.Table("folder_files").First(&model.FolderFile{},"folder_id = ? AND file_id = ?",fromId,fileId)
	//	f.db.Table("folder_files").
	//		FirstOrCreate(&model.FolderFile{
	//			FolderId: toId,
	//			FileId:   fileId,
	//			Filename:fileId,
	//		}, "folder_id = ? AND file_id = ?", toId, fileId)
	//}
	return nil
}

func (f *dbFile) MoveFile(fromId, toId int64, fileIds []int64) (err error) {
	err = f.db.Table("folder_files").
		Where("folder_id = ? AND file_id IN (?)", fromId, fileIds).
		Update(map[string]interface{}{
			"folder_id": toId,
		}).Error
	return
}

func (f *dbFile) DeleteFile(ids []int64, folderId int64) (err error) {
	err = f.db.Exec("DELETE FROM `folder_files` WHERE `folder_id` = ? AND `file_id` IN (?)", folderId, ids).Error
	return
}

func (f *dbFile) SaveFileToFolder(file *model.File, folder *model.Folder) (err error) {
	err = f.db.First(&file, "`hash` = ?", file.Hash).Error
	if gorm.IsRecordNotFoundError(err) {
		err = errors.RecordNotFound("文件不存在")
	}
	err = f.db.Table("folder_files").
		Create(&model.FolderFile{
			FolderId: folder.Id,
			FileId:   file.Id,
			Filename: file.Filename,
		}).Error
	return
}

func NewDBFile(db *gorm.DB) model.FileStore {
	return &dbFile{db}
}
