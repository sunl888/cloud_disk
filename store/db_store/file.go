package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
)

type dbFile struct {
	db *gorm.DB
}

func (f *dbFile) LoadFile(folderId, fileId, userId int64) (file *model.File, err error) {
	file = &model.File{}
	err = f.db.Table("folders fo").
		Select("ff.file_id as id, ff.filename, f.hash, f.format, f.extra, f.size, f.created_at, f.updated_at").
		Joins("LEFT JOIN `folder_files` ff ON ff.folder_id = fo.id").
		Joins("LEFT JOIN `files` f ON f.id = ff.origin_file_id").
		Where("fo.id = ? AND fo.user_id = ? AND ff.file_id = ?", folderId, userId, fileId).
		Limit(1).
		Scan(&file).
		Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = errors.RecordNotFound("文件不存在")
		}
		return
	}
	return
}

func (f *dbFile) RenameFile(folderId, fileId int64, newName string) (err error) {
	var count int
	f.db.Model(model.FolderFile{}).Where("folder_id = ? AND filename = ?", folderId, newName).Limit(1).Count(&count)
	if count > 0 {
		return errors.FileAlreadyExist("该目录下已经存在同名文件")
	} else {
		err = f.db.Model(model.FolderFile{}).
			Where("file_id = ?", fileId).
			Update("filename", newName).
			Error
		if gorm.IsRecordNotFoundError(err) {
			err = errors.RecordNotFound("文件不存在")
		}
	}
	return err
}

// fromId != toId
func (f *dbFile) CopyFile(fromId, toId int64, fileIds []int64) (totalSize uint64, err error) {
	savedFileIds := make([]int64, 0, len(fileIds))
	for _, fileId := range fileIds {
		var (
			fromFile model.FolderFile
			count    int
		)
		// 查询源文件信息
		err = f.db.Model(model.FolderFile{}).
			Where("`folder_id` = ? AND `file_id` = ?", fromId, fileId).
			First(&fromFile).
			Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				continue
			}
			return
		}
		// 查询目标目录有没有同名文件
		err = f.db.Model(model.FolderFile{}).
			Where("`folder_id` = ? AND `filename` = ?", toId, fromFile.Filename).
			Limit(1).
			Count(&count).
			Error
		if err != nil {
			return
		}
		// 移动到的目录已经存在同名文件
		if count > 0 {
			continue
		}
		err = f.db.Create(&model.FolderFile{
			OriginFileId: fromFile.OriginFileId,
			FolderId:     toId,
			Filename:     fromFile.Filename,
		}).Error
		if err != nil {
			return
		}
		savedFileIds = append(savedFileIds, fromFile.OriginFileId)
	}
	// 计算复制的文件大小
	if len(savedFileIds) > 0 {
		fileSizes := make([]int64, 0, len(savedFileIds))
		f.db.Table("files").Where("id IN (?)", savedFileIds).Pluck("size", &fileSizes)
		for _, size := range fileSizes {
			totalSize += uint64(size)
		}
	}
	return
}

func (f *dbFile) MoveFile(fromId, toId int64, fileIds []int64) (err error) {
	for _, fileId := range fileIds {
		var (
			fromFile model.FolderFile
			count    int
		)
		err = f.db.Model(model.FolderFile{}).
			Where("`folder_id` = ? AND `file_id` = ?", fromId, fileId).
			First(&fromFile).
			Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				continue
			}
			return
		}
		err = f.db.Model(model.FolderFile{}).
			Where("`folder_id` = ? AND `filename` = ?", toId, fromFile.Filename).
			Limit(1).
			Count(&count).
			Error
		if err != nil {
			return
		}
		// 移动到的目录已经存在同名文件
		if count > 0 {
			continue
		} else {
			err = f.db.Model(&fromFile).Update("folder_id", toId).Error
			if err != nil {
				return
			}
		}
	}
	return
}

func (f *dbFile) DeleteFile(ids []int64, folderId int64) (err error) {
	err = f.db.Exec("DELETE FROM `folder_files` WHERE `folder_id` = ? AND `file_id` IN (?)", folderId, ids).Error
	return
}

func (f *dbFile) SaveFileToFolder(file *model.File, folderId int64) (err error) {
	err = f.db.Model(model.File{}).Create(
		&model.FolderFile{
			FolderId:     folderId,
			OriginFileId: file.Id,
			Filename:     file.Filename,
		}).Error
	return
}

func NewDBFile(db *gorm.DB) model.FileStore {
	return &dbFile{db}
}
