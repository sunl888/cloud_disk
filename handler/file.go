package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/service"
	"net/http"
)

type fileHandler struct {
}

func (*fileHandler) RenameFile(c *gin.Context) {
	l := struct {
		FileId   int64  `json:"file_id" form:"file_id"`
		FolderId int64  `json:"folder_id" form:"folder_id"`
		NewName  string `json:"new_name" form:"new_name"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(err)
		return
	}
	authId := middleware.UserId(c)
	folder, err := service.LoadFolder(c.Request.Context(), l.FolderId, authId, false)
	if err != nil {
		_ = c.Error(err)
		return
	}
	err = service.RenameFile(c.Request.Context(), folder.Id, l.FileId, l.NewName)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func NewFileHandler() *fileHandler {
	return &fileHandler{}
}
