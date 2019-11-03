package services

import (
	"errors"
	"time"

	"github.com/phantom-atom/file-explorer/internal/log"

	"github.com/phantom-atom/file-explorer/mailer"

	"github.com/phantom-atom/file-explorer/internal/utils/random"

	"github.com/phantom-atom/file-explorer/internal/security"

	"github.com/phantom-atom/file-explorer/internal/utils/validation"

	"github.com/phantom-atom/file-explorer/internal/utils/password"

	"github.com/phantom-atom/file-explorer/config"
	"github.com/phantom-atom/file-explorer/internal/locker"
	"github.com/phantom-atom/file-explorer/models"
	"github.com/phantom-atom/file-explorer/repository"
)

var (
	//ErrUsernameAlreadyExists 用户名已经存在
	ErrUsernameAlreadyExists = errors.New("用户名已经存在")
	//ErrEmailAlreadyExists 邮箱已经存在
	ErrEmailAlreadyExists = errors.New("邮箱已经存在")
	//ErrIncorrectUnameOrPWD 用户名或密码错误
	ErrIncorrectUnameOrPWD = errors.New("用户名或密码错误")
	//ErrUserFieldIsInvalid 用户字段无效
	ErrUserFieldIsInvalid = errors.New("用户字段无效")
	//ErrEmailNotFound 邮箱不存在
	ErrEmailNotFound = errors.New("邮箱不存在")
	//ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("用户不存在")
	//ErrVerificationCodeIsInvalid 确认码无效
	ErrVerificationCodeIsInvalid = errors.New("验证码无效")
)

//UserError 用户错误
type UserError struct {
	Op  string //操作
	Err error  //错误
}

//Error error接口
func (ue *UserError) Error() string {
	return "[" + ue.Op + "]" + ue.Err.Error()
}

//NewUserError 创建UserError
func NewUserError(op string, err error) *UserError {
	return &UserError{
		Op:  op,
		Err: err,
	}
}

//UserService 用户服务
type UserService struct {
	config      func() *config.Config
	uuid        func() string
	now         func() time.Time
	dataContext repository.DataContext
	mailer      *mailer.Mailer
	namedLocker locker.NamedLocker
}

//NewUserService 创建用户服务
func NewUserService(
	configFunc func() *config.Config,
	uuid func() string,
	now func() time.Time,
	dataContext repository.DataContext,
	mailer *mailer.Mailer,
	namedLocker locker.NamedLocker) *UserService {
	return &UserService{
		config:      configFunc,
		uuid:        uuid,
		now:         now,
		dataContext: dataContext,
		mailer:      mailer,
		namedLocker: namedLocker,
	}
}

//UserRegisterCodeParams 用户注册码参数
type UserRegisterCodeParams struct {
	Email string `json:"email"`
}

//CreateEmailVerificationCode 创建邮箱确认码
func (us *UserService) CreateEmailVerificationCode(params *UserRegisterCodeParams) (string, error) {
	if params == nil {
		return "", invalidArgument("UserService", "params", "CreateRegisterCode")
	}

	emailCodeConf := &us.config().UserService.VerificationCode.Email
	code, err := us.createVerificationCode(
		params.Email,
		emailCodeConf.TypeName,
		emailCodeConf.Expiration,
		emailCodeConf.MailTemplName,
		true,
	)

	return code, err
}

//UserResetPWDCodeParams 用户修改密码验证码参数
type UserResetPWDCodeParams struct {
	Email string `json:"email"`
}

//CreateResetPWDVerificationCode 创建修改密码验证码
func (us *UserService) CreateResetPWDVerificationCode(params *UserResetPWDCodeParams) (string, error) {
	if params == nil {
		return "", invalidArgument("UserService", "params", "CreateResetPWDVerificationCode")
	}

	resestPasswordCodeConf := &us.config().UserService.VerificationCode.ResetPassword
	code, err := us.createVerificationCode(
		params.Email,
		resestPasswordCodeConf.TypeName,
		resestPasswordCodeConf.Expiration,
		resestPasswordCodeConf.MailTemplName,
		false,
	)
	return code, err
}

