package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/service"
	"net/http"
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
	folder := model.Folder{
		UserId:     authId,
		Level:      parentFolder.Level + 1,
		ParentId:   l.ParentId,
		Key:        parentFolder.Key + pId2String + model.FolderKeyPrefix,
		FolderName: l.FolderName,
	}
	err = service.CreateFolder(c.Request.Context(), &folder)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, folder)
}

func (*folderHandler) DeleteSource(c *gin.Context) {
	l := struct {
		FileIds         []int64 `json:"file_ids" form:"file_ids"`
		FolderIds       []int64 `json:"folder_ids" form:"folder_ids"`
		CurrentFolderId int64   `json:"current_folder_id" form:"current_folder_id"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(err)
		return
	}
	if len(l.FileIds) == 0 && len(l.FolderIds) == 0 {
		_ = c.Error(errors.BadRequest("请指定要删除的文件或者目录ID"))
		return
	}
	authId := middleware.UserId(c)
	// 删除指定的文件
	fmt.Println(l.FileIds)
	if len(l.FileIds) > 0 {
		currentFolder, err := service.LoadFolder(c.Request.Context(), l.CurrentFolderId, authId, false)
		if err != nil {
			_ = c.Error(err)
			return
		}
		err = service.DeleteFile(c.Request.Context(), l.FileIds, currentFolder.Id)
		if err != nil {
			_ = c.Error(err)
			return
		}
	}
	// 删除目录列表
	if len(l.FolderIds) > 0 {
		err := service.DeleteFolder(c.Request.Context(), l.FolderIds, authId)
		if err != nil {
			_ = c.Error(err)
			return
		}
	}
	c.Status(http.StatusNoContent)
}

func NewFolder() *folderHandler {
	return &folderHandler{}
}
