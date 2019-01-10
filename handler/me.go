package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/pkg/bytesize"
	"github.com/wq1019/cloud_disk/service"
	"github.com/zm-dev/go-image_uploader/image_url"
	"net/http"
)

type meHandler struct {
	imageUrl image_url.URL
}

func (m *meHandler) Show(c *gin.Context) {
	uid := middleware.UserId(c)
	user, err := service.UserLoad(c.Request.Context(), uid)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, convert2UserResp(user, m.imageUrl))
}

func (m *meHandler) UpdateInfo(c *gin.Context) {
	var authId = middleware.UserId(c)
	l := struct {
		Email      string `json:"email" form:"email"`
		Profile    string `json:"profile" form:"profile"`
		Nickname   string `json:"nickname" form:"nickname"`
		AvatarHash string `json:"avatar_hash" form:"avatar_hash"`
		Gender     int8   `json:"gender" form:"gender"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}
	err := service.UserUpdate(c.Request.Context(), authId, map[string]interface{}{
		"nickname":    l.Nickname,
		"avatar_hash": l.AvatarHash,
		"profile":     l.Profile,
		"email":       l.Email,
		"gender":      l.Gender,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(204, nil)
}

func convert2UserResp(user *model.User, imageUrl image_url.URL) map[string]interface{} {
	var gender string
	if user.Gender {
		gender = "男"
	} else {
		gender = "女"
	}
	return map[string]interface{}{
		"id":                user.Id,
		"name":              user.Name,
		"email":             user.Email,
		"gender":            gender,
		"profile":           user.Profile,
		"nickname":          user.Nickname,
		"created_at":        user.CreatedAt,
		"updated_at":        user.UpdatedAt,
		"avatar_url":        imageUrl.Generate(user.AvatarHash),
		"group_name":        user.Group.Name,
		"avatar_hash":       user.AvatarHash,
		"used_storage":      bytesize.ByteSize(user.UsedStorage),
		"is_allow_share":    user.Group.AllowShare,
		"max_allow_storage": bytesize.ByteSize(user.Group.MaxStorage),
	}
}

func NewMeHandler(imageUrl image_url.URL) *meHandler {
	return &meHandler{imageUrl: imageUrl}
}
