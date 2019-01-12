package minio

import (
	"github.com/minio/minio-go"
	"mime"
	"io"
	"path/filepath"
	. "github.com/zm-dev/go-file-uploader"
	"time"
	"net/url"
	"fmt"
)

type minioUploader struct {
	h           Hasher
	minioClient *minio.Client
	bucketName  string
	h2sn        Hash2StorageName
	s           Store
}

func (mu *minioUploader) saveToMinio(hashValue string, fh FileHeader) error {
	name, err := mu.h2sn.Convent(hashValue)
	if err != nil {
		return fmt.Errorf("hash to storage name error. err:%+v", err)
	}
	_, err = fh.File.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	//obj, _ := mu.minioClient.GetObject(mu.bucketName, name, minio.GetObjectOptions{})
	//_, err = obj.Stat()
	//if err != nil {
	//	if minio.ToErrorResponse(err).Code != "NoSuchKey" {
	//		return nil, fmt.Errorf("minio client get object error. err:%+v", err)
	//	}
	//} else {
	//	// 文件已经存在
	//	return nil
	//}

	ext := filepath.Ext(fh.Filename)
	// 在 apline 镜像中 mime.TypeByExtension 只能用 jpg
	if ext == "jpeg" {
		ext = "jpg"
	}

	_, err = mu.minioClient.PutObject(
		mu.bucketName,
		name,
		fh.File,
		fh.Size,
		minio.PutObjectOptions{ContentType: mime.TypeByExtension(ext)},
	)

	if err != nil {
		return fmt.Errorf("minio client put object error. err:%+v", err)
	}

	return nil
}

func (mu *minioUploader) Upload(fh FileHeader, extra string) (f *FileModel, err error) {
	hashValue, err := mu.h.Hash(fh.File)
	if err != nil {
		return nil, err
	}

	if exist, err := mu.s.FileExist(hashValue); exist && err == nil {
		// 文件已经存在
		return mu.s.FileLoad(hashValue)
	} else if err != nil {
		return nil, err
	}

	err = mu.saveToMinio(hashValue, fh)
	if err != nil {
		return nil, err
	}

	return SaveToStore(mu.s, hashValue, fh, extra)
}

func (mu *minioUploader) PresignedGetObject(hashValue string, expires time.Duration, reqParams url.Values) (u *url.URL, err error) {
	name, err := mu.h2sn.Convent(hashValue)
	if err != nil {
		return nil, err
	}
	return mu.minioClient.PresignedGetObject(mu.bucketName, name, expires, reqParams)
}

func (mu *minioUploader) Store() Store {
	return mu.s
}

type readFile struct {
	*minio.Object
}

func (rf *readFile) Stat() (*FileInfo, error) {
	info, err := rf.Object.Stat()
	if err != nil {
		return nil, err
	}
	return &FileInfo{
		LastModified: info.LastModified,
		Size:         info.Size,
		ContentType:  info.ContentType,
	}, nil
}

func (mu *minioUploader) ReadFile(hashValue string) (rf ReadFile, err error) {
	name, err := mu.h2sn.Convent(hashValue)
	if err != nil {
		return
	}
	obj, err := mu.minioClient.GetObject(mu.bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return
	}
	return &readFile{obj}, nil
}

func NewMinioUploader(h Hasher, minioClient *minio.Client, s Store, bucketName string, h2sn Hash2StorageName) Uploader {
	if h2sn == nil {
		h2sn = Hash2StorageNameFunc(DefaultHash2StorageNameFunc)
	}
	return &minioUploader{h: h, minioClient: minioClient, bucketName: bucketName, h2sn: h2sn, s: s}
}
