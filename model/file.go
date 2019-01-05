package model

import (
	"time"
)

type File struct {
	Id        int64     `gorm:"type:BIGINT AUTO_INCREMENT;PRIMARY_KEY;NOT NUll" json:"id"` // ID
	Filename  string    `gorm:"type:char(32); NOT NULL" json:"filename"`                   // 文件名称
	Hash      string    `gorm:"type:varchar(32);NOT NULL" json:"hash"`                     // 文件Hash
	Format    string    `gorm:"NOT NULL" json:"format"`                                    // 文件MimeType 例如: video/mp4 -> .mp4
	Extra     string    `gorm:"NOT NULL;type:TEXT" json:"extra"`                           // extra
	Size      int64     `gorm:"type:BIGINT" json:"size"`                                   // 文件大小
	CreatedAt time.Time `json:"created_at"`                                                // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                                                // 更新时间
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
	// 重命名文件
	RenameFile(folderId, fileId int64, newName string) (err error)
	// 加载文件
	LoadFile(folderId, fileId, userId int64) (file *File, err error)
}

type FileService interface {
	FileStore
}
