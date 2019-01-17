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

type FormData struct {
	FolderId    int64  `json:"folder_id" form:"folder_id"`
	ChunkIndex  int    `json:"chunk_index" form:"chunk_index"`
	TotalChunk  int    `json:"total_chunk" form:"total_chunk"`
	TotalSize   int64  `json:"total_size" form:"total_size"`
	FileHash    string `json:"file_hash" form:"file_hash"`
	IsLastChunk bool   `json:"is_last_chunk" form:"is_last_chunk"`
	Filename    string `json:"filename" form:"filename"`
	UploadId    string `json:"upload_id" form:"upload_id"`
}

type UploadResponse struct {
	UploadId   string `json:"upload_id"`
	ChunkIndex int    `json:"chunk_index"`
	FileHash   string `json:"file_hash"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

const (
	ChunkMaxSize = 100 << 20 // 分片上传最大 100MB
)

func (uf *uploadFile) UploadChunk(c *gin.Context) {
	l := FormData{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}
	// 验证表单
	if ok, err := validForm(&l); !ok {
		_ = c.Error(err)
		return
	}
	var (
		authId = middleware.UserId(c) // 没必要每次都获取 authID, 第一次上传和最后一次上传时获取一下就可以
		err    error                  // err
	)
	// 第一次上传或者最后一次上传时都检查有没有权限
	if l.ChunkIndex == 1 || l.IsLastChunk == true {
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
	}
	// 从 form-data 中获取数据块
	postChunkData, fh, err := c.Request.FormFile("file-data")
	if err != nil {
		_ = c.Error(errors.BadRequest("请上传文件", err))
		return
	}
	defer postChunkData.Close()
	if fh.Size > ChunkMaxSize {
		_ = c.Error(errors.BadRequest("上传失败, 数据块太大"))
		return
	}
	// 文件秒传
	exist, _ := uf.u.Store().FileExist(l.FileHash)
	if exist {
		file, err := uf.u.Store().FileLoad(l.FileHash)
		if err != nil {
			_ = c.Error(errors.BadRequest("文件秒传失败", err))
			return
		}
		fileModel := &model.File{}
		if file.Filename != l.Filename {
			file.Filename = l.Filename
			fileModel = convert2FileModel(file)
		} else {
			fileModel = convert2FileModel(file)
		}
		err = service.SaveFileToFolder(c.Request.Context(), fileModel, l.FolderId)
		if err != nil {
			_ = c.Error(err)
			return
		}
		// 更新用户已使用的空间
		err = service.UserUpdateUsedStorage(c.Request.Context(), authId, uint64(file.Size))
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, UploadResponse{
			Message:    "文件秒传成功",
			StatusCode: 1,
		})
		return
	} else {
		// 分片上传
		file, uploadId, err := uf.u.UploadChunk(go_file_uploader.ChunkHeader{
			ChunkNumber:    l.ChunkIndex,
			UploadId:       l.UploadId,
			OriginFilename: l.Filename,
			OriginFileHash: l.FileHash,
			OriginFileSize: l.TotalSize,
			IsLastChunk:    l.IsLastChunk,
			ChunkContent:   postChunkData,
			ChunkCount:     l.TotalChunk,
		}, "")
		if err != nil {
			_ = c.Error(errors.BadRequest("分片上传失败", err))
			return
		}
		// 非最后一个数据块
		if l.IsLastChunk == false {
			c.JSON(http.StatusOK, UploadResponse{
				Message:    "数据块上传成功",
				StatusCode: 2,
				UploadId:   uploadId,
				ChunkIndex: l.ChunkIndex,
				FileHash:   l.FileHash,
			})
			return
		} else {
			if file == nil {
				_ = c.Error(errors.BadRequest("文件上传失败, 所有数据块已经上传, 但是保存到数据库时可能出现了问题"))
				return
			}
			// 最后一个数据块上传完成后需要写入文件和目录的关系到数据库
			fileModel := &model.File{}
			if file.Filename != l.Filename {
				file.Filename = l.Filename
				fileModel = convert2FileModel(file)
			} else {
				fileModel = convert2FileModel(file)
			}
			err = service.SaveFileToFolder(c.Request.Context(), fileModel, l.FolderId)
			if err != nil {
				_ = c.Error(err)
				return
			}
			// 更新用户已使用的空间
			err = service.UserUpdateUsedStorage(c.Request.Context(), authId, uint64(file.Size))
			if err != nil {
				_ = c.Error(err)
				return
			}
			c.JSON(http.StatusOK, UploadResponse{
				Message:    "文件上传成功",
				StatusCode: 0,
				UploadId:   uploadId,
				ChunkIndex: l.ChunkIndex,
				FileHash:   l.FileHash,
			})
			return
		}
	}
}

func validForm(l *FormData) (ok bool, err error) {
	ok = false
	err = nil
	if l.Filename == "" {
		err = errors.BadRequest("filename 不存在", nil)
		return
	}
	if l.FileHash == "" {
		err = errors.BadRequest("filehash 不存在", nil)
		return
	}
	if l.TotalSize <= 0 {
		err = errors.BadRequest("totalSize 必须大于 0", nil)
		return
	}
	// 验证传入的 chunk 是否合法
	if l.ChunkIndex > l.TotalChunk || l.ChunkIndex < 1 {
		err = errors.BadRequest("chunk 必须大于 0 小于等于 totalChunk", nil)
		return
	}
	if l.FolderId == 0 {
		err = errors.BadRequest("请指定上传的文件夹", nil)
		return
	}
	if l.ChunkIndex != 1 && l.UploadId == "" {
		err = errors.BadRequest("除第一次上传文件, 每次上传都要传 uploadId")
		return
	}
	return true, nil
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
