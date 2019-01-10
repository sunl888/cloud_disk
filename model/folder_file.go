package model

type FolderFile struct {
	FolderId int64  `gorm:"type:BIGINT;index:folder_id_file_id_unique_index"`
	FileId   int64  `gorm:"type:BIGINT;index:folder_id_file_id_unique_index"`
	Filename string `gorm:"type:varchar(255);NOT NULL" json:"filename"`
}

func (*FolderFile) TableName() string {
	return "folder_files"
}

type FolderFileStore interface {
	// 加载指定目录的文件s
	LoadFolderFilesByFolderIds(folderIds []int64, userId int64) (folderFiles []*FolderFile, err error)
	// 加载指定目录的指定文件s的详细信息
	LoadFolderFilesByFolderIdAndFileIds(folderId int64, fileIds []int64, userId int64) (folderFiles []*FolderFile, err error)
}

type FolderFileService interface {
	FolderFileStore
}
