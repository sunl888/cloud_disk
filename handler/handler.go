package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/wq1019/cloud_disk/handler/middleware"
	"github.com/wq1019/cloud_disk/server"
	"net/http"
	"strconv"
)

func CreateHTTPHandler(s *server.Server) http.Handler {
	authHandler := NewAuthHandler()
	meHandler := NewMeHandler(s.ImageUrl)
	userHandler := NewUserHandler(s.ImageUrl)
	uploadFileHandler := NewUploadFileHandler(s.FileUploader)
	uploadImageHandler := NewUploadImage(s.ImageUploader, s.ImageUrl)
	folderHandler := NewFolderHandler()
	fileHandler := NewFileHandler()
	downloadHandler := NewDownloadHandler(s.FileUploader)

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
	authRouter.POST("/register", authHandler.Register)
	authRouter.POST("/login", authHandler.Login)

	authorized := api.Group("/")
	authorized.Use(middleware.AuthMiddleware)
	{
		// 显示我的基本信息
		authorized.GET("/auth/me", meHandler.Show)
		// 更新我的基本信息
		authorized.PUT("/auth/me", meHandler.UpdateInfo)
		// 退出登录
		authorized.GET("/auth/logout", authHandler.Logout)
		// 上传文件
		authorized.POST("/upload_file", uploadFileHandler.UploadFile)
		// 上传图片
		authorized.POST("/upload_image", uploadImageHandler.UploadImage)
		// 指定目录下第一层的资源列表
		authorized.GET("/folder", folderHandler.LoadFolder)
		// 创建目录
		authorized.POST("/folder", folderHandler.CreateFolder)
		// 删除文件和目录资源
		authorized.DELETE("/source", folderHandler.DeleteSource)
		// 移动到指定目录
		authorized.PUT("/source/move", folderHandler.Move2Folder)
		// 复制到指定目录
		authorized.PUT("/source/copy", folderHandler.Copy2Folder)
		// 重命名文件
		authorized.PUT("/file/rename", fileHandler.RenameFile)
		// 重命名目录
		authorized.PUT("/folder/rename", folderHandler.RenameFolder)
		// 文件下载
		authorized.GET("/download", downloadHandler.Download)
		// 获取要下载的文件和目录的详细信息
		authorized.GET("/pre_download", downloadHandler.PreDownload)
	}

	adminRouter := api.Group("/admin")
	adminRouter.Use(middleware.AuthMiddleware, middleware.AdminMiddleware)
	{
		// 用户列表
		adminRouter.GET("/user", userHandler.UserList)
		// 更新用户的禁用状态
		adminRouter.PUT("/user/:id/ban_status", userHandler.UpdateBanStatus)
	}

	// 文档
	api.GET("/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}

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
