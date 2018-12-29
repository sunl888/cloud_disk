package service

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/model"
)

type folderService struct {
	model.FolderStore
}

func LoadFolder(c *gin.Context, id int64) (folder *model.Folder, err error) {
	if service, ok := c.Value("service").(Service); ok {
		return service.LoadFolder(id)
	}
	return nil, ServiceError
}

func NewFolderService(ds model.FolderStore) model.FolderService {
	return &folderService{ds}
}
