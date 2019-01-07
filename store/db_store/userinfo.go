package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/model"
)

type dbUserInfo struct {
	db *gorm.DB
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

func (u *dbUserInfo) CreateUserInfo(userInfo *model.UserInfo) (err error) {
	return u.db.Create(&userInfo).Error
}

func NewDBUserInfo(db *gorm.DB) model.UserInfoStore {
	return &dbUserInfo{db}
}
