package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/service"
	"github.com/zm-dev/go-image_uploader/image_url"
	"strconv"
)

type userHandler struct {
	imageUrl image_url.URL
}

func (m *userHandler) UpdateBanStatus(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}
	if userId <= 0 {
		_ = c.Error(errors.ErrAccountNotFound())
		return
	}
	l := struct {
		IsBan bool `json:"is_ban" form:"is_ban"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}
	user, err := service.UserLoad(c.Request.Context(), userId)
	if user.IsAdmin == true {
		_ = c.Error(errors.UserNotAllowBeBan("管理员账号不允许被 ban"))
		return
	}
	err = service.UserUpdateBanStatus(c.Request.Context(), userId, l.IsBan)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(204)
}

func NewUserHandler(imageUrl image_url.URL) *userHandler {
	return &userHandler{imageUrl: imageUrl}
}
