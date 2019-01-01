package model

import "time"

type Folder struct {
	Id         int64      `gorm:"primary_key" json:"id"`                                     // ID
	Files      []*File    `gorm:"many2many:folders_file;" json:"files"`                      // many2many
	Folders    []*Folder  `gorm:"foreignkey:ParentId" json:"folders"`                        // one2many 当前目录下的目录
	UserId     int64      `gorm:"index:user_id_folder_name_unique_index" json:"user_id"`     // 创建者(组合唯一)
	FolderName string     `gorm:"index:user_id_folder_name_unique_index" json:"folder_name"` // 目录名称(组合唯一)
	ParentId   int64      `gorm:"default:0" json:"parent_id"`                                // 父目录
	Key        string     `gorm:"default:''" json:"key"`                                     // 辅助键
	Level      int64      `gorm:"default:1" json:"level"`                                    // 辅助键
	CreatedAt  time.Time  `json:"created_at"`                                                // 创建时间
	UpdatedAt  time.Time  `json:"updated_at"`                                                // 更新时间
	DeletedAt  *time.Time `json:"deleted_at"`                                                // 软删除时间
}

const (
	FolderKeyPrefix = "-"
)

type FolderStore interface {
	// 创建一个目录
	CreateFolder(folder *Folder) (err error)
	ExistFolder(userId int64, folderName string) (isExist bool)
	// 加载指定目录 可以选择同时加载其[子目录和文件列表]
	// 当 id != 0 则表示加载指定目录, 当 id == 0 则表示加载根目录
	// 第三个参数表示是否加载关联模型(Files, Folders)
	LoadFolder(id, userId int64, isLoadRelated bool) (folder *Folder, err error)
}

type FolderService interface {
	FolderStore
}
