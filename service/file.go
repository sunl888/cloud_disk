package service

import (
	"context"
	"github.com/wq1019/cloud_disk/model"
)

type fileService struct {
	model.FileStore
}

func (fileService) MoveFile(fromId, toId int64) {
	panic("implement me")
}

func UpdateFile(ctx context.Context, id int64, file *model.File) (err error) {
	return FromContext(ctx).UpdateFile(id, file)
}

func SaveFileToFolder(ctx context.Context, file *model.File, folder *model.Folder) (err error) {
	return FromContext(ctx).SaveFileToFolder(file, folder)
}

func NewFileService(fs model.FileStore) model.FileService {
	return &fileService{fs}
}
