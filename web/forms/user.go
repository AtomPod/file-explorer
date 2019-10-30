package forms

//UserRegister 用户注册表单
type UserRegister struct {
	Username string `json:"username" form:"username" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,alphanum"`
	Code     string `json:"code" form:"code" binding:"required,len=6,numeric"`
}

//UserLogin 用户登录表单
type UserLogin struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required,alphanum"`
}

//UserEmailCode 用户验证码
type UserEmailCode struct {
	Email string `json:"email" form:"email" binding:"required,email"`
}

//UserResetPassword 用户修改密码表单
type UserResetPassword struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,alphanum"`
	Code     string `json:"code" form:"code" binding:"required,len=6,numeric"`
}
