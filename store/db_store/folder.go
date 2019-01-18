package db_store

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
	"strconv"
	"strings"
)

type dbFolder struct {
	db *gorm.DB
}

func (f *dbFolder) ListFolder(folderIds []int64, userId int64) (folders []*model.Folder, err error) {
	// 去重
	ids := hashset.New()
	for _, v := range folderIds {
		ids.Add(v)
	}
	if ids.Size() <= 0 {
		return nil, errors.RecordNotFound("没有目录")
	}
	folders = make([]*model.Folder, 10)
	err = f.db.Model(&model.Folder{}).
		Where("user_id = ? AND id IN (?)", userId, ids.Values()).
		Find(&folders).
		Error
	return
}

func (f *dbFolder) RenameFolder(id, currentFolderId int64, newName string) (err error) {
	var count int
	err = f.db.Table("`folders` fo").
		Where("fo.parent_id = ? AND fo.folder_name = ?", currentFolderId, newName).
		Limit(1).
		Count(&count).
		Error
	if err != nil {
		return
	}
	if count > 0 {
		return model.FolderAlreadyExisted
	} else {
		err = f.db.Model(model.Folder{}).
			Where("id = ?", id).
			Update("folder_name", newName).
			Error
	}
	return
}

func (f *dbFolder) CopyFolder(to *model.Folder, waitCopyFoders []*model.Folder) (totalSize uint64, err error) {
	var (
		toId2Str string      // 移动到的目录 ID 字符串形式
		userId   = to.UserId // 用户 id
	)
	toId2Str = strconv.FormatInt(to.Id, 10)
	for _, waitFolder := range waitCopyFoders {
		var (
			waitCopyId2Str = strconv.FormatInt(waitFolder.Id, 10) // 等待移动的目录ID string 形式
			children       = make([]*model.Folder, 0, 5)
			idMap          = make(map[int64]int64, 3)
			pIdMap         = make(map[int64]int64, 3)
		)
		// 查询所有子目录
		f.db.Model(model.Folder{}).Where("`key` LIKE ?", waitFolder.Key+waitCopyId2Str+"-%").Order("id ASC").Find(&children)
		newRootFolder := model.Folder{
			UserId:     userId,
			FolderName: waitFolder.FolderName,
			Level:      to.Level + 1,
			ParentId:   to.Id, // new parentID
			Key:        to.Key + toId2Str + model.FolderKeyPrefix,
		}
		var count int
		f.db.Model(model.Folder{}).Where("user_id = ? AND folder_name = ? AND parent_id = ?",
			userId, newRootFolder.FolderName, newRootFolder.ParentId).Limit(1).Count(&count)
		if count > 0 {
			continue // 目录已存在, 跳过直接复制下一个目录
		}
		// 创建一个与原根目录相等的根目录
		f.db.Create(&newRootFolder)
		idMap[waitFolder.Id] = newRootFolder.Id
		pIdMap[newRootFolder.Id] = 0

		// 创建子目录
		newFolders := make(map[int64]*model.Folder, len(children))
		for i := 0; i < len(children); i++ {
			newChildFolder := model.Folder{
				UserId:     userId,
				FolderName: children[i].FolderName,
				Key:        children[i].Key,                                              // default
				ParentId:   children[i].ParentId,                                         // default
				Level:      newRootFolder.Level + (children[i].Level - waitFolder.Level), // must >0
			}
			f.db.Create(&newChildFolder)
			idMap[children[i].Id] = newChildFolder.Id
			pIdMap[newChildFolder.Id] = newChildFolder.ParentId

			newFolders[newChildFolder.Id] = &newChildFolder
		}
		// 更新所有新的 child 目录 的 key 和 parentId
		for id, folder := range newFolders {
			key := newRootFolder.Key
			tmpKey := ""
			pId := folder.ParentId
			for i := int64(0); i < folder.Level-newRootFolder.Level; i++ {
				tmpKey = fmt.Sprintf("%d-", idMap[pId]) + tmpKey
				pId = pIdMap[idMap[pId]]
			}
			newParentId := idMap[folder.ParentId]
			f.db.Model(model.Folder{}).Where("id = ?", id).Updates(model.Folder{
				Key:      key + tmpKey,
				ParentId: newParentId,
			})
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
		err = f.db.Table("folder_files").Where("folder_id IN (?)", oldFolderIds).Scan(&folderFiles).Error
		if err != nil {
			return
		}
		// 文件索引创建,因为目录都是新创建的,所以不可能会出现文件已存在的情况
		sql := "INSERT INTO `folder_files` SELECT ?,`origin_file_id`,`filename`,NULL FROM `folder_files` WHERE `folder_id` = ? AND `file_id` = ?"
		for _, v := range folderFiles {
			newFolderId := idMap[v.FolderId]
			rowsAffected := f.db.Exec(sql, newFolderId, v.FolderId, v.FileId).RowsAffected
			if rowsAffected > 0 {
				sizes := make([]int64, 0, 1)
				// 成功复制一个文件索引就为用户的使用空间加上这个文件占用的空间
				f.db.Table("folder_files ff").
					Joins("LEFT JOIN `files` f ON ff.origin_file_id = f.id").
					Where("ff.file_id = ?", v.FileId).Pluck("f.size", &sizes)
				totalSize += uint64(sizes[0])
			}
		}
	}
	return
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
		f.db.Model(model.Folder{}).Where("`key` LIKE ?", rootFolder.Key+id2Str+"-%").Find(&children)

		tmpFolder = rootFolder
		// 更新根目录的信息
		f.db.Model(&rootFolder).Updates(model.Folder{
			Level:    to.Level + 1,
			ParentId: to.Id,
			Key:      to.Key + toId2Str + model.FolderKeyPrefix,
		})
		for _, child := range children {
			f.db.Model(&child).Updates(model.Folder{
				Level: rootFolder.Level + (child.Level - tmpFolder.Level),
				Key:   updateKey(rootFolder.Key, child.Key, id2Str),
			})
		}
		children = nil
		rootFolder = model.Folder{}
	}
	return nil
}

