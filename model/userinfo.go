package model

import "time"

type UserInfo struct {
	Id          int64     `gorm:"primary_key" json:"id"`                     // ID
	UserId      int64     `gorm:"unique_index" json:"user_id"`               // UID
	Nickname    string    `json:"nickname"`                                  // 昵称
	AvatarHash  string    `gorm:"type:varchar(32)" json:"avatar_hash"`       // 头像
	Profile     string    `json:"profile"`                                   // 简介
	Email       string    `json:"email"`                                     // 用户邮箱
	IsBan       bool      `gorm:"type:tinyint;default:0" json:"is_ban"`      // 是否禁用
	UsedStorage int64     `gorm:"type:bigint;default:0" json:"used_storage"` // 已使用的空间大小/KB
	GroupId     int64     `gorm:"type:bigint;not null" json:"group_id"`      // 所属用户组
	CreatedAt   time.Time `json:"created_at"`                                // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`                                // 更新时间
}

type UserInfoStore interface {
}

type UserInfoService interface {
	UserInfoStore
}
