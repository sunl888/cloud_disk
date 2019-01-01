package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/service"
	"strconv"
)

type folderHandler struct {
}

func (*folderHandler) LoadFolder(c *gin.Context) {
	l := struct {
		FolderId int64 `json:"folder_id" form:"folder_id"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(errors.BadRequest("id 格式不正确", err))
		return
	}
	authId := middleware.UserId(c)
	folder, err := service.LoadFolder(c.Request.Context(), l.FolderId, authId, true)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if authId != folder.UserId {
		_ = c.Error(errors.Unauthorized("没有访问权限"))
		return
	}
	c.JSON(200, folder)
}

func (*folderHandler) CreateFolder(c *gin.Context) {
	l := struct {
		ParentId   int64  `json:"parent_id" form:"parent_id"`
		FolderName string `json:"folder_name" form:"folder_name"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(err)
		return
	}
	if l.FolderName == "" {
		_ = c.Error(errors.BadRequest("目录名称不能为空"))
		return
	}
	authId := middleware.UserId(c)
	parentFolder, err := service.LoadFolder(c.Request.Context(), l.ParentId, authId, false)
	if err != nil {
		_ = c.Error(err)
		return
	}
	// 通过 userID 和 ID 组合查询,因此这里不用判断了
	//if authId != parentFolder.UserId {
	//	_ = c.Error(errors.Unauthorized("没有访问权限"))
	//	return
	//}
	isExist := service.ExistFolder(c.Request.Context(), authId, l.FolderName)
	if isExist {
		_ = c.Error(errors.BadRequest("目录已经存在"))
		return
	}
	pId2String := strconv.FormatInt(parentFolder.Id, 10)
	err = service.CreateFolder(c.Request.Context(), &model.Folder{
		UserId:     authId,
		Level:      parentFolder.Level + 1,
		ParentId:   l.ParentId,
		Key:        parentFolder.Key + pId2String + model.FolderKeyPrefix,
		FolderName: l.FolderName,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
}

func NewFolder() *folderHandler {
	return &folderHandler{}
}
