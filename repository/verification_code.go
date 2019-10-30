package repository

import (
	"github.com/phantom-atom/file-explorer/models"
)

//VerificationCodeRepository 验证码仓库接口
type VerificationCodeRepository interface {
	CreateVerificationCode(*models.VerificationCode) error
	DeleteVerificationCode(typ string, code string, target string) error
	GetVerificationCodeByCodeTypeTarget(typ string, code string, target string) (*models.VerificationCode, error)
}
