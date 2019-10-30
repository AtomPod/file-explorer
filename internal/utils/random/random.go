package random

import (
	"crypto/rand"
	mathRand "math/rand"

	"github.com/phantom-atom/file-explorer/internal/log"
)

//String 随机一个len长度的字符串
func String(len int) string {
	byt := make([]byte, len)
	if _, err := rand.Read(byt); err != nil {
		log.Warn("msg", "occur a error when random string", "err", err.Error())
		return ""
	}
	return string(byt)
}

//DigitString 随机一个len长度的数字字符串
func DigitString(len int) string {
	byt := make([]byte, len)

	for i := 0; i < len; i++ {
		n := mathRand.Intn(10)
		byt[i] = byte(n + '0')
	}
	return string(byt)
}
