package security

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

//SecretKey 安全Key
type SecretKey []byte

//GenerateJWT 创建新的JWT
func GenerateJWT(secretKey SecretKey, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encoded, e := token.SignedString([]byte(secretKey))
	if e != nil {
		return "", e
	}
	return encoded, nil
}

//DecodeJWT 解析并验证JWT
func DecodeJWT(secretKey SecretKey, token string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unknow token method")
		}
		return []byte(secretKey), nil
	})
}
