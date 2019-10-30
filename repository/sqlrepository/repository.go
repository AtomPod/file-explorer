package sqlrepository

import (
	"github.com/jinzhu/gorm"
	"github.com/phantom-atom/file-explorer/config"
)

//Repository 数据仓库
type Repository struct {
	db     *gorm.DB
	config func() *config.Config
}

//NewRepository 创建数据仓库
func NewRepository(db *gorm.DB, configFunc func() *config.Config) *Repository {
	return &Repository{
		db:     db,
		config: configFunc,
	}
}
