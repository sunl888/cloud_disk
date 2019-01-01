package db_store

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
	"strconv"
)

type dbFolder struct {
	db *gorm.DB
}

// 查找d的所有子孙节点: select * from table_name where `key` like "${d.key}${d.id}-%"
// 查找子节点: select * from table_name where `key` like "${d.key}${d.id}-%" and level=${d.level}+1

func (f *dbFolder) DeleteFolder(ids []int64, userId int64) (err error) {
	var (
		parentFolder     model.Folder
		conditions       string
		waitDelFolderIds []int64
	)
	for _, v := range ids {
		conditions = fmt.Sprintf("id = %d AND user_id = %d", v, userId)
		err := f.db.First(&parentFolder, conditions).Error
		if gorm.IsRecordNotFoundError(err) {
			continue
		}
		// 将父目录的 ID 放到待删除的目录列表, 准备删除该目录下面的文件
		waitDelFolderIds = append(waitDelFolderIds, parentFolder.Id)
		// 删除父目录
		f.db.Delete(&parentFolder)
		// 在数据库中列出所有子目录 ID
		id2String := strconv.FormatInt(parentFolder.Id, 10)
		f.db.Model(model.Folder{}).
			Where("`key` LIKE ?", parentFolder.Key+id2String+"-%").
			Pluck("id", &waitDelFolderIds)
		// 删除父目录下面的所有子目录
		f.db.Delete(&model.Folder{}, "id IN (?)", waitDelFolderIds)
		// 删除父目录下面所有子目录中的文件
		f.db.Exec("DELETE FROM `folder_files` WHERE folder_id IN (?)", waitDelFolderIds)
	}
	return nil
}

func (f *dbFolder) ExistFolder(userId int64, folderName string) (isExist bool) {
	var (
		count uint8
	)
	f.db.Model(model.Folder{}).
		Where("user_id = ? AND folder_name = ?", userId, folderName).
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
