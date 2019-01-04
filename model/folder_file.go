package model

type FolderFile struct {
	FolderId int64  `gorm:"type:BIGINT;index:folder_id_file_id_unique_index"`
	FileId   int64  `gorm:"type:BIGINT;index:folder_id_file_id_unique_index"`
	Filename string `gorm:"type:varchar(255);NOT NULL" json:"filename"`
}

func (*FolderFile) TableName() string {
	return "folder_files"
}
