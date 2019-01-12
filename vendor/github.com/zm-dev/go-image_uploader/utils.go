package image_uploader

import (
	"io"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	// todo 下面这两个类型可以考虑不要
	//_ "golang.org/x/image/bmp"
	//_ "golang.org/x/image/tiff"
)

type ImageInfo struct {
	width, height uint
	format        string
}

func DecodeImageInfo(file File) (info ImageInfo, err error) {
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return
	}
	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return ImageInfo{}, err
	}

	return ImageInfo{
		width:  uint(config.Width),
		height: uint(config.Height),
		format: format,
	}, nil

}

func IsUnknownFormat(err error) bool {
	return image.ErrFormat == err
}
