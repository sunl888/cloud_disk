package handler

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/service"
	"github.com/wq1019/go-file-uploader"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	folder, err := service.LoadSimpleFolder(c.Request.Context(), l.FolderId, authId)
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
	// 判断目录是否存在同名文件
	for _, file := range folder.Files {
		if file.Filename == fh.Filename {
			_ = c.Error(errors.FileAlreadyExist("上传失败, 该目录下存在同名文件"))
			return
		}
	}
	// 计算上传的文件大小是否超过用户可使用的总大小
	newTotalSize := uint64(fh.Size) + auth.UsedStorage
	if newTotalSize > auth.Group.MaxStorage {
		_ = c.Error(errors.BadRequest("您的空间已经用完啦, 快去求求攻城狮大哥吧 ^_^", err))
		return
	}
	defer uploadFile.Close()
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

func (uf *uploadFile) Upload(c *gin.Context) {
	l := struct {
		FolderId    int64  `json:"folder_id" form:"folder_id"`
		ChunkIndex  int64  `json:"chunk_index" form:"chunk_index"`
		TotalChunk  int64  `json:"total_chunk" form:"total_chunk"`
		TotalSize   int64  `json:"total_size" form:"total_size"`
		FileHash    string `json:"file_hash" form:"file_hash"`
		IsLastChunk bool   `json:"is_last_chunk" form:"is_last_chunk"`
		Filename    string `json:"filename" form:"filename"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}
	// 验证传入的 chunk 是否合法
	if l.ChunkIndex > l.TotalChunk || l.ChunkIndex < 1 {
		_ = c.Error(errors.BadRequest("chunk 必须大于 0 小于等于 total chunk", nil))
		return
	}
	if l.FolderId == 0 {
		_ = c.Error(errors.BadRequest("请指定上传的文件夹", nil))
		return
	}
	authId := middleware.UserId(c)
	// 判断用户有没有上传到该目录的权限
	folder, err := service.LoadSimpleFolder(c.Request.Context(), l.FolderId, authId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if authId != folder.UserId {
		_ = c.Error(errors.Unauthorized("该目录没有访问权限"))
		return
	}
	// 判断目录是否存在同名文件
	for _, file := range folder.Files {
		if file.Filename == l.Filename {
			_ = c.Error(errors.FileAlreadyExist("上传失败, 该目录下存在同名文件"))
			return
		}
	}
	// 从 form-data 中获取数据块
	postChunkData, _, err := c.Request.FormFile("file-data")
	if err != nil {
		_ = c.Error(errors.BadRequest("请上传文件", err))
		return
	}
	defer postChunkData.Close()

	var (
		tmpDir = os.TempDir()
	)
	tmpFile, err := os.OpenFile(fmt.Sprintf("%s/%s-%d", tmpDir, l.FileHash, l.ChunkIndex), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		_ = c.Error(errors.InternalServerError(fmt.Sprintf("上传第%d个数据块失败: %+v", l.ChunkIndex, err), err))
		return
	}
	_, err = io.Copy(tmpFile, postChunkData)
	if err != nil {
		_ = c.Error(errors.InternalServerError(fmt.Sprintf("上传第%d个数据块失败: %+v", l.ChunkIndex, err), err))
		return
	}
	tmpFile.Close()
	// 如果不是最后一个数据块则到这里就上传完成了
	if l.IsLastChunk == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "上传数据块成功",
		})
		return
	} else if l.IsLastChunk == true && l.TotalChunk == l.ChunkIndex {
		// 合并所有数据块 mode: 0644 - 0022 = 622  这个model 必须是 622 否则在最后验证 hash 的时候会和源文件不一致
		file, err := os.OpenFile(fmt.Sprintf("%s/%s", tmpDir, l.Filename), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			_ = c.Error(errors.InternalServerError(fmt.Sprintf("创建待合并的文件失败:%+v", err), err))
			return
		}
		defer file.Close()
		for i := int64(1); i <= l.TotalChunk; i++ {
			tmpChunkFileName := fmt.Sprintf("%s/%s-%d", tmpDir, l.FileHash, i)
			f, err := os.Open(tmpChunkFileName)
			fBytes, err := ioutil.ReadAll(f)
			if err != nil {
				_ = c.Error(errors.InternalServerError(fmt.Sprintf("读取第%d个数据块失败: %+v", i, err), err))
				return
			}
			_, err = file.WriteAt(fBytes, 0)
			if err != nil {
				_ = c.Error(errors.InternalServerError(fmt.Sprintf("写入第%d个数据块失败: %+v", l.ChunkIndex, err), err))
				return
			}
			// 删除分片文件
			err = os.Remove(tmpChunkFileName)
			if err != nil {
				log.Printf("移除临时文件失败: %+v", err)
			}
		}
		// 获取 fileInfo 信息
		fileStat, err := file.Stat()
		if err != nil {
			_ = c.Error(errors.InternalServerError(fmt.Sprintf("FileInfo 获取失败: %+v", err)))
			return
		}
		md5hash := md5.New()
		if _, err := io.Copy(md5hash, file); err != nil {
			_ = c.Error(errors.InternalServerError(fmt.Sprintf("获取文件 md5 失败: %+v", err)))
			return
		}
		md5sum := fmt.Sprintf("%x", md5hash.Sum(nil))
		if md5sum != l.FileHash {
			_ = c.Error(errors.InternalServerError(fmt.Sprintf("文件 md5 不匹配: %+v", err), err))
			return
		}
		// auth info
		auth, err := service.UserLoad(c.Request.Context(), authId)
		if err != nil {
			_ = c.Error(err)
			return
		}
		// 计算上传的文件大小是否超过用户可使用的总大小
		newTotalSize := uint64(fileStat.Size()) + auth.UsedStorage
		if newTotalSize > auth.Group.MaxStorage {
			_ = c.Error(errors.BadRequest("您的空间已经用完啦, 快去求求攻城狮大哥吧 ^_^", err))
			return
		}
		// 上传到 minio
		uFile, err := uf.u.Upload(go_file_uploader.FileHeader{Filename: fileStat.Name(), Size: fileStat.Size(), File: file}, "")
		if err != nil {
			_ = c.Error(errors.InternalServerError(fmt.Sprintf("上传失败: %+v", err), err))
			return
		}
		// 删除临时文件
		err = os.Remove(file.Name())
		if err != nil {
			log.Printf("移除临时文件失败: %+v", err)
		}
		var fileModel *model.File
		// hash相同文件名不同, 虽然不用上传文件, 但是需要创建一个不同的folder<->file_name关系
		if uFile.Filename != fileStat.Name() {
			uFile.Filename = fileStat.Name()
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
	} else {
		// default
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "上传文件失败",
		})
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
