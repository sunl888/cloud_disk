package image_uploader

import (
	"github.com/jinzhu/gorm"
)

type dbStore struct {
	db *gorm.DB
}

func (is *dbStore) ImageLoad(hash string) (image *Image, err error) {
	image = &Image{}
	err = is.db.Where(Image{Hash: hash}).First(image).Error
	return
}

func (is *dbStore) ImageCreate(image *Image) error {
	return is.db.Create(image).Error
}

func (is *dbStore) ImageExist(hash string) (bool, error) {
	var count uint
	err := is.db.Model(&Image{}).Where(Image{Hash: hash}).Count(&count).Error
	return count > 0, err
}

//
//// ImageSave 方法会判断图片在数据库中是否已经存在如果存在直接放回 否则创建它
//func (is *dbStore) ImageSave(image *Image) (storedImage *Image, exist bool, err error) {
//	originDB := is.db
//	is.db = is.db.Begin()
//	defer func() {
//		if r := recover(); r != nil {
//			is.db.Rollback()
//		}
//		is.db = originDB
//	}()
//	exist, err = is.ImageExist(image.Hash)
//
//	// tx.Commit().Error
//	if err != nil {
//		return nil, false, err
//	} else {
//		if exist {
//			// 图片已经存在
//			storedImage, err = is.ImageLoad(image.Hash)
//			return storedImage, true, err
//		} else {
//			// 图片不存在
//			is.ImageCreate(image)
//		}
//	}
//}

func NewDBStore(db *gorm.DB) Store {
	return &dbStore{db}
}
