package model

// 次序千万不能更改,否则 gorm 的 select 就不能用了
type FolderFile struct {
	FolderId     int64  `gorm:"type:BIGINT;NOT NUll" json:"folder_id"`
	OriginFileId int64  `gorm:"type:BIGINT;NOT NUll" json:"origin_file_id"`
	Filename     string `gorm:"type:varchar(255);NOT NULL" json:"filename"`
	FileId       int64  `gorm:"type:BIGINT AUTO_INCREMENT;PRIMARY_KEY;NOT NUll" json:"file_id"`
}

type WrapFolderFile struct {
	FileId       int64  `json:"file_id"`
	FolderId     int64  `json:"folder_id"`
	FileSize     int64  `json:"file_size"`
	Filename     string `json:"filename"`
	Format       string `json:"format"`
	RelativePath string `json:"relative_path"`
}

func (*FolderFile) TableName() string {
	return "folder_files"
}

type FolderFileStore interface {
	// 加载指定目录的文件s
	LoadFolderFilesByFolderIds(folderIds []int64, userId int64) (folderFiles []*WrapFolderFile, err error)
	// 加载指定目录的指定文件s的详细信息
	LoadFolderFilesByFolderIdAndFileIds(folderId int64, fileIds []int64, userId int64) (folderFiles []*WrapFolderFile, err error)
	// 是否存在
	ExistFile(filename string, folderId, userId int64) (isExist bool, err error)
}

type FolderFileService interface {
	FolderFileStore
}
