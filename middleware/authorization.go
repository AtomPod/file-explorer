package middleware

import (
	"errors"
	"strings"

	v1 "github.com/phantom-atom/file-explorer/web/api/v1"

	"github.com/dgrijalva/jwt-go"
	"github.com/phantom-atom/file-explorer/config"
	"github.com/phantom-atom/file-explorer/internal/security"

	"github.com/gin-gonic/gin"
	"github.com/phantom-atom/file-explorer/models"
	"github.com/phantom-atom/file-explorer/services"
)

const (
	authHeader = "X-REQUEST-TOKEN"
)

var (
	errUserNotAuthenticated = errors.New("user not authenticated")
	errPermissionDenied     = errors.New("permission denied")
)

//Authorization 认证中间件
type Authorization struct {
	service *services.UserService
	config  func() *config.Config
}

//NewAuthorization 新建认证中间件
func NewAuthorization(services *services.UserService,
	configFunc func() *config.Config) *Authorization {
	return &Authorization{
		service: services,
		config:  configFunc,
	}
}

//HandlerFunc 认证中间件
func (a *Authorization) HandlerFunc(role string) v1.GinFunc {
	return func(c *gin.Context) *v1.APIResult {
		user, token, err := a.verifyAndTest(c, role)
		if err == nil {
			a.setUserToContext(c, user, token)
			c.Next()
			return nil
		}

		switch err {
		case errUserNotAuthenticated:
			return v1.Unauthuenticated(err, nil)
		case errPermissionDenied:
			return v1.PermissionDenied(err, nil)
		default:
			return v1.Internal(err, nil)
		}
	}
}

func (a *Authorization) tokenFromBearer(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	BTA := strings.Split(bearerToken, " ")
	if len(BTA) < 2 {
		return ""
	}

	if strings.ToLower(BTA[0]) == "bearer" {
		return BTA[1]
	}
	return ""
}

func (a *Authorization) tokenFromHeader(c *gin.Context) string {
	token := c.GetHeader(authHeader)
	return token
}

func (a *Authorization) tokenFromQuery(c *gin.Context) string {
	token := c.Query("token")
	return token
}

func (a *Authorization) tokenFromAny(c *gin.Context) string {
	token := a.tokenFromBearer(c)

	if token == "" {
		token = a.tokenFromHeader(c)
	}

	if token == "" {
		token = a.tokenFromQuery(c)
	}
	return token
}

func (a *Authorization) verifyAndTest(c *gin.Context, role string) (*models.User, string, error) {
	token := a.tokenFromAny(c)

	if token == "" {
		return nil, "", errUserNotAuthenticated
	}

	user, err := a.decodeUserToken(token)

	if err != nil {
		return nil, "", err
	}

	if user == nil {
		return nil, "", errUserNotAuthenticated
	}

	if role != user.Role {
		return nil, "", errPermissionDenied
	}

	return user, token, nil
}

func (a *Authorization) decodeUserToken(tokenString string) (*models.User, error) {
	jwtConf := a.config().UserService.JWT

	var claim *models.UserTokenClaims
	token, err := security.DecodeJWT(security.SecretKey(jwtConf.SigningKey), tokenString, new(models.UserTokenClaims))

	if token.Valid {
		var ok bool
		if claim, ok = token.Claims.(*models.UserTokenClaims); !ok {
			return nil, nil
		}
	} else if _, ok := err.(*jwt.ValidationError); ok {
		return nil, nil
	} else {
		return nil, err
	}

	if claim.Subject != "user_auth" {
		return nil, nil
	}

	user := &models.User{}
	user.ID = claim.ID
	user.Role = claim.Role
	return user, nil
}

func (a *Authorization) setUserToContext(c *gin.Context, user *models.User, token string) {
	c.Set("user", user)
	c.Set("token", token)
	c.Set("userID", user.ID)
	c.Set("userRole", user.Role)
}

//GetUserFromContext 从Context中获取用户
func GetUserFromContext(c *gin.Context) *models.User {
	v, exists := c.Get("user")
	if !exists {
		return nil
	}
	if user, ok := v.(*models.User); ok {
		return user
	}
	return nil
}

//GetTokenFromContext 从Context中获取token
func GetTokenFromContext(c *gin.Context) string {
	v, exists := c.Get("token")
	if !exists {
		return ""
	}
	if token, ok := v.(string); ok {
		return token
	}
	return ""
}
