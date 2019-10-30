package password

import (
	"github.com/phantom-atom/file-explorer/internal/log"
	"golang.org/x/crypto/bcrypt"
)

//CreateHashPassword 创建Hash密码
func CreateHashPassword(pwd string) string {
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Info("msg", "occur a error when create hash passowrd", "err", err.Error())
		return ""
	}
	return string(hashPwd)
}

//CompareHashPassword 比较Hash密码与原密码
func CompareHashPassword(pwd string, hashPwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(pwd)) == nil
}
