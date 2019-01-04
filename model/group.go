package model

import "time"

type Group struct {
	Id         int64     `gorm:"type:BIGINT;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`  // ID
	Name       string    `gorm:"type:varchar(32)" json:"name"`                      // 组名
	MaxStorage int64     `gorm:"type:BIGINT;default:1073741824" json:"max_storage"` // 最大容量/KB 默认1TB
	AllowShare bool      `gorm:"type:TINYINT;default:0" json:"allow_share"`         // 是否允许分享文件
	CreatedAt  time.Time `json:"created_at"`                                        // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                                        // 更新时间
}

type GroupStore interface {
}

type GroupService interface {
	GroupStore
}
