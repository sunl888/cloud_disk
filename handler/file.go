package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/service"
	"net/http"
)

type fileHandler struct{}

// RenameFile godoc
// @Tags 文件
// @Summary 重命名文件
// @Description 通过文件 ID 重命名文件
// @ID rename-file
// @Accept json,multipart/form-data
// @Produce json,multipart/form-data
// @Param file_id query uint64 true "文件 ID" Format(uint64)
// @Param folder_id query uint64 true "文件所属的目录 ID" Format(uint64)
// @Param new_name query string true "新的文件名" Format(string)
// @Success 204
// @Failure 404 {object} errors.GlobalError "文件不存在" | "目录不存在"
// @Failure 500 {object} errors.GlobalError
// @Router /file/rename [PUT]
func (*fileHandler) RenameFile(c *gin.Context) {
	l := struct {
		FileId   int64  `json:"file_id" form:"file_id"`
		FolderId int64  `json:"folder_id" form:"folder_id"`
		NewName  string `json:"new_name" form:"new_name"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(err)
		return
	}
	authId := middleware.UserId(c)
	folder, err := service.LoadFolder(c.Request.Context(), l.FolderId, authId, false)
	if err != nil {
		_ = c.Error(err)
		return
	}
	err = service.RenameFile(c.Request.Context(), folder.Id, l.FileId, l.NewName)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func NewFileHandler() *fileHandler {
	return &fileHandler{}
}
