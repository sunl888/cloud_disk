package image_uploader

import (
	"github.com/minio/minio-go"
	"mime"
	"path/filepath"
	"io"
)

type minioUploader struct {
	h           Hasher
	s           Store
	minioClient *minio.Client
	bucketName  string
	h2sn        Hash2StorageName
}

func (mu *minioUploader) saveToMinio(hashValue string, fh FileHeader, info ImageInfo) error {
	name, err := mu.h2sn.Convent(hashValue)
	if err != nil {
		return err
	}
	_, err = fh.File.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	// 在 apline 镜像中 mime.TypeByExtension 只能用 jpg
	if info.format == "jpeg" {
		info.format = "jpg"
	}

	_, err = mu.minioClient.PutObject(
		mu.bucketName,
		name,
		fh.File,
		fh.Size,
		minio.PutObjectOptions{ContentType: mime.TypeByExtension("." + info.format)},
	)
	return err
}

func (mu *minioUploader) Upload(fh FileHeader) (*Image, error) {
	info, err := DecodeImageInfo(fh.File)
	if err != nil {
		return nil, err
	}
	hashValue, err := mu.h.Hash(fh.File)
	if err != nil {
		return nil, err
	}
	if exist, err := mu.s.ImageExist(hashValue); exist && err == nil {
		// 图片已经存在
		return mu.s.ImageLoad(hashValue)
	} else if err != nil {
		return nil, err
	}

	if err := mu.saveToMinio(hashValue, fh, info); err != nil {
		return nil, err
	}

	return saveToStore(mu.s, hashValue, fh.Filename, info)
}

func (mu *minioUploader) UploadFromURL(u string, filename string) (*Image, error) {
	if filename == "" {
		filename = filepath.Base(u)
	}
	file, size, err := DownloadImage(u)

	if err != nil {
		return nil, err
	}

	defer removeFile(file)

	fh := FileHeader{
		Filename: filename,
		Size:     size,
		File:     file,
	}

	return mu.Upload(fh)
}

func NewMinioUploader(h Hasher, s Store, minioClient *minio.Client, bucketName string, h2sn Hash2StorageName) Uploader {
	if h2sn == nil {
		h2sn = Hash2StorageNameFunc(DefaultHash2StorageNameFunc)
	}
	return &minioUploader{h: h, s: s, minioClient: minioClient, bucketName: bucketName, h2sn: h2sn}
}
