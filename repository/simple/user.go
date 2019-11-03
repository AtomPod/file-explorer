package simple

import (
	"github.com/jinzhu/gorm"
	"github.com/phantom-atom/file-explorer/models"
)

func (r *dbRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *dbRepository) DeleteUser(id string) error {
	db := r.db.Delete("id = ?", id).Delete(&models.User{})
	if err := db.Error; err != nil {
		return err
	}
	if db.RowsAffected > 0 {
		err := r.DeleteFile(&models.File{
			Owner: id,
			IsDir: true,
			FID:   id,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *dbRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *dbRepository) GetUser(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *dbRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *dbRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *dbRepository) GetUserList(limit int, offset int) ([]*models.User, error) {
	var users = make([]*models.User, 0)
	db := r.db

	if limit != 0 {
		db.Limit(limit)
	}

	if offset != 0 {
		db.Offset(offset)
	}

	err := db.Find(&users).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *dbRepository) CheckUserExists(user *models.User) (string, error) {
	var conflictingUser models.User

	db := r.db
	db = db.Where("username = ? OR email = ?", user.Username, user.Email)
	err := db.First(&conflictingUser).Error
	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	if user.Username == conflictingUser.Username {
		return "username", nil
	}
	return "email", nil
}
