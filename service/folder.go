package service

import (
	"context"
	"github.com/wq1019/cloud_disk/model"
)

type folderService struct {
	model.FolderStore
}

func LoadFolder(ctx context.Context, id int64) (folder *model.Folder, err error) {
	return FromContext(ctx).LoadFolder(id)
}

func LoadFolderSesource(ctx context.Context, id int64) (folder *model.Folder, err error) {
	return FromContext(ctx).LoadFolderSesource(id)
}

func NewFolderService(ds model.FolderStore) model.FolderService {
	return &folderService{ds}
}
