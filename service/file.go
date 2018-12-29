package service

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/model"
)

type fileService struct {
	model.FileStore
}

func (fileService) MoveFile(fromId, toId int64) {
	panic("implement me")
}

func UpdateFile(c *gin.Context, id int64, file *model.File) (err error) {
	if service, ok := c.Value("service").(Service); ok {
		return service.UpdateFile(id, file)
	}
	return ServiceError
}

func SaveFileToFolder(c *gin.Context, file *model.File, folder *model.Folder) (err error) {
	if service, ok := c.Value("service").(Service); ok {
		return service.SaveFileToFolder(file, folder)
	}
	return ServiceError
}

func NewFileService(fs model.FileStore) model.FileService {
	return &fileService{fs}
}
