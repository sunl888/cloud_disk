package image_uploader

import (
	"io"
	"context"
	"errors"
	"io/ioutil"
	"os"
	"net/http"
	"fmt"
)

type File interface {
	io.Reader
	io.Seeker
}

type FileHeader struct {
	Filename string
	Size     int64
	File     File
}

type Uploader interface {
	Upload(fh FileHeader) (*Image, error)
	UploadFromURL(u string, filename string) (*Image, error)
}

func removeFile(f *os.File) {
	f.Close()
	os.Remove(f.Name())
}

func DownloadImage(u string) (f *os.File, size int64, err error) {
	f, err = ioutil.TempFile("", "image_uploader")
	if err != nil {
		return nil, 0, fmt.Errorf("create temp file faild. err=%+v", err)
	}
	resp, err := http.Get(u)

	if err != nil {
		removeFile(f)
		return nil, 0, fmt.Errorf("http.Get image faild. err=%+v", err)
	}

	size, err = io.Copy(f, resp.Body)

	if err != nil {
		removeFile(f)
		return nil, 0, fmt.Errorf(" io.Copy faild. err=%+v", err)
	}

	return
}

func Upload(ctx context.Context, fh FileHeader) (*Image, error) {
	u, ok := FromContext(ctx)
	if !ok {
		return nil, errors.New("uploader不存在")
	}
	return u.Upload(fh)
}


