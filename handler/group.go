package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/service"
	"net/http"
	"strconv"
)

type groupHandler struct {
}

func (g *groupHandler) GroupCreate(c *gin.Context) {
	l := struct {
		Name       string `json:"name" form:"name"`
		MaxStorage uint64 `json:"max_storage" form:"max_storage"`
		AllowShare bool   `json:"allow_share" form:"allow_share"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}
	group := model.Group{
		Name:       l.Name,
		MaxStorage: l.MaxStorage,
		AllowShare: l.AllowShare,
	}
	err := service.GroupCreate(c.Request.Context(), &group)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, group)
}

func (g *groupHandler) GroupUpdate(c *gin.Context) {
	groupId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}
	if groupId <= 0 {
		_ = c.Error(model.ErrGroupNotExist)
		return
	}
	l := struct {
		Name       string `json:"name" form:"name"`
		MaxStorage uint64 `json:"max_storage" form:"max_storage"`
		AllowShare bool   `json:"allow_share" form:"allow_share"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}
	err = service.GroupUpdate(c.Request.Context(), groupId, map[string]interface{}{
		"name":        l.Name,
		"max_storage": l.MaxStorage,
		"allow_share": l.AllowShare,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusCreated)
}

func (g *groupHandler) GroupList(c *gin.Context) {
	limit, offset := getInt64LimitAndOffset(c)
	groups, count, err := service.GroupList(c.Request.Context(), offset, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(200, gin.H{
		"count": count,
		"data":  groups,
	})
}

func (g *groupHandler) GroupDelete(c *gin.Context) {
	groupId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}
	err = service.GroupDelete(c.Request.Context(), groupId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(204)
}

func NewGroupHandler() *groupHandler {
	return &groupHandler{}
}
