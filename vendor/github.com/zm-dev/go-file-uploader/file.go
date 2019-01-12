package go_file_uploader

import "time"

type FileModel struct {
	Hash      string    `json:"hash" gorm:"primary_key;type:char(32)"`
	Format    string    `json:"format" gorm:"not null"`
	Filename  string    `json:"filename" gorm:"not null"`
	Size      int64     `json:"size" gorm:"not null"`
	Extra     string    `json:"extra" gorm:"not null;type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (FileModel) TableName() string {
	return "files"
}
