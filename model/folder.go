package model

import (
	"errors"
	"time"
)

type Folder struct {
	Id         int64     `gorm:"type:BIGINT AUTO_INCREMENT;PRIMARY_KEY;NOT NUll" json:"id"` // ID
	FolderName string    `gorm:"type:varchar(255)" json:"folder_name"`                      // 目录名称
	ParentId   int64     `gorm:"type:BIGINT;default:0" json:"parent_id"`                    // 父目录
	UserId     int64     `gorm:"type:BIGINT;index:user_id" json:"user_id"`                  // 创建者
	Key        string    `gorm:"type:varchar(255);default:''" json:"key"`                   // 辅助键
	Level      int64     `gorm:"type:INT;default:1" json:"level"`                           // 辅助键
	Files      []*File   `json:"files"`                                                     // many2many
	Folders    []*Folder `gorm:"foreignkey:ParentId" json:"folders"`                        // one2many 当前目录下的目录
	CreatedAt  time.Time `json:"created_at"`                                                // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                                                // 更新时间
}

const (
	FolderKeyPrefix = "-"
)

var FolderAlreadyExisted = errors.New("目录已存在")

type FolderStore interface {
	// 创建一个目录
	CreateFolder(folder *Folder) (err error)
	// 目录是否存在
	ExistFolder(userId, parentId int64, folderName string) (isExist bool)
	// 当 id != 0 则表示加载指定目录, 当 id == 0 则表示加载根目录
	LoadFolder(id, userId int64, isLoadRelated bool) (folder *Folder, err error)
	// 删除指定目录
	DeleteFolder(ids []int64, userId int64) (err error)
	// 移动目录
	MoveFolder(to *Folder, ids []int64) (err error)
	// 复制目录
	CopyFolder(to *Folder, foders []*Folder) (totalSize uint64, err error)
	// 重命名目录
	RenameFolder(id int64, newName string) (err error)
	// 目录列表
	ListFolder(folderIds []int64, userId int64) (folder []*Folder, err error)
}

type FolderService interface {
	FolderStore
}
