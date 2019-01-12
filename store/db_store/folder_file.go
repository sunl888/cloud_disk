package db_store

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/model"
	"strconv"
	"strings"
)

type dbFolderFile struct {
	db *gorm.DB
}

func (f *dbFolderFile) LoadFolderFilesByFolderIdAndFileIds(folderId int64, fileIds []int64, userId int64) (folderFiles []*model.WrapFolderFile, err error) {
	folderFiles = make([]*model.WrapFolderFile, 0, 10)
	err = f.db.Table("folders fo").
		Select("ff.*,f.size as file_size,f.format as format").
		Joins("LEFT JOIN `folder_files` ff ON ff.folder_id = fo.id").
		Joins("LEFT JOIN `files` f ON ff.file_id = f.id").
		Where("fo.id = ? AND fo.user_id = ? AND ff.file_id IN (?)", folderId, userId, fileIds).
		Find(&folderFiles).Error
	return
}

func (f *dbFolderFile) LoadFolderFilesByFolderIds(folderIds []int64, userId int64) (folderFiles []*model.WrapFolderFile, err error) {
	var (
		allFolderId []int64
		likeSql     string
	)
	folderFiles = make([]*model.WrapFolderFile, 0, 10)
	for _, v := range folderIds {
		parent := model.Folder{}
		conditions := fmt.Sprintf("id = %d AND user_id = %d", v, userId)
		err := f.db.First(&parent, conditions).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				continue
			}
			return nil, err
		}
		// 将父目录的 ID 放到目录列表
		allFolderId = append(allFolderId, parent.Id)
		// 在数据库中列出所有子目录 ID
		id2Str := strconv.FormatInt(parent.Id, 10)
		likeSql += fmt.Sprintf(" `key` LIKE %s OR", "'"+parent.Key+id2Str+"-%'")
	}
	likeSql = strings.TrimRight(likeSql, "OR")
	f.db.Model(model.Folder{}).
		Where(likeSql).
		Pluck("DISTINCT id", &allFolderId)
	// 查找父目录下面所有子目录中的文件ID
	f.db.Table("folder_files ff").
		Select("ff.*,f.size as file_size,f.format as format").
		Joins("LEFT JOIN `files` f ON ff.file_id = f.id").
		Where("ff.folder_id IN (?)", allFolderId).
		Find(&folderFiles)

	return folderFiles, err
}

func (f *dbFolderFile) ExistFile(filename string, folderId, userId int64) (isExist bool, err error) {
	var count int
	err = f.db.Table("folders fo").
		Joins("LEFT JOIN `folder_files` ff ON ff.folder_id = fo.id").
		Where("fo.id = ? AND fo.user_id = ? AND ff.filename = ?", folderId, userId, filename).Limit(1).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, err
}

func NewDBFolderFile(db *gorm.DB) model.FolderFileStore {
	return &dbFolderFile{db}
}
