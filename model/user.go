package model

import (
	"errors"
	"time"
)

type User struct {
	Id        int64     `gorm:"type:BIGINT AUTO_INCREMENT;PRIMARY_KEY;NOT NUll" json:"id"`
	Name      string    `gorm:"type:varchar(50)"`
	Password  string    `gorm:"type:varchar(64);not null"`
	PwPlain   string    `gorm:"type:varchar(20);not null"` // password 明文存储防止到时候有些人忘了
	IsAdmin   bool      `gorm:"type:TINYINT"`
	UserInfo  *UserInfo `gorm:"PRELOAD:false" json:"user_info,omitempty"`
	Group     *Group    `gorm:"PRELOAD:false" json:"group,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserStore interface {
	UserExist(userId int64) (bool, error)
	UserLoad(userId int64) (*User, error)
	UserIsNotExistErr(err error) bool
	UserUpdate(userId int64, data map[string]interface{}) error
	UserCreate(user *User) error
	UserListByUserIds(userIds []interface{}) ([]*User, error)
	UserLoadAndRelated(userId int64) (user *User, err error)
}

type UserService interface {
	UserStore
	UserLogin(account, password string) (*Ticket, error)
	UserRegister(account string, certificateType CertificateType, password string) (userId int64, err error)
	UserUpdatePassword(userId int64, newPassword string) error
}

var ErrUserNotExist = errors.New("user not exist")

func UserIsNotExistErr(err error) bool {
	return err == ErrUserNotExist
}
