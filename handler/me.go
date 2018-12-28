package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/service"
	"net/http"
)

func convert2UserResp(user *model.User) map[string]interface{} {
	return map[string]interface{}{
		"id":         user.Id,
		"name":       user.Name,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
}

type meHandler struct {
}

func (*meHandler) Show(c *gin.Context) {
	uid := middleware.UserId(c)
	user, err := service.UserLoad(c.Request.Context(), uid)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, convert2UserResp(user))
}

func NewMeHandler() *meHandler {
	return &meHandler{}
}
