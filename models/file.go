package models

import "time"

//File 文件结构体
type File struct {
	ID        uint       `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
	FID       string     `gorm:"column:fid;unique_index" json:"file_id,omitempty"`
	Owner     string     `gorm:"column:owner" json:"owner"`
	IsDir     bool       `gorm:"column:isdir" json:"isdir"`
	Directory string     `gorm:"column:directory;index" json:"directory"`
	Filename  string     `gorm:"column:filename;index" json:"filename"`
	Size      int64      `gorm:"column:size" json:"size"`
	PFID      string     `gorm:"column:pfid" json:"parent_id"`
}
