package model

import (
	"errors"
	"time"
)

type User struct {
	Id          int64  `gorm:"type:BIGINT AUTO_INCREMENT;PRIMARY_KEY;NOT NUll" json:"id"` // id
	Name        string `gorm:"type:varchar(50)" json:"name"`                              // 账号
	Email       string `gorm:"type:varchar(255)" json:"email"`                            // 用户邮箱
	IsBan       bool   `gorm:"type:TINYINT;default:0" json:"is_ban"`                      // 是否禁用
	Group       *Group `gorm:"PRELOAD:false" json:"group,omitempty"`                      // 用户组
	Gender      bool   `gorm:"type:TINYINT;default:0" json:"gender"`                      // 性别
	Profile     string `gorm:"type:varchar(255)" json:"profile"`                          // 简介
	GroupId     int64  `gorm:"type:BIGINT;NOT NULL" json:"group_id"`                      // 所属用户组
	IsAdmin     bool   `gorm:"type:TINYINT" json:"is_admin"`                              // 是否为超级管理员
	PwPlain     string `gorm:"type:varchar(20);not null" json:"pw_plain"`                 // password 明文存储防止到时候有些人忘了
	Password    string `gorm:"type:varchar(64);not null" json:"password"`                 // hash(密码)
	Nickname    string `gorm:"type:varchar(255)" json:"nickname"`                         // 昵称
	AvatarHash  string `gorm:"type:varchar(32)" json:"avatar_hash"`                       // 头像
	UsedStorage uint64 `gorm:"type:BIGINT;default:0" json:"used_storage"`                 // 已使用的空间大小/KB
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserStore interface {
	UserExist(userId int64) (bool, error)
	UserLoad(userId int64) (*User, error)
	UserIsNotExistErr(err error) bool
	UserUpdate(userId int64, data map[string]interface{}) error
	UserCreate(user *User) error
	UserList(offset, limit int64) (user []*User, count int64, err error)
	UserListByUserIds(userIds []interface{}) ([]*User, error)
}

type UserService interface {
	UserStore
	UserLogin(account, password string) (*Ticket, error)
	UserRegister(account string, certificateType CertificateType, password string) (userId int64, err error)
	UserUpdatePassword(userId int64, newPassword string) (err error)
	UserUpdateUsedStorage(userId int64, usedStorage uint64) (err error)
	UserUpdateBanStatus(userId int64, newBanStatus bool) (err error)
}

var ErrUserNotExist = errors.New("user not exist")

func UserIsNotExistErr(err error) bool {
	return err == ErrUserNotExist
}
