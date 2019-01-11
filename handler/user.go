package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/pkg/bytesize"
	"github.com/wq1019/cloud_disk/service"
	"github.com/zm-dev/go-image_uploader/image_url"
	"strconv"
)

type userHandler struct {
	imageUrl image_url.URL
}

func (*userHandler) UpdateBanStatus(c *gin.Context) {
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

func (u *userHandler) UserList(c *gin.Context) {
	limit, offset := getInt64LimitAndOffset(c)
	users, err := service.UserList(c.Request.Context(), offset, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(200, convert2UserListResp(users, u.imageUrl))
}

func convert2UserListResp(users []*model.User, imageUrl image_url.URL) []map[string]interface{} {
	userList := make([]map[string]interface{}, 0, len(users))
	for _, v := range users {
		userList = append(userList, convert2UserResp(v, imageUrl))
	}
	return userList
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

func NewUserHandler(imageUrl image_url.URL) *userHandler {
	return &userHandler{imageUrl: imageUrl}
}
