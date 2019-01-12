package go_file_uploader

import (
	"mime"
	"path/filepath"
	"errors"
)

type Store interface {
	FileExist(hash string) (bool, error)
	FileLoad(hash string) (*FileModel, error)
	FileIsNotExistError(error) bool
	FileCreate(file *FileModel) error
}

func SaveToStore(s Store, hashValue string, fh FileHeader, extra string) (fileModel *FileModel, err error) {
	ext := filepath.Ext(fh.Filename)
	// 在 apline 镜像中 mime.TypeByExtension 只能用 jpg
	if ext == "jpeg" {
		ext = "jpg"
	}
	fileModel = &FileModel{Hash: hashValue, Format: mime.TypeByExtension(ext), Filename: fh.Filename, Size: fh.Size, Extra: extra}
	err = s.FileCreate(fileModel)
	return
}

var ErrFileNotExist = errors.New("file not exist")

func FileIsNotExistError(err error) bool {
	return err == ErrFileNotExist
}
