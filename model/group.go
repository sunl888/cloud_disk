package model

import (
	"errors"
	"time"
)

type Group struct {
	Id         int64     `gorm:"type:BIGINT AUTO_INCREMENT;PRIMARY_KEY;NOT NUll" json:"id"` // ID
	Name       string    `gorm:"type:varchar(32)" json:"name"`                              // 组名
	MaxStorage uint64    `gorm:"type:BIGINT" json:"max_storage"`                            // 最大容量/KB 默认1TB
	AllowShare bool      `gorm:"type:TINYINT;default:0" json:"allow_share"`                 // 是否允许分享文件
	Users      []*User   `json:"users,omitempty"`                                           // 用户列表
	CreatedAt  time.Time `json:"created_at"`                                                // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                                                // 更新时间
}

type WrapGroupList struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	MaxStorage uint64    `json:"max_storage"`
	AllowShare bool      `json:"allow_share"`
	UserCount  int64     `json:"user_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

var ErrGroupNotExist = errors.New("group not exist")
var ErrGroupAlreadyExist = errors.New("group already exist")

const (
	DefaultAllowSize = 5 << 30 // 5GB
	VipAllowSize     = 5 << 40 // 5TB
)

type GroupStore interface {
	GroupCreate(group *Group) (err error)
	GroupDelete(id int64) (err error)
	GroupExist(name string) (isExist bool, err error)
	GroupUpdate(id int64, data map[string]interface{}) (err error)
	GroupList(offset, limit int64) (groups []*WrapGroupList, count int64, err error)
}

type GroupService interface {
	GroupStore
}
