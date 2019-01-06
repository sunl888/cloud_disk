package handler

import (
	"archive/zip"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/service"
	uploader "github.com/zm-dev/go-file-uploader"
	"io"
	"strconv"
	"strings"
)

type downloadHandler struct {
	u uploader.Uploader
}

const (
	MaxCount = 1000    // 最大 1000 个文件
	MaxSize  = 2 << 30 // 1024*1024*1024*2 = 2^31 = 2GB = 2147483648Byte
)

type FolderData struct {
	Filename string
	Key      string
}

type FileData struct {
	Id       int64
	FolderId int64
	Filename string
	Hash     string
	Size     int64
}

// 批量打包下载文件
func downloadMultiple(c *gin.Context, u uploader.Uploader, userId, currentFolderId int64, folderFiles []*model.FolderFile) (err error) {
	var (
		size       int64
		count      int64
		folderIds  = make([]int64, 0, 5)
		folderMaps = make(map[int64]FolderData, 10)
		fileLists  = make([]*FileData, 10)
	)
	// 查找每个文件所在目录的信息
	folderIds = append(folderIds, currentFolderId)
	for _, v := range folderFiles {
		folderIds = append(folderIds, v.FolderId)
	}
	folders, err := service.ListFolder(c.Request.Context(), folderIds, userId)
	if err != nil {
		return err
	}
	// 将目录的 id 与 name 写入 Map
	for _, v := range folders {
		folderMaps[v.Id] = FolderData{
			Filename: v.FolderName,
			Key:      v.Key,
		}
	}

	// 查找所有文件的详细信息
	for _, v := range folderFiles {
		file, err := service.LoadFile(c.Request.Context(), v.FolderId, v.FileId, userId)
		if err != nil {
			return err
		}
		count++
		size += file.Size
		if count > MaxCount {
			return errors.BadRequest(fmt.Sprintf("文件数量超过%d个，不允许下载", MaxCount))
		} else if size > MaxSize {
			return errors.BadRequest(fmt.Sprintf("文件总大小超过%dByte，不允许下载", MaxSize))
		}
		fileLists = append(fileLists, &FileData{
			Id:       v.FileId,
			FolderId: v.FolderId,
			Filename: file.Filename,
			Hash:     file.Hash,
			Size:     file.Size,
		})
	}

	// 将文件全部写入 Zip 文件流中
	c.Writer.Header().Add("Content-Disposition", "attachment;filename=批量下载.zip")
	w := zip.NewWriter(c.Writer)
	defer w.Close()
	for _, file := range fileLists {
		rFile, err := u.ReadFile(file.Hash)
		if err != nil {
			return err
		}
		path := generatePath(folderMaps, currentFolderId, folderMaps[file.FolderId].Key, file.FolderId)
		err = compress(rFile, path, w, file.Filename)
		rFile.Close()
	}

	return nil
}

func download(c *gin.Context, u uploader.Uploader, userId, fileId, folderId int64) (err error) {
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
	return
}

func compress(file uploader.ReadFile, prefix string, zw *zip.Writer, filename string) error {
	writer, err := zw.CreateHeader(&zip.FileHeader{
		Name: prefix + "/" + filename,
	})
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}
	return nil
}

func generatePath(folderMap map[int64]FolderData, currentId int64, key string, ownId int64) (path string) {
	if currentId == ownId {
		return ""
	}
	key2Arr := strings.Split(key, "-")
	for _, v := range key2Arr {
		id2Int64, _ := strconv.ParseInt(v, 10, 64)
		if id2Int64 > currentId {
			path += folderMap[id2Int64].Filename + "/"
		}
	}
	path += folderMap[ownId].Filename
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

func (d *downloadHandler) DownloadFile(c *gin.Context) {
	folderIdStr, _ := c.GetQuery("current_folder_id")
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
	fileIds2Int64 := strArr2Int64Arr(fileIds)
	folderIds2Int64 := strArr2Int64Arr(folderIds)
	// 合并目录和文件
	ffs := make([]*model.FolderFile, 0, 10)
	authId := middleware.UserId(c)
	if len(folderIds2Int64) > 0 {
		ffs, err = service.LoadFilesByFolderIds(c.Request.Context(), folderIds2Int64, authId)
		if err != nil {
			_ = c.Error(err)
			return
		}
	}
	// 当前目录中要下载的文件
	if len(fileIds2Int64) > 0 {
		for _, v := range fileIds2Int64 {
			ffs = append(ffs, &model.FolderFile{
				FileId:   v,
				FolderId: currentFolderId,
			})
		}
	}
	if len(ffs) == 0 {
		_ = c.Error(errors.BadRequest("没有要下载的文件"))
		return
	}
	if len(ffs) == 1 {
		// 下载单个文件
		download(c, d.u, authId, ffs[0].FileId, ffs[0].FolderId)
	} else {
		// 多个文件打包下载
		err = downloadMultiple(c, d.u, authId, currentFolderId, ffs)
		if err != nil {
			_ = c.Error(err)
			return
		}
	}
	c.Status(200)
}

func NewDownloadHandler(u uploader.Uploader) *downloadHandler {
	return &downloadHandler{u: u}
}
