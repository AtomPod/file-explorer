package validation

import (
	"regexp"

	"github.com/phantom-atom/file-explorer/internal/log"
)

//IsEmail 验证是否为邮箱
func IsEmail(email string) bool {
	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+\\.[a-zA-Z0-9_-]+$", email)
	if err != nil {
		log.Error("msg", "occur a error when matching a email", "error", err.Error())
		return false
	}
	return matched
}

//IsUsername 验证是否为用户名
func IsUsername(username string) bool {
	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]{4,16}$", username)
	if err != nil {
		log.Error("msg", "occur a error when matching a username", "error", err.Error())
		return false
	}
	return matched
}
