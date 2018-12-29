package model

import "time"

type Share struct {
	Id          int64      `gorm:"primary_key" json:"id"`                                        // ID
	Type        string     `gorm:"type:enum('private','publish') default:'publish'" json:"type"` // 分享类型 私有分享和公开分享
	SourceType  string     `gorm:"type:enum('file','dir')" json:"source_type"`                   // 资源类型 文件还是目录
	SourceId    string     `json:"source_id"`                                                    // 资源ID 对应files表或者folders表的ID字段
	SharePwd    string     `gorm:"type:varchar(64);not null" json:"share_pwd"`                   // 分享密码
	DownloadNum int64      `gorm:"default:0" json:"download_num"`                                // 下载次数
	ViewNum     int64      `gorm:"default:0" json:"view_num"`                                    // 浏览次数
	EndAt       *time.Time `json:"end_at"`                                                       // 结束时间
	CreatedAt   time.Time  `json:"created_at"`                                                   // 创建时间
	UpdatedAt   time.Time  `json:"updated_at"`                                                   // 更新时间
}

type ShareStore interface {
}

type ShareService interface {
	ShareStore
}
