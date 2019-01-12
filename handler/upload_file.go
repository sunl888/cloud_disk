package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/service"
	"github.com/wq1019/go-file-uploader"
	"net/http"
)

type uploadFile struct {
	u go_file_uploader.Uploader
}

func (uf *uploadFile) UploadFile(c *gin.Context) {
	l := struct {
		FolderId int64 `json:"folder_id" form:"folder_id"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}
	if l.FolderId == 0 {
		_ = c.Error(errors.BadRequest("请指定上传的文件夹", nil))
		return
	}
	authId := middleware.UserId(c)
	folder, err := service.LoadFolder(c.Request.Context(), l.FolderId, authId, false)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if authId != folder.UserId {
		_ = c.Error(errors.Unauthorized("没有访问权限"))
		return
	}
	auth, err := service.UserLoad(c.Request.Context(), authId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	uploadFile, fh, err := c.Request.FormFile("file")
	if err != nil {
		_ = c.Error(errors.BadRequest("请上传文件", err))
		return
	}
	// 计算上传的文件大小是否超过用户可使用的总大小
	newTotalSize := uint64(fh.Size) + auth.UsedStorage
	if newTotalSize > auth.Group.MaxStorage {
		_ = c.Error(errors.BadRequest("您的空间已经用完啦, 快去求求攻城狮大哥吧 ^_^", err))
		return
	}
	defer uploadFile.Close()
	isExist, err := service.ExistFile(c.Request.Context(), fh.Filename, l.FolderId, authId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if isExist {
		_ = c.Error(errors.FileAlreadyExist("上传失败, 该目录下存在同名文件"))
		return
	} else {
		uFile, err := uf.u.Upload(go_file_uploader.FileHeader{Filename: fh.Filename, Size: fh.Size, File: uploadFile}, "")
		if err != nil {
			_ = c.Error(errors.InternalServerError("上传失败", err))
			return
		}
		var fileModel *model.File
		// hash相同文件名不同, 虽然不用上传文件, 但是需要创建一个不同的folder<->file_name关系
		if uFile.Filename != fh.Filename {
			uFile.Filename = fh.Filename
			fileModel = convert2FileModel(uFile)
		} else {
			fileModel = convert2FileModel(uFile)
		}
		err = service.SaveFileToFolder(c.Request.Context(), fileModel, folder)
		if err != nil {
			_ = c.Error(err)
			return
		}
		// 更新用户已使用的空间
		err = service.UserUpdateUsedStorage(c.Request.Context(), authId, newTotalSize)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusCreated, fileModel)
	}
}

func convert2FileModel(upload *go_file_uploader.FileModel) *model.File {
	return &model.File{
		Id:       upload.Id,
		Hash:     upload.Hash,
		Format:   upload.Format,
		Filename: upload.Filename,
		Size:     upload.Size,
	}
}

func NewUploadFileHandler(u go_file_uploader.Uploader) *uploadFile {
	return &uploadFile{u: u}
}
