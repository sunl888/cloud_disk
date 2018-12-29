package model

import "time"

type Folder struct {
	Id         int64      `gorm:"primary_key" json:"id"`                // ID
	UserId     int64      `gorm:"index" json:"user_id"`                 // 创建者
	ParentId   int64      `gorm:"default:null" json:"parent_id"`        // 父目录
	FolderName string     `gorm:"type:varchar(255)" json:"folder_name"` // 目录名称
	CreatedAt  time.Time  `json:"created_at"`                           // 创建时间
	UpdatedAt  time.Time  `json:"updated_at"`                           // 更新时间
	DeletedAt  *time.Time `gorm:"-" json:"deleted_at"`                  // 软删除时间
}

type FolderStore interface {
}

type FolderService interface {
	FolderStore
}

