package image_uploader

import (
	"github.com/spf13/afero"
	"io"
	"path/filepath"
	"os"
)

type aferoUploader struct {
	h    Hasher
	s    Store
	fs   afero.Fs
	h2sn Hash2StorageName
}

func (au *aferoUploader) saveToFs(hashValue string, f File) error {
	name, err := au.h2sn.Convent(hashValue)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	baseDir:=filepath.Dir(name)
	_, err = au.fs.Stat(baseDir)
	if os.IsNotExist(err) {
		au.fs.MkdirAll(baseDir, 0755)
	}
	// todo savepath
	newFile, err := au.fs.Create(name)
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, f)
	return err
}

func (au *aferoUploader) Upload(fh FileHeader) (*Image, error) {
	info, err := DecodeImageInfo(fh.File)
	if err != nil {
		return nil, err
	}

	hashValue, err := au.h.Hash(fh.File)
	if err != nil {
		return nil, err
	}

	if exist, err := au.s.ImageExist(hashValue); exist && err == nil {
		// 图片已经存在
		return au.s.ImageLoad(hashValue)
	} else if err != nil {
		return nil, err
	}

	if err := au.saveToFs(hashValue, fh.File); err != nil {
		return nil, err
	}
	return saveToStore(au.s, hashValue, fh.Filename, info)
}

func (au *aferoUploader) UploadFromURL(u string, filename string) (*Image, error) {
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
	return au.Upload(fh)
}

func NewAferoUploader(h Hasher, s Store, fs afero.Fs, h2sn Hash2StorageName) Uploader {
	if h2sn == nil {
		h2sn = Hash2StorageNameFunc(DefaultHash2StorageNameFunc)
	}
	return &aferoUploader{h: h, s: s, fs: fs, h2sn: h2sn}
}
