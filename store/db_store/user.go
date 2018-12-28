package db_store

import (
	"github.com/jinzhu/gorm"
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
	return
}

func (u *dbUser) UserUpdate(userId int64, data map[string]interface{}) error {
	if userId <= 0 {
		return model.ErrUserNotExist
	}
	return u.db.Model(model.User{Id: userId}).Select("name", "student_num", "password", "pw_plain", "class_name", "is_admin").Updates(data).Error
}

func (u *dbUser) UserCreate(user *model.User) error {
	return u.db.Create(user).Error
}

func (u *dbUser) UserListByUserIds(userIds []interface{}) (users []*model.User, err error) {
	if len(userIds) == 0 {
		return
	}
	users = make([]*model.User, 0, len(userIds))
	err = u.db.Where("id in (?)", userIds).Find(&users).Error
	return
}

func NewDBUser(db *gorm.DB) model.UserStore {
	return &dbUser{db: db}
}
