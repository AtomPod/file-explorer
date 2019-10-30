package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/phantom-atom/file-explorer/models"
	"github.com/phantom-atom/file-explorer/services"
	"github.com/phantom-atom/file-explorer/web/forms"
)

func userErrorToAPIResult(err error) *APIResult {
	userErr, ok := err.(*services.UserError)
	if !ok {
		return Internal(err, nil)
	}

	switch userErr.Err {
	case services.ErrUserNotFound:
		return NotFound(err, nil)
	case services.ErrEmailAlreadyExists, services.ErrUsernameAlreadyExists:
		return AlreadyExists(err, nil)
	case services.ErrUserFieldIsInvalid, services.ErrVerificationCodeIsInvalid:
		return InvalidArgument(err, nil)
	case services.ErrIncorrectUnameOrPWD:
		return InvalidArgument(err, nil)
	default:
		return Internal(err, nil)
	}
}

//UserEmailCodeGenerator 用户邮箱验证码生成API
//POST /api/v1/user/email_code
func (api *API) UserEmailCodeGenerator(c *gin.Context, form *forms.UserEmailCode) *APIResult {
	_, err := api.userServ.CreateEmailVerificationCode(&services.UserRegisterCodeParams{
		Email: form.Email,
	})

	if err != nil {
		return userErrorToAPIResult(err)
	}
	return OK("Email verification code has been sent", nil)
}

//UserRegister 用户注册API
//POST /api/v1/user/register
func (api *API) UserRegister(c *gin.Context, form *forms.UserRegister) *APIResult {
	newUser, err := api.userServ.RegisterUser(&services.UserRegisterParams{
		Username: form.Username,
		Email:    form.Email,
		Password: form.Password,
		Role:     models.UserRoleUser,
		Code:     form.Code,
	})

	if err != nil {
		return userErrorToAPIResult(err)
	}

	return OK(newUser, nil)
}

//UserLogin 用户登录API
//POST /api/v1/user/login
func (api *API) UserLogin(c *gin.Context, form *forms.UserLogin) *APIResult {
	token, err := api.userServ.LoginUser(&services.UserLoginParams{
		Identity:    form.Username,
		Certificate: form.Password,
	})

	if err != nil {
		return userErrorToAPIResult(err)
	}

	return OK(token, nil)
}

//UserResetPasswordCode 用户修改密码验证码请求API
//POST /api/v1/user/password/reset_code
func (api *API) UserResetPasswordCode(c *gin.Context, form *forms.UserEmailCode) *APIResult {
	_, err := api.userServ.CreateResetPWDVerificationCode(&services.UserResetPWDCodeParams{
		Email: form.Email,
	})

	if err != nil {
		return userErrorToAPIResult(err)
	}

	return OK("Reset password verification code has been sent", nil)
}

//UserResetPassword 用户修改密码请求API
//POST /api/v1/user/password/reset
func (api *API) UserResetPassword(c *gin.Context, form *forms.UserResetPassword) *APIResult {
	err := api.userServ.ResetPassword(&services.UserResetPasswordParams{
		Email:    form.Email,
		Code:     form.Code,
		Password: form.Password,
	})

	if err != nil {
		return userErrorToAPIResult(err)
	}

	return OK("The password has been modified", nil)
}

//UserCurrentInfo 当前用户信息请求API
//POST /api/v1/user/current
func (api *API) UserCurrentInfo(c *gin.Context) *APIResult {
	userID := c.GetString("userID")
	user, err := api.userServ.GetUserByID(userID)

	if err != nil {
		return userErrorToAPIResult(err)
	}

	return OK(user, nil)
}
