package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/service"
	"github.com/zm-dev/go-file-uploader"
	"io"
)

type downloadHandler struct {
	u go_file_uploader.Uploader
}

func (d *downloadHandler) DownloadFile(c *gin.Context) {
	l := struct {
		FileId   int64 `json:"file_id" form:"file_id"`
		FolderId int64 `json:"folder_id" form:"folder_id"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(err)
		return
	}
	authId := middleware.UserId(c)
	file, err := service.LoadFile(c.Request.Context(), l.FolderId, l.FileId, authId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	r, err := d.u.ReadFile(file.Hash)
	if err != nil {
		_ = c.Error(err)
		return
	}
	defer r.Close()

	c.Writer.Header().Add("Content-Disposition", "attachment;filename="+file.Filename)
	_, err = io.Copy(c.Writer, r)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(200)
}

func generageUrl() {

}

func NewDownloadHandler(u go_file_uploader.Uploader) *downloadHandler {
	return &downloadHandler{u: u}
}
