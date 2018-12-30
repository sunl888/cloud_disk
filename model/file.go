package model

import (
	"time"
)

type File struct {
	Id        int64      `gorm:"primary_key" json:"id"`                 // ID
	Folders   []*Folder  `gorm:"many2many:folders_file;" json:"-"`      // 仅仅用来作为关联,表示该文件可以被多个目录引用,而不需要显示 所以用 - 忽略
	Filename  string     `gorm:"not null" json:"filename"`              // 文件名称
	Hash      string     `gorm:"type:varchar(32);not null" json:"hash"` // 文件Hash
	Format    string     `gorm:"not null" json:"format"`                // 文件MimeType 例如: video/mp4 -> .mp4
	Extra     string     `json:"extra" gorm:"not null;type:text"`       // extra
	Size      int64      `gorm:"type:bigint" json:"size"`               // 文件大小
	CreatedAt time.Time  `json:"created_at"`                            // 创建时间
	UpdatedAt time.Time  `json:"updated_at"`                            // 更新时间
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`                // 软删除时间
}

type FileStore interface {
	//DeleteFile(id int64) error
	//DeletedFileList(limit, offset int64) (files []*File, count int64, err error)
	//RecoverFile(id int64) error
	//ListFile(limit, offset int64) (files []*File, count int64, err error)
	//IsExistFile(id, userId int64) (isExist bool, err error)
	//UpdateFile(id int64, file *File) (err error)
	//BetchDeleteFile(ids []int64) (err error)
	SaveFileToFolder(file *File, folder *Folder) (err error)
}

type FileService interface {
	FileStore
	MoveFile(fromId, toId int64) // 移动文件或者目录到指定的目录下
}
