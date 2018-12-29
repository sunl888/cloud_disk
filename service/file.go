package service

import "github.com/wq1019/cloud_disk/model"

type fileService struct {
	model.FileStore
}

func (fileService) MoveFile(fromId, toId int64) {
	panic("implement me")
}

func NewFileService(fs model.FileStore) model.FileService {
	return &fileService{fs}
}
