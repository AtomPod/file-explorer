package simple

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/phantom-atom/file-explorer/cache"
	"github.com/phantom-atom/file-explorer/config"
	"github.com/phantom-atom/file-explorer/repository"
)

type context struct {
	*dataRepository
}

//NewContext 创建context
func NewContext(
	db *gorm.DB,
	cache cache.Cache,
	configFun func() *config.Config,
	now func() time.Time) (repository.DataContext, error) {

	verificationCodeRepo := newVerificationCodeRepository(
		cache, configFun, now,
	)

	dbRepo := newDBRepository(db, configFun)

	return &context{
		&dataRepository{
			dbRepository:     dbRepo,
			verificationCode: verificationCodeRepo,
		},
	}, nil
}

func (c *context) Unit() (repository.UnitOfWork, error) {
	return newUnixOfWork(c)
}
