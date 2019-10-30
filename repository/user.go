package repository

import (
	"github.com/phantom-atom/file-explorer/models"
)

//UserRepository 用户仓库接口
type UserRepository interface {
	CreateUser(*models.User) error
	DeleteUser(id string) error
	UpdateUser(*models.User) error
	GetUser(id string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserList(limit int, offset int) ([]*models.User, error)
	CheckUserExists(*models.User) (string, error)
}
