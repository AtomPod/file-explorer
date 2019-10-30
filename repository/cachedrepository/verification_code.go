package cachedrepository

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/phantom-atom/file-explorer/cache"

	"github.com/phantom-atom/file-explorer/config"
	"github.com/phantom-atom/file-explorer/models"
)

//VerificationCodeRepository 验证码仓库实现
type VerificationCodeRepository struct {
	cache  cache.Cache
	config func() *config.Config
	now    func() time.Time
}

//NewVerificationCodeRepository 创建VerificationCodeRepository
func NewVerificationCodeRepository(
	cache cache.Cache,
	configFun func() *config.Config,
	now func() time.Time,
) *VerificationCodeRepository {
	return &VerificationCodeRepository{
		cache:  cache,
		config: configFun,
		now:    now,
	}
}

func (r *VerificationCodeRepository) getCodeKey(typ string, code string, target string) string {
	return fmt.Sprintf("%s.%s.%s", typ, target, code)
}

//CreateVerificationCode *
func (r *VerificationCodeRepository) CreateVerificationCode(code *models.VerificationCode) error {
	if code == nil {
		return errors.New("invalid argument 'code'")
	}

	key := r.getCodeKey(code.Type, code.Code, code.Target)
	if code.CreatedAt.IsZero() {
		code.CreatedAt = r.now()
	}

	data, err := json.Marshal(code)
	if err != nil {
		return err
	}

	if err := r.cache.Set(&cache.Entity{
		Key:        key,
		Value:      data,
		Expiration: code.Expiration,
	}); err != nil {
		return err
	}
	return nil
}

//DeleteVerificationCode *
func (r *VerificationCodeRepository) DeleteVerificationCode(typ string, code string, target string) error {
	if typ == "" && code == "" && target == "" {
		return errors.New("invalid argument 'typ&code&target'")
	}

	if typ == "" {
		typ = "*"
	}

	if code == "" {
		code = "*"
	}

	if target == "" {
		target = "*"
	}

	key := r.getCodeKey(typ, code, target)
	if err := r.cache.Del(key); err != nil {
		return err
	}
	return nil
}

//GetVerificationCodeByCodeTypeTarget *
func (r *VerificationCodeRepository) GetVerificationCodeByCodeTypeTarget(typ string, code string, target string) (*models.VerificationCode, error) {
	if typ == "" {
		return nil, errors.New("invalid argument 'typ'")
	}

	if code == "" {
		return nil, errors.New("invalid argument 'code'")
	}

	if target == "" {
		return nil, errors.New("invalid argument 'target'")
	}

	key := r.getCodeKey(typ, code, target)
	entitys, err := r.cache.Get(key)

	if err != nil {
		if err == cache.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	entity := entitys[0]
	codeModel := &models.VerificationCode{}
	if err := json.Unmarshal(entity.Value, codeModel); err != nil {
		return nil, err
	}

	return codeModel, nil
}
