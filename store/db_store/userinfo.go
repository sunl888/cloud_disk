package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/model"
)

type dbUserInfo struct {
	db *gorm.DB
}

func (u *dbUserInfo) UpdateUsedStorage(userId int64, usedStorage uint64) (err error) {
	if userId <= 0 {
		return model.ErrUserNotExist
	}
	err = u.db.Model(model.UserInfo{}).
		Where("user_id = ?", userId).
		Update("used_storage", usedStorage).
		Error
	return
}

func (u *dbUserInfo) LoadUserInfo(userId int64) (userInfo *model.UserInfo, err error) {
	if userId <= 0 {
		return nil, model.ErrUserNotExist
	}
	userInfo = &model.UserInfo{}
	err = u.db.First(&userInfo, "user_id = ?", userId).Error
	if gorm.IsRecordNotFoundError(err) {
		err = model.ErrUserNotExist
	}
	return
}

func NewDBUserInfo(db *gorm.DB) model.UserInfoStore {
	return &dbUserInfo{db}
}
