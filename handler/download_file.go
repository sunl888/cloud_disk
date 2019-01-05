package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/service"
	"github.com/zm-dev/go-file-uploader"
	"io"
	"strconv"
	"strings"
)

type downloadHandler struct {
	u go_file_uploader.Uploader
}

func download(c *gin.Context, u go_file_uploader.Uploader, userId, fileId, folderId int64, ch chan int) (err error) {
	file, err := service.LoadFile(c.Request.Context(), folderId, fileId, userId)
	if err != nil {
		return err
	}
	r, err := u.ReadFile(file.Hash)
	if err != nil {
		return err
	}
	defer r.Close()

	c.Writer.Header().Add("Content-Disposition", "attachment;filename="+file.Filename)
	_, err = io.Copy(c.Writer, r)
	if err != nil {
		return err
	}
	defer func() {
		<-ch
	}()
	return
}

func (d *downloadHandler) DownloadFile(c *gin.Context) {
	folderIdStr,_ :=c.GetQuery("current_folder_id")
	currentFolderId, err := strconv.ParseInt(strings.TrimSpace(folderIdStr), 10, 64)
	if err != nil || currentFolderId <= 0 {
		_ = c.Error(errors.BadRequest("请指定当前目录ID"))
		return
	}
	fileIds, _ := c.GetQueryArray("file_ids[]")
	folderIds, _ := c.GetQueryArray("folder_ids[]")
	if len(fileIds) == 0 && len(folderIds) == 0 {
		_ = c.Error(errors.BadRequest("请指定要下载的文件或者目录ID"))
		return
	}
	var (
		folderIds2Int64 []int64
		fileIds2Int64   []int64
	)
	for _, v := range folderIds {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		folderIds2Int64 = append(folderIds2Int64, id)
	}
	for _, v := range fileIds {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		fileIds2Int64 = append(fileIds2Int64, id)
	}

	ch := make(chan int, 4)
	authId := middleware.UserId(c)
	if len(folderIds2Int64) > 0 {
		folderFiles, err := service.LoadFileIds(c.Request.Context(), folderIds2Int64, authId)
		if err != nil {
			_ = c.Error(err)
			return
		}
		for _, v := range folderFiles {
			ch <- 1
			go download(c, d.u, authId, v.FileId, v.FolderId, ch)
		}
	}
	if len(fileIds2Int64) > 0 {
		for _, v := range fileIds2Int64 {
			ch <- 1
			go download(c, d.u, authId, v, currentFolderId, ch)
		}
	}
	c.Status(200)
}

func NewDownloadHandler(u go_file_uploader.Uploader) *downloadHandler {
	return &downloadHandler{u: u}
}
