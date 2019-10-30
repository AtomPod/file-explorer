package security_test

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/phantom-atom/file-explorer/internal/security"
)

func TestJWT(t *testing.T) {
	claims := &jwt.StandardClaims{}
	claims.Id = "15"
	token, _ := security.GenerateJWT(security.SecretKey("hello"), claims)
	_, err := security.DecodeJWT(security.SecretKey("hello"), token, new(jwt.StandardClaims))
	if err != nil {
		t.Error(err)
	}
}
