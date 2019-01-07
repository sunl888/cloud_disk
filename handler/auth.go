package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/service"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type authHandler struct{}

func (authHandler) Login(c *gin.Context) {
	req := &struct {
		Account  string `form:"account" json:"account"`
		Password string `form:"password" json:"password"`
	}{}
	if err := c.ShouldBind(req); err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}
	ticket, err := service.UserLogin(c.Request.Context(), strings.TrimSpace(req.Account), strings.TrimSpace(req.Password))
	if err != nil {
		_ = c.Error(err)
		return
	}
	setAuthCookie(c, ticket.Id, ticket.UserId, int(ticket.ExpiredAt.Sub(time.Now()).Seconds()))
	c.JSON(http.StatusNoContent, nil)
}

func (authHandler) Logout(c *gin.Context) {
	ticketId, err := c.Cookie("ticket_id")
	if err != nil {
		c.JSON(http.StatusNoContent, nil)
		return
	}
	removeAuthCookie(c)
	_ = service.TicketDestroy(c.Request.Context(), ticketId)
	c.JSON(http.StatusNoContent, nil)
}

func (authHandler) Register(c *gin.Context) {
	l := struct {
		Account  string `form:"account" json:"account"`
		Password string `form:"password" json:"password"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(err)
		return
	}
	// 注册账号
	userId, err := service.UserRegister(c.Request.Context(), strings.TrimSpace(l.Account), model.CertificateType(0), l.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}
	// 为新账号添加一个根目录
	err = service.CreateFolder(c.Request.Context(), &model.Folder{
		UserId:     userId,
		Level:      1,
		ParentId:   0,
		Key:        "",
		FolderName: "根目录",
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	err = service.CreateUserInfo(c.Request.Context(), &model.UserInfo{
		UserId:   userId,
		Nickname: l.Account,
		Profile:  "这货很懒,什么都没有说哦",
		IsBan:    false,
		GroupId:  1,
	})
	c.Status(201)
}

func setAuthCookie(c *gin.Context, ticketId string, userId int64, maxAge int) {
	c.SetCookie("ticket_id", ticketId, maxAge, "", "", false, false)
	c.SetCookie("user_id", strconv.FormatInt(userId, 10), maxAge, "", "", false, false)
}

func removeAuthCookie(c *gin.Context) {
	c.SetCookie("ticket_id", "", -1, "", "", false, true)
	c.SetCookie("user_id", "", -1, "", "", false, false)
}

func NewAuthHandler() *authHandler {
	return &authHandler{}
}
