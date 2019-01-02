package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/server"
	"net/http"
	"strconv"
)

func getInt32LimitAndOffset(c *gin.Context) (limit, offset int32) {
	var err error
	limitI64, err := strconv.ParseInt(c.Query("limit"), 10, 32)
	if err != nil {
		limit = 10
	} else {
		limit = int32(limitI64)
	}
	if limit > 50 {
		limit = 50
	}

	offsetI64, err := strconv.ParseInt(c.Query("offset"), 10, 32)
	if err != nil {
		offset = 0
	} else {
		offset = int32(offsetI64)
	}
	return limit, offset
}

func getInt64LimitAndOffset(c *gin.Context) (limit, offset int64) {
	var err error
	limit, err = strconv.ParseInt(c.Query("limit"), 10, 32)
	if err != nil {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	offset, err = strconv.ParseInt(c.Query("offset"), 10, 32)
	if err != nil {
		offset = 0
	}
	return limit, offset
}

func CreateHTTPHandler(s *server.Server) http.Handler {
	authHandler := NewAuthHandler()
	meHandler := NewMeHandler()
	uploadFileHandler := NewUploadFileHandler(s.FileUploader)
	folderHandler := NewFolderHandler()
	fileHandler := NewFileHandler()

	if s.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(middleware.Gorm(s.DB))
	router.Use(middleware.Service(s.Service))
	router.Use(middleware.NewHandleErrorMiddleware(s.Conf.ServiceName))
	api := router.Group("/api")

	authRouter := api.Group("/auth")
	// 注册
	authRouter.POST("/register", authHandler.Register)
	// 登录
	authRouter.POST("/login", authHandler.Login)
	{
		authRouter.GET("/logout", authHandler.Logout).Use(middleware.AuthMiddleware)
		authRouter.GET("/me", meHandler.Show).Use(middleware.AuthMiddleware)
	}

	authorized := api.Group("/")
	authorized.Use(middleware.AuthMiddleware)
	{
		// 上传文件
		authorized.POST("/upload_file", uploadFileHandler.UploadFile)
		// 指定目录下第一层的资源列表
		authorized.GET("/folder", folderHandler.LoadFolder)
		// 创建目录
		authorized.POST("/folder", folderHandler.CreateFolder)
		// 删除文件和目录资源 (file_ids, folder_ids)
		authorized.DELETE("/source", folderHandler.DeleteSource)
		// TODO 如何避免往子目录移动,怎样复制到当前目录(未完成)
		// 移动到指定目录
		authorized.PUT("/source/move", folderHandler.Move2Folder)
		// 复制到指定目录
		authorized.PUT("/source/copy", folderHandler.Copy2Folder)
		// 重命名文件
		authorized.PUT("/file/rename", fileHandler.RenameFile)
		// 重命名目录
		authorized.PUT("/folder/rename", folderHandler.RenameFolder)
	}

	adminRouter := api.Group("/")
	adminRouter.Use(middleware.AuthMiddleware, middleware.AdminMiddleware)
	{
	}
	return router
}
