package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/jinzhu/gorm"
)

var (
	//UserRoleAdmin 用户权限(管理)
	UserRoleAdmin = "admin"
	//UserRoleUser 用户权限(用户)
	UserRoleUser = "user"
)

//User 用户模型
type User struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`
	Username  string     `gorm:"column:username;unique_index" json:"username"`
	Password  string     `gorm:"column:password" json:"-"`
	Email     string     `gorm:"column:email;unique_index" json:"email"`
	Role      string     `gorm:"column:role" json:"role"`
}

//BeforeCreate 添加UUID主键
func (user *User) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	return scope.SetColumn("ID", uuid.String())
}

//UserTokenClaims 用户JWT Claims
type UserTokenClaims struct {
	jwt.StandardClaims
	ID   string `json:"id"`
	Role string `json:"role"`
}
