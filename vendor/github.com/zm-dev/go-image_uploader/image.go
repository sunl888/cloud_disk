package image_uploader

import "time"

// please create image table
type Image struct {
	Hash      string    `json:"hash" gorm:"primary_key;type:char(32)"`
	// Url       string    `json:"image_url" gorm:"-"`
	Format    string    `json:"format" gorm:"not null"`
	Title     string    `json:"title" gorm:"not null"`
	Width     uint      `json:"width" gorm:"type:MEDIUMINT UNSIGNED;not null"`
	Height    uint      `json:"height" gorm:"type:MEDIUMINT UNSIGNED;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
