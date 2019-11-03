package simple

import (
	"github.com/phantom-atom/file-explorer/repository"
)

type unixOfWork struct {
	*dataRepository
}

func newUnixOfWork(c *context) (repository.UnitOfWork, error) {
	dbRepo, err := c.dbRepository.Begin()
	if err != nil {
		return nil, err
	}
	return &unixOfWork{
		&dataRepository{
			dbRepository:     dbRepo,
			verificationCode: c.verificationCode,
		},
	}, nil
}

func (u *unixOfWork) Commit() error {
	return u.dbRepository.Commit()
}

func (u *unixOfWork) Rollback() error {
	return u.dbRepository.Rollback()
}
