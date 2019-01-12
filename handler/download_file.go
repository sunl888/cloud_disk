package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/service"
	uploader "github.com/wq1019/go-file-uploader"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type downloadHandler struct {
	u uploader.Uploader
}

type FolderData struct {
	Filename string
	Key      string
}

func (d *downloadHandler) PreDownload(c *gin.Context) {
	folderIdStr, _ := c.GetQuery("current_folder_id")
	currentFolderId, err := strconv.ParseInt(strings.TrimSpace(folderIdStr), 10, 64)
	if err != nil || currentFolderId <= 0 {
		_ = c.Error(errors.BadRequest("请指定当前目录ID"))
		return
	}
	// 选中的文件
	fileIdsReq, _ := c.GetQueryArray("file_ids[]")
	// 选中的目录
	folderIdsReq, _ := c.GetQueryArray("folder_ids[]")
	if len(fileIdsReq) == 0 && len(folderIdsReq) == 0 {
		_ = c.Error(errors.BadRequest("请指定要下载的文件或者目录ID"))
		return
	}
	var (
		authId          = middleware.UserId(c)
		fileIds2Int64   = strArr2Int64Arr(fileIdsReq)
		folderIds2Int64 = strArr2Int64Arr(folderIdsReq)
		foldersLen      = len(folderIds2Int64)
		filesLen        = len(fileIds2Int64)
		folderFiles     = make([]*model.WrapFolderFile, 0, foldersLen+filesLen)
	)
	// 对于用户指定的所有目录下的文件都要查出来并返回
	if foldersLen > 0 {
		folderFiles, err = service.LoadFolderFilesByFolderIds(c.Request.Context(), folderIds2Int64, authId)
		if err != nil {
			_ = c.Error(err)
			return
		}
	}
	// 用户明确选中需要下载的文件(注: 就是当前目录下用户选中的文件)
	currentFolderFiles, err := service.LoadFolderFilesByFolderIdAndFileIds(c.Request.Context(), currentFolderId, fileIds2Int64, authId)
	if len(fileIds2Int64) > 0 {
		for _, v := range currentFolderFiles {
			folderFiles = append(folderFiles, v)
		}
	}
	if len(folderFiles) == 0 {
		_ = c.Error(errors.BadRequest("没有要下载的文件"))
		return
	}

	// Wrap Response Data
	var (
		folderIds  = make([]int64, 0, 5)
		folderMaps = make(map[int64]FolderData, 10)
	)
	// 查找每个文件所在目录的信息
	folderIds = append(folderIds, currentFolderId)
	for _, v := range folderFiles {
		folderIds = append(folderIds, v.FolderId)
	}
	folders, err := service.ListFolder(c.Request.Context(), folderIds, authId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	// 将目录的 id 与 name 写入 Map
	for _, v := range folders {
		folderMaps[v.Id] = FolderData{
			Filename: v.FolderName,
			Key:      v.Key,
		}
	}
	for i := 0; i < len(folderFiles); i++ {
		relativePath := mergePath(folderMaps, currentFolderId, folderMaps[folderFiles[i].FolderId].Key, folderFiles[i].FolderId)
		folderFiles[i].RelativePath = relativePath
	}
	c.JSON(http.StatusOK, folderFiles)
}

// 文件下载
// example:
// curl -H "Range: bytes=0-12929" http://localhost:8080/api/download?folder_id=1\&file_id=3 -v --output 3.png
func (d *downloadHandler) Download(c *gin.Context) {
	l := struct {
		FolderId int64 `json:"folder_id" form:"folder_id"`
		FileId   int64 `json:"file_id" form:"file_id"`
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
	readFile, err := d.u.ReadFile(file.Hash)
	if err != nil {
		_ = c.Error(err)
		return
	}
	defer readFile.Close()
	c.Writer.Header().Add("Content-Disposition", "attachment;filename="+file.Filename)
	// Range: bytes=start-end
	rangeVal := c.Request.Header.Get("Range")
	if rangeVal != "" {
		var (
			start       int64 // not null
			end         int64 // not null
			prefixIndex = strings.Index(rangeVal, "-")
		)
		if prefixIndex == -1 {
			_ = c.Error(errors.BadRequest("Http range header error, not found prefix `-`"))
			return
		}
		start, err = strconv.ParseInt(rangeVal[6:prefixIndex], 10, 64)
		if err != nil {
			_ = c.Error(err)
			return
		}
		if rangeVal[prefixIndex+1:] == "" {
			_ = c.Error(errors.BadRequest("Http range header error, end value not exist."))
			return
		}
		end, err = strconv.ParseInt(rangeVal[prefixIndex+1:], 10, 64)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.Writer.Header().Add("Accept-Ranges", "bytes")
		c.Writer.Header().Add("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, file.Size))
		// 这里必须要提前设置状态吗为 206 否则会 Warning https://github.com/gin-gonic/gin/issues/471#issuecomment-190186203
		c.Status(http.StatusPartialContent)

		buff := make([]byte, end-start)
		_, err = readFile.ReadAt(buff, start)
		if err != nil {
			_ = c.Error(err)
			return
		}
		_, err = c.Writer.Write(buff)
		if err != nil {
			_ = c.Error(err)
			return
		}
	} else {
		// 整个文件下载
		_, err = io.Copy(c.Writer, readFile)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.Status(http.StatusOK)
	}
}

func mergePath(folderMap map[int64]FolderData, currentId int64, key string, withId int64) (path string) {
	if currentId == withId {
		return "./"
	}
	key2Arr := strings.Split(key, "-")
	for _, v := range key2Arr {
		id2Int64, _ := strconv.ParseInt(v, 10, 64)
		if id2Int64 > currentId {
			path += folderMap[id2Int64].Filename + "/"
		}
	}
	path += folderMap[withId].Filename
	return path
}

func strArr2Int64Arr(str []string) []int64 {
	var int64Arr []int64
	for _, v := range str {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		int64Arr = append(int64Arr, id)
	}
	return int64Arr
}

func NewDownloadHandler(u uploader.Uploader) *downloadHandler {
	return &downloadHandler{u: u}
}
