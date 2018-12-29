package model

import (
	"time"
)

type File struct {
	Id        int64      `gorm:"primary_key" json:"id"`                 // ID
	Filename  string     `gorm:"not null" json:"filename"`              // 文件名称
	Hash      string     `gorm:"type:varchar(32);not null" json:"hash"` // 文件Hash
	Format    string     `gorm:"not null" json:"format"`                // 文件MimeType 例如: video/mp4 -> .mp4
	Size      int64      `gorm:"type:bigint" json:"size"`               // 文件大小
	Extra     string     `gorm:"not null;type:text" json:"extra"`       // 扩展名
	FolderId  int64      `gorm:"type:bigint;not null" json:"folder_id"` // 所属目录
	CreatedAt time.Time  `json:"created_at"`                            // 创建时间
	UpdatedAt time.Time  `json:"updated_at"`                            // 更新时间
	DeletedAt *time.Time `json:"deleted_at"`                            // 软删除时间
}

type FileStore interface {
	DeleteFile(id int64) error
	DeletedFileList(limit, offset int64) (files []*File, count int64, err error)
	RecoverFile(id int64) error
	ListFile(limit, offset int64) (files []*File, count int64, err error)
	IsExistFile(id, userId int64) (isExist bool, err error)
	UpdateFile(id int64, file *File) (err error)
}

type FileService interface {
	FileStore
	MoveFile(fromId, toId int64) // 移动文件或者目录到指定的目录下
}