//UserResetPasswordParams 用户修改密码参数
type UserResetPasswordParams struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

//ResetPassword 修改密码
func (us *UserService) ResetPassword(params *UserResetPasswordParams) error {
	if params == nil {
		return invalidArgument("UserService", "params", "ResetPassword")
	}

	if params.Password == "" {
		return NewUserError("reset password", ErrUserFieldIsInvalid)
	}

	if params.Code == "" {
		return NewUserError("reset password", ErrVerificationCodeIsInvalid)
	}

	codeRepository, err := us.dataContext.VerificationCode()
	if err != nil {
		return err
	}

	userRepository, err := us.dataContext.User()
	if err != nil {
		return err
	}

	resestPasswordCodeConf := &us.config().UserService.VerificationCode.ResetPassword
	code, err := codeRepository.GetVerificationCodeByCodeTypeTarget(
		resestPasswordCodeConf.TypeName,
		params.Code,
		params.Email,
	)

	if err != nil {
		return err
	}

	if code == nil {
		return NewUserError("reset password", ErrVerificationCodeIsInvalid)
	}

	if err := codeRepository.DeleteVerificationCode(
		resestPasswordCodeConf.TypeName,
		params.Code,
		params.Email,
	); err != nil {
		return err
	}

	user, err := userRepository.GetUserByEmail(params.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return NewUserError("reset password", ErrEmailNotFound)
	}

	user.Password = password.CreateHashPassword(params.Password)
	if err := userRepository.UpdateUser(user); err != nil {
		return err
	}
	return nil
}

//UserRegisterParams 用户注册参数
type UserRegisterParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Code     string `json:"code"`
}

