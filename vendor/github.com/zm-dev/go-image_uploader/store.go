package image_uploader

type Store interface {
	ImageExist(hash string) (bool, error)
	ImageLoad(hash string) (*Image, error)
	ImageCreate(image *Image) error
}

func saveToStore(s Store, hashValue, title string, info ImageInfo) (imageModel *Image, err error) {
	imageModel = &Image{Hash: hashValue, Format: info.format, Title: title, Width: info.width, Height: info.height}
	err = s.ImageCreate(imageModel)
	return
}
