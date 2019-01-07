package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/service"
	"github.com/zm-dev/go-file-uploader"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type uploadFile struct {
	u go_file_uploader.Uploader
}

func copy2TmpFile(file io.Reader) (tmpFileName string, err error) {
	tmpFile, err := ioutil.TempFile("", "cloud-")
	_, err = io.Copy(tmpFile, file)
	if cerr := tmpFile.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		_ = os.Remove(tmpFile.Name())
		return
	}
	return tmpFile.Name(), nil
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
	auth, err := service.UserLoadAndRelated(c.Request.Context(), authId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	uploadFile, fh, err := c.Request.FormFile("file")
	if err != nil {
		_ = c.Error(errors.BadRequest("请上传文件", err))
		return
	}
	newTotalSize := fh.Size + auth.UserInfo.UsedStorage
	if newTotalSize > auth.Group.MaxStorage {
		_ = c.Error(errors.BadRequest("您的空间已经用完啦, 开会员吧", err))
		return
	}
	// 写入最新已用容量
	err = service.UpdateUsedStorage(c.Request.Context(), authId, newTotalSize)
	if err != nil {
		_ = c.Error(err)
		return
	}
	defer uploadFile.Close()
	tmpFileName, err := copy2TmpFile(uploadFile)
	if err != nil {
		_ = c.Error(err)
		// 还原最新已用容量
		err = service.UpdateUsedStorage(c.Request.Context(), authId, auth.UserInfo.UsedStorage)
		if err != nil {
			_ = c.Error(err)
			return
		}
		return
	}
	defer os.Remove(tmpFileName)
	uFile, err := uf.u.Upload(go_file_uploader.FileHeader{Filename: fh.Filename, Size: fh.Size, File: uploadFile}, "")
	if err != nil {
		_ = c.Error(errors.InternalServerError("上传失败", err))
		// 还原最新已用容量
		err = service.UpdateUsedStorage(c.Request.Context(), authId, auth.UserInfo.UsedStorage)
		if err != nil {
			_ = c.Error(err)
			return
		}
		return
	}
	fileModel := convert2FileModel(uFile)
	err = service.SaveFileToFolder(c.Request.Context(), fileModel, folder)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, fileModel)
}

func (uf *uploadFile) ShowFile(c *gin.Context) {
	hash := strings.TrimSpace(c.Param("hash"))
	s := uf.u.Store()
	fileModel, err := s.FileLoad(hash)
	if s.FileIsNotExistError(err) {
		_ = c.Error(errors.NotFound("文件不存在"))
		return
	}
	c.JSON(200, fileModel)
}

func convert2FileModel(upload *go_file_uploader.FileModel) *model.File {
	return &model.File{
		Hash:      upload.Hash,
		Format:    upload.Format,
		Filename:  upload.Filename,
		Size:      upload.Size,
		CreatedAt: upload.CreatedAt,
		UpdatedAt: upload.UpdatedAt,
	}
}

func NewUploadFileHandler(u go_file_uploader.Uploader) *uploadFile {
	return &uploadFile{u: u}
}
