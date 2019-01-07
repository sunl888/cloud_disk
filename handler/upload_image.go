package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/zm-dev/go-image_uploader"
	"github.com/zm-dev/go-image_uploader/image_url"
)

type uploadImage struct {
	u        image_uploader.Uploader
	imageUrl image_url.URL
}

func (ui *uploadImage) UploadImage(c *gin.Context) {
	file, fh, err := c.Request.FormFile("image")
	if err != nil {
		_ = c.Error(errors.BadRequest("请上传图片", err))
		return
	}
	defer file.Close()
	image, err := ui.u.Upload(image_uploader.FileHeader{Filename: fh.Filename, Size: fh.Size, File: file})

	if err != nil {
		if image_uploader.IsUnknownFormat(err) {
			_ = c.Error(errors.BadRequest("不支持的图片类型", nil))
			return
		} else {
			_ = c.Error(errors.InternalServerError("图片上传失败", err))
			return
		}
	}
	u := ui.imageUrl.Generate(image.Hash)

	c.JSON(200, gin.H{
		"image_url":  u,
		"image_hash": image.Hash,
	})
}

func NewUploadImage(u image_uploader.Uploader, imageUrl image_url.URL) *uploadImage {
	return &uploadImage{u: u, imageUrl: imageUrl}
}