func (f *dbFolder) DeleteFolder(ids []int64, userId int64) (allowDelFileHashList []string, err error) {
	var (
		waitDelFolderIds []int64
		likeSql          string
	)
	allowDelFileHashList = make([]string, 0, len(ids)*2)
	for _, v := range ids {
		relativeRootFolder := model.Folder{}
		conditions := fmt.Sprintf("id = %d AND user_id = %d", v, userId)
		err := f.db.First(&relativeRootFolder, conditions).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				continue
			}
			return nil, err
		}
		// 将父目录的 ID 放到待删除的目录列表, 准备删除该目录下面的文件
		waitDelFolderIds = append(waitDelFolderIds, relativeRootFolder.Id)
		// 在数据库中列出所有子目录 ID
		id2Str := strconv.FormatInt(relativeRootFolder.Id, 10)
		likeSql += fmt.Sprintf(" `key` LIKE %s OR", "'"+relativeRootFolder.Key+id2Str+"-%'")
	}
	if likeSql == "" {
		return nil, errors.RecordNotFound("没有要删除的记录")
	}
	likeSql = strings.TrimRight(likeSql, "OR")
	f.db.Model(model.Folder{}).
		Where(likeSql).
		Pluck("DISTINCT id", &waitDelFolderIds)

	// 删除父目录以及下面的所有子目录
	f.db.Delete(&model.Folder{}, "id IN (?)", waitDelFolderIds)

	// 统计每个文件的引用次数, 如果该文件只被引用了一次, 则可以去 minio 中将这个文件直接删除
	var originFileIds []int64
	f.db.Table("folder_files").
		Where("folder_id IN (?)", waitDelFolderIds).
		Pluck("DISTINCT origin_file_id", &originFileIds)
	for _, id := range originFileIds {
		var count int8
		f.db.Table("folder_files").
			Where("`folder_id` NOT IN (?)", waitDelFolderIds).
			Where("`origin_file_id` = ?", id).
			Limit(1).Count(&count)
		// 如果源文件被引用超过一次则表示别的目录或者别的用户也使用了这个文件, 就不用删除该文件, 只要删除该目录和该文件之间的关联即可
		if count <= 0 {
			f.db.Table("files").Where("`id` = ?", id).Pluck("`hash`", &allowDelFileHashList)
		}
	}

	// 删除父目录下面所有子目录中的文件
	f.db.Exec("DELETE FROM `folder_files` WHERE folder_id IN (?)", waitDelFolderIds)

	if len(allowDelFileHashList) > 0 {
		// 在数据库中删除所有被引用了一次的文件
		f.db.Exec("DELETE FROM `files` WHERE `hash` IN (?)", allowDelFileHashList)
	}
	return
}

func (f *dbFolder) ExistFolder(userId, parentId int64, folderName string) (isExist bool) {
	var (
		count uint8
	)
	f.db.Model(model.Folder{}).
		Where("user_id = ? AND folder_name = ? AND parent_id = ?", userId, folderName, parentId).
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
	var (
		files []*model.File
	)
	folder = &model.Folder{}
	files = make([]*model.File, 0, 1)
	q := f.db.Model(model.Folder{})
	if isLoadRelated {
		q = q.Preload("Folders", "user_id = ?", userId) // 此语句是在 #232 行时才执行的
	}
	q = q.Where("user_id = ?", userId)
	// 如果没有传目录id表示加载根目录
	if id == 0 {
		q = q.Where("level = 1")
	} else {
		q = q.Where("id = ?", id)
	}
	err = q.First(&folder).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = errors.RecordNotFound("目录不存在")
		}
		return nil, err
	}
	if isLoadRelated {
		f.db.Table("folders fo").
			Select("ff.file_id as id, ff.filename, f.hash, f.format, f.extra, f.size, f.created_at, f.updated_at").
			Joins("INNER JOIN `folder_files` ff ON ff.folder_id = fo.id").
			Joins("INNER JOIN `files` f ON f.id = ff.origin_file_id").
			Where("fo.id = ?", folder.Id).Find(&files)
	}
	folder.Files = files
	return
}

func (f *dbFolder) LoadSimpleFolder(id, userId int64) (folder *model.SimpleFolder, err error) {
	folder = &model.SimpleFolder{}
	files := make([]*model.SimpleFile, 0, 1)

	if id == 0 {
		return nil, errors.NotFound("目录 ID 不能为空")
	}
	if userId == 0 {
		return nil, errors.NotFound("用户 ID 不能为空")
	}
	q := f.db.Model(model.Folder{})
	f.db.Table("folders fo").
		Select("ff.file_id as id, ff.filename").
		Joins("LEFT JOIN `folder_files` ff ON ff.folder_id = fo.id").
		Where("fo.id = ?", id).Scan(&files)

	q = q.Where("id = ? AND user_id = ?", id, userId)
	err = q.Scan(&folder).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = errors.RecordNotFound("目录不存在")
		}
		return nil, err
	}
	folder.Files = files
	return
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

func NewDBFolder(db *gorm.DB) model.FolderStore {
	return &dbFolder{db}
}
