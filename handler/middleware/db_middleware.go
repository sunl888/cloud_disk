package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/store/db_store"
)

func Gorm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(db_store.NewDBContext(c.Request.Context(), db))
		c.Next()
	}
}
