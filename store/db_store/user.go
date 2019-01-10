package db_store

import (
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
)

type dbUser struct {
	db *gorm.DB
}

func (u *dbUser) UserExist(id int64) (bool, error) {
	var count uint8
	err := u.db.Model(model.User{}).Where(model.User{Id: id}).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u *dbUser) UserIsNotExistErr(err error) bool {
	return model.UserIsNotExistErr(err)
}

func (u *dbUser) UserLoad(id int64) (user *model.User, err error) {
	if id <= 0 {
		return nil, model.ErrUserNotExist
	}
	user = &model.User{}
	err = u.db.Where(model.User{Id: id}).First(user).Error
	if gorm.IsRecordNotFoundError(err) {
		err = model.ErrUserNotExist
	}

	group := &model.Group{}
	err = u.db.Where("id = ?", user.GroupId).First(&group).Error
	if gorm.IsRecordNotFoundError(err) {
		err = errors.RecordNotFound("用户组不存在")
	}
	user.Group = group
	return
}

func (u *dbUser) UserUpdate(userId int64, data map[string]interface{}) error {
	if userId <= 0 {
		return model.ErrUserNotExist
	}
	return u.db.Model(model.User{Id: userId}).
		Select(
			"name", "gender", "used_storage", "password", "is_ban", "used_storage", "group_id",
			"is_admin", "nickname", "email", "avatar_hash", "profile",
		).
		Updates(data).Error
}

func (u *dbUser) UserCreate(user *model.User) (err error) {
	err = u.db.Create(&user).Error
	return
}

func NewDBUser(db *gorm.DB) model.UserStore {
	return &dbUser{db: db}
}
