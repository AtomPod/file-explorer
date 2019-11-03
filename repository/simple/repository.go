package simple

import (
	"github.com/phantom-atom/file-explorer/repository"
)

type dataRepository struct {
	dbRepository     *dbRepository
	verificationCode *VerificationCodeRepository
}

func (d *dataRepository) File() (repository.FileRepository, error) {
	return d.dbRepository, nil
}
func (d *dataRepository) User() (repository.UserRepository, error) {
	return d.dbRepository, nil
}

func (d *dataRepository) VerificationCode() (repository.VerificationCodeRepository, error) {
	return d.verificationCode, nil
}
