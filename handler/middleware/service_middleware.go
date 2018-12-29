package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/service"
)

var ServiceKey = "service"

func SetService(c *gin.Context, s service.Service) {
	c.Set(ServiceKey, s)
}

func ServiceMiddleware(s service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		SetService(c, s)
		c.Next()
	}
}
