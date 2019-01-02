package model

import (
	"time"
)

type File struct {
	Id        int64     `gorm:"primary_key" json:"id"`                 // ID
	Folders   []*Folder `gorm:"many2many:folder_files;" json:"-"`      // 仅仅用来作为关联,而不需要显示 所以用 - 忽略
	Filename  string    `gorm:"not null" json:"filename"`              // 文件名称
	Hash      string    `gorm:"type:varchar(32);not null" json:"hash"` // 文件Hash
	Format    string    `gorm:"not null" json:"format"`                // 文件MimeType 例如: video/mp4 -> .mp4
	Extra     string    `json:"extra" gorm:"not null;type:text"`       // extra
	Size      int64     `gorm:"type:bigint" json:"size"`               // 文件大小
	CreatedAt time.Time `json:"created_at"`                            // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                            // 更新时间
}

type FileStore interface {
	// 保存文件到指定目录
	SaveFileToFolder(file *File, folder *Folder) (err error)
	// 删除文件
	DeleteFile(ids []int64, folderId int64) (err error)
	// 移动文件
	MoveFile(fromId, toId int64, fileIds []int64) (err error)
	// 复制文件
	CopyFile(toId int64, fileIds []int64) (err error)
}

type FileService interface {
	FileStore
}
