package simple

import (
	"github.com/jinzhu/gorm"
	"github.com/phantom-atom/file-explorer/config"
)

//dbRepository 数据仓库
type dbRepository struct {
	db     *gorm.DB
	config func() *config.Config
}

//newDBRepository 创建数据仓库
func newDBRepository(db *gorm.DB, configFunc func() *config.Config) *dbRepository {
	return &dbRepository{
		db:     db,
		config: configFunc,
	}
}

//Begin 开始事务
func (r *dbRepository) Begin() (*dbRepository, error) {
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}
	return newDBRepository(tx, r.config), nil
}

//Commit 提交事务
func (r *dbRepository) Commit() error {
	return r.db.Commit().Error
}

//Rollback 回滚事务
func (r *dbRepository) Rollback() error {
	return r.db.Rollback().Error
}
