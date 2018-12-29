package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/service"
)

var (
	isLoginKey    = "is_login"
	userIdKey     = "user_id"
	loggedUserKey = "logged_user"
)

func AuthMiddleware(c *gin.Context) {
	// TODO authId
	isLogin := true
	//isLogin := check(c)
	if !isLogin {
		_ = c.Error(errors.Unauthorized())
		c.Abort()
		return
	}
	c.Next()
}

func AdminMiddleware(c *gin.Context) {
	user := LoggedUser(c)
	if user == nil || !user.IsAdmin {
		_ = c.Error(errors.Forbidden("没有权限.", nil))
		c.Abort()
		return
	}
	c.Next()
}

func check(c *gin.Context) bool {
	var (
		isLogin bool
	)
	if ticketId, err := c.Cookie("ticket_id"); err == nil {
		isValid, userId, err := service.TicketIsValid(c.Request.Context(), ticketId)
		if err == nil {
			isLogin = isValid
			setIsLogin(c, isLogin)
			setUserId(c, userId)
		}
	} else {
		// cookie不存在
		isLogin = false
	}
	return isLogin
}

func setIsLogin(c *gin.Context, isLogin bool) {
	c.Set(isLoginKey, isLogin)
}

func setUserId(c *gin.Context, userId int64) {
	c.Set(userIdKey, userId)
}

func CheckLogin(c *gin.Context) bool {
	isLogin, ok := c.Get(isLoginKey)
	if !ok {
		return check(c)
	}
	return isLogin.(bool)

}

func UserId(c *gin.Context) int64 {
	userId, ok := c.Get(userIdKey)
	if !ok {
		check(c)
		return c.GetInt64(userIdKey)
	}
	return userId.(int64)
}

func LoggedUser(c *gin.Context) *model.User {
	user, ok := c.Get(loggedUserKey)
	if !ok {
		userId := UserId(c)
		if userId == 0 {
			return nil
		}
		userModel, err := service.UserLoad(c.Request.Context(), userId)
		if err != nil {
			return nil
		}
		c.Set("loggedUserKey", userModel)
		return userModel
	}
	return user.(*model.User)
}
