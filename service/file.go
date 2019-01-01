package service

import (
	"context"
	"github.com/wq1019/cloud_disk/model"
)

type fileService struct {
	model.FileStore
}

func SaveFileToFolder(ctx context.Context, file *model.File, folder *model.Folder) (err error) {
	return FromContext(ctx).SaveFileToFolder(file, folder)
}

func DeleteFile(ctx context.Context, ids []int64, folderId int64) (err error) {
	return FromContext(ctx).DeleteFile(ids, folderId)
}

func MoveFile(ctx context.Context, fromId, toId int64, fileIds []int64) (err error) {
	return FromContext(ctx).MoveFile(fromId, toId, fileIds)
}

func CopyFile(ctx context.Context, toId int64, fileIds []int64) (err error) {
	return FromContext(ctx).CopyFile(toId, fileIds)
}

func NewFileService(fs model.FileStore) model.FileService {
	return &fileService{fs}
}
