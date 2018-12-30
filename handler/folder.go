package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/service"
	"strconv"
)

type folderHandler struct {
}

func (*folderHandler) LoadFolderSesource(c *gin.Context) {
	folderId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.BadRequest("id 格式不正确", nil))
		return
	}
	if folderId == 0 {
		_ = c.Error(errors.BadRequest("请指定上传的文件夹", nil))
		return
	}
	authId := middleware.UserId(c)
	folder, err := service.LoadFolder(c.Request.Context(), folderId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if authId != folder.UserId {
		_ = c.Error(errors.Unauthorized("没有访问权限"))
		return
	}
	folder, err = service.LoadFolderSesource(c.Request.Context(), folderId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, folder)
}

func NewFolder() *folderHandler {
	return &folderHandler{}
}
