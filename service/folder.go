package service

import (
	"context"
	"github.com/wq1019/cloud_disk/model"
)

type folderService struct {
	model.FolderStore
}

func ListFolder(ctx context.Context, folderIds []int64, userId int64) (folder []*model.Folder, err error) {
	return FromContext(ctx).ListFolder(folderIds, userId)
}

func LoadFolder(ctx context.Context, id, userId int64, isLoadRelated bool) (folder *model.Folder, err error) {
	return FromContext(ctx).LoadFolder(id, userId, isLoadRelated)
}

func CreateFolder(ctx context.Context, folder *model.Folder) (err error) {
	return FromContext(ctx).CreateFolder(folder)
}

func ExistFolder(ctx context.Context, userId, parentId int64, folderName string) (isExist bool) {
	return FromContext(ctx).ExistFolder(userId, parentId, folderName)
}

func DeleteFolder(ctx context.Context, ids []int64, userId int64) (err error) {
	return FromContext(ctx).DeleteFolder(ids, userId)
}

func MoveFolder(ctx context.Context, to *model.Folder, ids []int64) (err error) {
	return FromContext(ctx).MoveFolder(to, ids)
}

func CopyFolder(ctx context.Context, to *model.Folder, ids []int64) (err error) {
	return FromContext(ctx).CopyFolder(to, ids)
}

func RenameFolder(ctx context.Context, id int64, newName string) (err error) {
	return FromContext(ctx).RenameFolder(id, newName)
}

func NewFolderService(ds model.FolderStore) model.FolderService {
	return &folderService{ds}
}
