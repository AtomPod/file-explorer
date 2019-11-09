package register

import (
	"github.com/gin-gonic/gin"
	"github.com/phantom-atom/file-explorer/internal/utils/httputil"
	"github.com/phantom-atom/file-explorer/middleware"
	"github.com/phantom-atom/file-explorer/models"
	v1 "github.com/phantom-atom/file-explorer/web/api/v1"
)

//APIGINRegister 注册gin路由
func APIGINRegister(
	api *v1.API,
	router *gin.RouterGroup,
	authMiddleware *middleware.Authorization,
) {
	ginAPIFunc := func(f interface{}) gin.HandlerFunc {
		return api.Gin(v1.APIWithModel(f))
	}

	router.Use(func(c *gin.Context) {
		httputil.CORS(c.Writer, c.Request)
		c.Next()
	})

	router.OPTIONS("/*path", api.Gin(func(c *gin.Context) *v1.APIResult {
		return v1.OK(nil, nil)
	}))

	authAPIMiddleware := api.Gin(authMiddleware.HandlerFunc(models.UserRoleUser))
	apiRouter := router.Group("/api/v1")
	userRouter := apiRouter.Group("/user")
	userRouter.POST("/register", ginAPIFunc(api.UserRegister))
	userRouter.POST("/login", ginAPIFunc(api.UserLogin))
	userRouter.POST("/email_code", ginAPIFunc(api.UserEmailCodeGenerator))
	userRouter.POST("/password/reset_code", ginAPIFunc(api.UserResetPasswordCode))
	userRouter.POST("/password/reset", ginAPIFunc(api.UserResetPassword))
	userRouter.GET("/current", authAPIMiddleware, api.Gin(api.UserCurrentInfo))

	fileRouter := apiRouter.Group("/file")
	fileRouter.Use(authAPIMiddleware)
	fileRouter.POST("/upload", ginAPIFunc(api.FileUpload))
	fileRouter.POST("/mkdir", ginAPIFunc(api.FileMkdir))
	fileRouter.PUT("/:id/rename", ginAPIFunc(api.FileRename))
	fileRouter.PUT("/:id/move", ginAPIFunc(api.FileMove))
	fileRouter.GET("/", api.Gin(api.FileGetRootList))
	fileRouter.GET("/:id", ginAPIFunc(api.FileDownload))
	fileRouter.GET("/:id/info", ginAPIFunc(api.FileGetInfo))
	fileRouter.GET("/:id/list", ginAPIFunc(api.FileGetList))
	fileRouter.DELETE("/:id", ginAPIFunc(api.FileDelete))
}