//RegisterUser 注册用户
func (us *UserService) RegisterUser(params *UserRegisterParams) (*models.User, error) {
	if params == nil {
		return nil, invalidArgument("UserService", "params", "RegisterUser")
	}

	if !validation.IsUsername(params.Username) {
		return nil, NewUserError("register", ErrUserFieldIsInvalid)
	}

	if params.Password == "" {
		return nil, NewUserError("register", ErrUserFieldIsInvalid)
	}

	if !validation.IsEmail(params.Email) {
		return nil, NewUserError("register", ErrUserFieldIsInvalid)
	}

	if params.Code == "" {
		return nil, NewUserError("register", ErrUserFieldIsInvalid)
	}

	if params.Role == "" {
		params.Role = models.UserRoleUser
	}

	codeRepository, err := us.dataContext.VerificationCode()
	if err != nil {
		return nil, err
	}

	userRepository, err := us.dataContext.User()
	if err != nil {
		return nil, err
	}

	emailCodeConf := &us.config().UserService.VerificationCode.Email
	code, err := codeRepository.GetVerificationCodeByCodeTypeTarget(
		emailCodeConf.TypeName,
		params.Code,
		params.Email,
	)

	if err != nil {
		return nil, err
	}

	if code == nil {
		return nil, NewUserError("register", ErrVerificationCodeIsInvalid)
	}

	if err := codeRepository.DeleteVerificationCode(
		emailCodeConf.TypeName,
		params.Code,
		params.Email,
	); err != nil {
		return nil, err
	}

	newUser := &models.User{
		Username: params.Username,
		Email:    params.Email,
	}

	us.namedLocker.Lock("$file-explorer:user-service$")
	defer us.namedLocker.UnLock("$file-explorer:user-service$")
	confling, err := userRepository.CheckUserExists(newUser)
	if err != nil {
		return nil, err
	}

	switch confling {
	case "username":
		return nil, NewUserError("register", ErrUsernameAlreadyExists)
	case "email":
		return nil, NewUserError("register", ErrEmailAlreadyExists)
	}

	newUser.Password = password.CreateHashPassword(params.Password)
	newUser.Role = params.Role
	if err := userRepository.CreateUser(newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

//UserLoginParams 用户登录参数
type UserLoginParams struct {
	Identity    string `json:"identity"`    //用户名
	Certificate string `json:"certificate"` //密码
}

//LoginUser 用户登录
func (us *UserService) LoginUser(params *UserLoginParams) (*models.Token, error) {
	if params == nil {
		return nil, invalidArgument("UserService", "params", "LoginUser")
	}

	if params.Identity == "" {
		return nil, invalidArgument("UserService", "params.Username", "LoginUser")
	}

	if params.Certificate == "" {
		return nil, invalidArgument("UserService", "params.Password", "LoginUser")
	}

	userRepository, err := us.dataContext.User()
	if err != nil {
		return nil, err
	}

	var matchedUser *models.User

	if validation.IsEmail(params.Identity) {
		matchedUser, err = userRepository.GetUserByEmail(params.Identity)
	} else if validation.IsUsername(params.Identity) {
		matchedUser, err = userRepository.GetUserByUsername(params.Identity)
	} else {
		return nil, NewUserError("login", ErrIncorrectUnameOrPWD)
	}

	if err != nil {
		return nil, err
	}

	if matchedUser == nil {
		return nil, NewUserError("login", ErrIncorrectUnameOrPWD)
	}

	if !password.CompareHashPassword(params.Certificate, matchedUser.Password) {
		return nil, NewUserError("login", ErrIncorrectUnameOrPWD)
	}

	return us.generateUserToken(matchedUser)
}

//GetUserByID 通过ID获取用户信息
func (us *UserService) GetUserByID(id string) (*models.User, error) {
	if id == "" {
		return nil, invalidArgument("UserService", "id", "GetUserByID")
	}

	userRepository, err := us.dataContext.User()
	if err != nil {
		return nil, err
	}

	user, err := userRepository.GetUser(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, NewUserError("get", ErrUserNotFound)
	}

	return user, nil
}

func (us *UserService) generateUserToken(user *models.User) (*models.Token, error) {
	jwtConf := us.config().UserService.JWT

	claim := models.UserTokenClaims{
		ID:   user.ID,
		Role: user.Role,
	}
	claim.Subject = "user_auth"
	claim.ExpiresAt = us.now().Add(jwtConf.Expire).Unix()
	tokenString, err := security.GenerateJWT(security.SecretKey(jwtConf.SigningKey), &claim)

	if err != nil {
		return nil, err
	}

	return &models.Token{
		Token:  tokenString,
		Expire: int64(jwtConf.Expire.Seconds()),
	}, nil
}

func (us *UserService) createVerificationCode(
	email string,
	typ string,
	expiration time.Duration,
	tmplName string,
	emailNotExists bool) (string, error) {
	if !validation.IsEmail(email) {
		return "", NewUserError("verification code", ErrUserFieldIsInvalid)
	}
	
	userRepository, err := us.dataContext.User()
	if err != nil {
		return "", err
	}

	codeRepository, err := us.dataContext.VerificationCode()
	if err != nil {
		return "", err
	}

	confling, err := userRepository.CheckUserExists(&models.User{
		Email: email,
	})

	if err != nil {
		return "", err
	}

	if confling == "email" && emailNotExists {
		return "", NewUserError("verification code", ErrEmailAlreadyExists)
	} else if confling != "email" && !emailNotExists {
		return "", NewUserError("verification code", ErrEmailNotFound)
	}

	config := us.config()
	emailConf := &config.Email
	code := random.DigitString(6)
	if err := codeRepository.CreateVerificationCode(&models.VerificationCode{
		Type:       typ,
		Target:     email,
		Code:       code,
		Expiration: expiration,
	}); err != nil {
		return "", err
	}

	tpl := mailer.GetMailTemplate(tmplName)
	if tpl != nil {
		data := &struct {
			ServerName string
			Email      string
			Code       string
			Expiration time.Duration
		}{config.ServerName, email, code, expiration}

		body, err := tpl.Body(data)
		if err != nil {
			return "", err
		}

		if err := us.mailer.Send(
			emailConf.Username,
			[]string{email},
			tpl.Subject,
			tpl.ContentType,
			body); err != nil {
			return "", err
		}
	} else {
		log.Warn("msg", "mail template file does not exist", "name", tmplName)
	}
	return code, nil
}
