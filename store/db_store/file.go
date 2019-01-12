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

func (f *dbFile) CopyFile(fromId, toId int64, fileIds []int64) (totalSize uint64, err error) {
	// 复制指定的文件索引并插入(创建)到指定目录
	// 不能用 IN  因为只要有一个文件已存在就会导致所有文件都不会被复制,因此这里必须循环检测每个文件是否已经存在
	// EXPLAIN INSERT INTO `folder_files` SELECT 4,`file_id`,`filename` FROM `folder_files` WHERE (`folder_id` = '1' AND `file_id` IN ('1','2')) AND NOT EXISTS (SELECT * FROM `folder_files` WHERE `folder_id` = '4' AND `file_id` IN ('1','2'))
	sql := "INSERT INTO `folder_files` " +
		"SELECT ?,`file_id`,`filename` FROM `folder_files` WHERE (`folder_id` = ? AND `file_id` = ?) AND " +
		"NOT EXISTS (SELECT `folder_id` FROM `folder_files` WHERE `folder_id` = ? AND `file_id` = ?)"
	savedFileIds := make([]int64, 0, len(fileIds))
	for _, fileId := range fileIds {
		rowsAffected := f.db.Exec(sql, toId, fromId, fileId, toId, fileId).RowsAffected
		if rowsAffected > 0 {
			savedFileIds = append(savedFileIds, fileId)
		}
	}
	if len(savedFileIds) > 0 {
		files := make([]*model.File, 0, len(savedFileIds))
		f.db.Model(model.File{}).Where("id IN (?) ", savedFileIds).First(&files)
		for _, file := range files {
			totalSize += uint64(file.Size)
		}
	}
	return
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
