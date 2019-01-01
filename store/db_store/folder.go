package db_store

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
	"strconv"
	"strings"
)

type dbFolder struct {
	db *gorm.DB
}

func updateKey(parentKey, key, startId string) string {
	keys := strings.Split(key, "-")
	for index, key := range keys {
		if key == startId {
			return parentKey + strings.Join(keys[index:], "-")
		}
	}
	return ""
}

func replaceKey(idMap map[int64]int64, parentKey, key, startId string) string {
	newKey := updateKey(parentKey, key, startId)
	if newKey == "" {
		return ""
	}
	keys := strings.Split(key, "-")
	for index, key := range keys {
		key2Int64, _ := strconv.ParseInt(key, 10, 64)
		if newId, ok := idMap[key2Int64]; ok {
			id2Str := strconv.FormatInt(newId, 10)
			keys[index] = id2Str
		}
	}
	return strings.Join(keys, "-")
}

func (f *dbFolder) CopyFolder(to *model.Folder, ids []int64) (err error) {
	var (
		rootFolder model.Folder    // 将要移动的第一层目录
		toId2Str   string          // 移动到的目录 ID 字符串形式
		tmpFolder  model.Folder    // 临时 folder
		children   []*model.Folder // 子目录
		id2Str     string          // 移动的目录的 ID 字符串形式
		idMap      map[int64]int64
	)
	children = make([]*model.Folder, 0, 5)
	idMap = make(map[int64]int64, len(ids))
	toId2Str = strconv.FormatInt(to.Id, 10)
	for _, id := range ids {
		id2Str = strconv.FormatInt(id, 10)
		err := f.db.First(&rootFolder, "id = ?", id).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				continue
			}
			return err
		}
		tmpFolder = rootFolder
		// 查询所有子目录
		f.db.Model(model.Folder{}).Where("`key` LIKE ?", rootFolder.Key+id2Str+"-%").Order("id ASC").Find(children)
		newFolder := model.Folder{
			UserId:     rootFolder.UserId,
			FolderName: rootFolder.FolderName,
			Level:      to.Level + 1,
			ParentId:   to.Id,
			Key:        to.Key + toId2Str + model.FolderKeyPrefix,
		}
		// 创建一个与原根目录相等的根目录
		f.db.Create(&newFolder)
		idMap[rootFolder.Id] = newFolder.Id
		// 创建子目录
		for i := 0; i < len(children); i++ {
			pId := strconv.FormatInt(children[i].ParentId, 10)
			key := replaceKey(idMap, newFolder.Key, children[i].Key, pId)
			newChildFolder := model.Folder{
				UserId:     rootFolder.UserId,
				FolderName: rootFolder.FolderName,
				Key:        key,                         // todo
				ParentId:   idMap[children[i].ParentId], // todo
				Level:      newFolder.Level + (children[i].Level - tmpFolder.Level),
			}
			f.db.Create(&newChildFolder)
			idMap[children[i].Id] = newChildFolder.Id
		}
		// 创建新的文件关联
		type FolderFile struct {
			FolderId int64
			FileId   int64
		}
		var (
			oldFolderIds []int64
			folderFiles  []*FolderFile
		)
		for k := range idMap {
			oldFolderIds = append(oldFolderIds, k)
		}
		f.db.Table("folder_files").Where("folder_id IN (?)", oldFolderIds).Scan(&folderFiles)
		// INSERT INTO `users` VALUES (?,?,?),(?,?,?)
		sql := "INSERT INTO `folder_files` (`folder_id`,`file_id`)VALUES "
		for _, v := range folderFiles {
			sql += fmt.Sprintf("(%d,%d),", idMap[v.FolderId], v.FileId)
		}
		sql = strings.TrimRight(sql, ",")
		f.db.Exec(sql)
		children = nil
		rootFolder = model.Folder{}
	}
	return nil
}

func (f *dbFolder) MoveFolder(to *model.Folder, ids []int64) (err error) {
	var (
		rootFolder model.Folder    // 将要移动的第一层目录
		toId2Str   string          // 移动到的目录 ID 字符串形式
		tmpFolder  model.Folder    // 临时 folder
		children   []*model.Folder // 子目录
		id2Str     string          // 移动的目录的 ID 字符串形式
	)
	toId2Str = strconv.FormatInt(to.Id, 10)
	for _, id := range ids {
		id2Str = strconv.FormatInt(id, 10)
		err := f.db.First(&rootFolder, "id = ?", id).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				continue
			}
			return err
		}
		// 查询所有子目录
		f.db.Model(model.Folder{}).Where("`key` LIKE ?", rootFolder.Key+id2Str+"-%").Find(children)

		tmpFolder = rootFolder
		// 更新根目录的信息
		f.db.Model(rootFolder).Updates(model.Folder{
			Level:    to.Level + 1,
			ParentId: to.Id,
			Key:      to.Key + toId2Str + model.FolderKeyPrefix,
		})
		for _, child := range children {
			f.db.Model(child).Updates(model.Folder{
				Level: rootFolder.Level + (child.Level - tmpFolder.Level),
				Key:   updateKey(rootFolder.Key, child.Key, id2Str),
			})
		}
		children = nil
		rootFolder = model.Folder{}
	}
	return nil
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
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				continue
			}
			return err
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
