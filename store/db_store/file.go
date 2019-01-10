package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
)

type dbFile struct {
	db *gorm.DB
}

//SELECT f.id... FROM folders fo LEFT JOIN `folder_files` ff ON ff.folder_id = fo.id LEFT JOIN `files` f ON f.id = ff.file_id
// WHERE (fo.id = '1' AND fo.user_id = '1' AND ff.file_id = '2') LIMIT 1
func (f *dbFile) LoadFile(folderId, fileId, userId int64) (file *model.File, err error) {
	file = &model.File{}
	err = f.db.Table("folders fo").
		Select("f.id, ff.filename, f.hash, f.format, f.extra, f.size, f.created_at, f.updated_at").
		Joins("LEFT JOIN `folder_files` ff ON ff.folder_id = fo.id").
		Joins("LEFT JOIN `files` f ON f.id = ff.file_id").
		Where("fo.id = ? AND fo.user_id = ? AND ff.file_id = ?", folderId, userId, fileId).Limit(1).
		Scan(&file).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = errors.RecordNotFound("文件不存在")
		}
		return
	}
	return
}

func (f *dbFile) RenameFile(folderId, fileId int64, newName string) (err error) {
	err = f.db.Model(model.FolderFile{}).
		Where("folder_id = ? AND file_id = ?", folderId, fileId).
		Update("filename", newName).
		Error
	if gorm.IsRecordNotFoundError(err) {
		err = errors.RecordNotFound("文件不存在")
	}
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
	var (
		count int8
	)
	f.db.Model(model.FolderFile{}).
		Where("folder_id = ? AND file_id = ?", folder.Id, file.Id).
		Limit(1).
		Count(&count)

	// 文件已经存在
	if count > 0 {
		return errors.FileAlreadyExist(nil)
	} else {
		err = f.db.Model(model.File{}).Create(
			&model.FolderFile{
				FolderId: folder.Id,
				FileId:   file.Id,
				Filename: file.Filename,
			}).Error
	}
	return
}

func NewDBFile(db *gorm.DB) model.FileStore {
	return &dbFile{db}
}
