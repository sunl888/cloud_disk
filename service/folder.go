package service

import "github.com/wq1019/cloud_disk/model"

type folderService struct {
	model.FolderStore
}

func NewFolderService(ds model.FolderStore) model.FolderService {
	return &folderService{ds}
}
