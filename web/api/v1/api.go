package v1

import (
	"net/http"

	"github.com/phantom-atom/file-explorer/internal/log"

	"github.com/phantom-atom/file-explorer/services"

	"github.com/phantom-atom/file-explorer/config"

	"github.com/gin-gonic/gin"
)

//ErrorType 错误类型
type ErrorType int

//Error 错误码
const (
	ErrNone               ErrorType = 0
	ErrFailedPrecondition ErrorType = 1
	ErrUnauthuenticated   ErrorType = 2
	ErrInternal           ErrorType = 3
	ErrAlreadyExists      ErrorType = 4
	ErrInvalidArgument    ErrorType = 5
	ErrPermissionDenied   ErrorType = 6
	ErrNotFound           ErrorType = 7
)

var errorName = []string{
	"",
	"failed_precondition",
	"unauthuenticated",
	"internal",
	"already_exists",
	"invalid_argument",
	"permisttion_denied",
	"not_found",
}

//ErrorTypeToName 转换错误码为字符串
func ErrorTypeToName(e ErrorType) string {
	if e > ErrNotFound {
		return "unknow"
	}
	return errorName[e]
}

type responseError struct {
	Code    int    `json:"code"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

type response struct {
	Status string         `json:"status"`
	Data   interface{}    `json:"data,omitempty"`
	Error  *responseError `json:"error,omitempty"`
}

//APIError API错误
type APIError struct {
	Code ErrorType `json:"code"`
	Err  error     `json:"error"`
}

//APIResult API返回值
type APIResult struct {
	Data      interface{}
	Responder func(c *gin.Context, data interface{}) error
	Error     *APIError
}

//GinFunc GIN转换函数
type GinFunc func(*gin.Context) *APIResult

//API API接口
type API struct {
	config   func() *config.Config
	fileServ *services.FileService
	userServ *services.UserService
}

//NewAPI 创建API
func NewAPI(configFunc func() *config.Config,
	fileServ *services.FileService,
	userServ *services.UserService) *API {
	return &API{
		config:   configFunc,
		fileServ: fileServ,
		userServ: userServ,
	}
}

//Gin 转换为gin使用的函数
func (api *API) Gin(f GinFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := f(c)
		if result != nil {
			if result.Responder != nil {
				if err := result.Responder(c, result.Data); err != nil {
					api.respondError(c, &APIError{
						Code: ErrInternal,
						Err:  err,
					}, nil)
				}
			} else if result.Error != nil {
				api.respondError(c, result.Error, result.Data)
			} else if result.Data != nil {
				api.respond(c, result.Data)
			} else {
				api.successRespond(c)
			}
			c.Abort()
		}
	}
}

func (api *API) successRespond(c *gin.Context) {
	resp := &response{
		Status: "success",
	}
	c.JSON(http.StatusOK, resp)
}

func (api *API) respond(c *gin.Context, data interface{}) {
	resp := &response{
		Status: "success",
		Data:   data,
	}

	c.JSON(http.StatusOK, resp)
}

func (api *API) respondError(c *gin.Context, e *APIError, data interface{}) {
	resp := &response{
		Status: "error",
		Data:   data,
		Error: &responseError{
			Code:    int(e.Code),
			Type:    ErrorTypeToName(e.Code),
			Message: e.Err.Error(),
		},
	}
	statusCode := api.errorTypeToHTTPStatusCode(e.Code)

	if statusCode == http.StatusInternalServerError {
		mode := api.config().Mode
		//在release模式下，不打印InternalServerError
		if mode == "release" {
			resp.Error.Message = resp.Error.Type
		}
		log.Error("msg", "occur an error when request", "error", e.Err.Error())
	}

	c.JSON(statusCode, resp)
}

func (api *API) errorTypeToHTTPStatusCode(e ErrorType) int {
	switch e {
	case ErrNone:
		return http.StatusOK
	case ErrFailedPrecondition,
		ErrAlreadyExists,
		ErrInvalidArgument:
		return http.StatusBadRequest
	case ErrUnauthuenticated:
		return http.StatusUnauthorized
	case ErrNotFound:
		return http.StatusNotFound
	case ErrPermissionDenied:
		return http.StatusForbidden
	case ErrInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

//OK OK
func OK(data interface{}, responder func(c *gin.Context, data interface{}) error) *APIResult {
	return &APIResult{
		Data:      data,
		Responder: responder,
	}
}

//FailedPrecondition FailedPrecondition
func FailedPrecondition(err error, data interface{}) *APIResult {
	return &APIResult{
		Data: data,
		Error: &APIError{
			Code: ErrFailedPrecondition,
			Err:  err,
		},
	}
}

//Unauthuenticated  Unauthuenticated
func Unauthuenticated(err error, data interface{}) *APIResult {
	return &APIResult{
		Data: data,
		Error: &APIError{
			Code: ErrUnauthuenticated,
			Err:  err,
		},
	}
}

//Internal Internal
func Internal(err error, data interface{}) *APIResult {
	return &APIResult{
		Data: data,
		Error: &APIError{
			Code: ErrInternal,
			Err:  err,
		},
	}
}

//AlreadyExists AlreadyExists
func AlreadyExists(err error, data interface{}) *APIResult {
	return &APIResult{
		Data: data,
		Error: &APIError{
			Code: ErrAlreadyExists,
			Err:  err,
		},
	}
}

//InvalidArgument InvalidArgument
func InvalidArgument(err error, data interface{}) *APIResult {
	return &APIResult{
		Data: data,
		Error: &APIError{
			Code: ErrInvalidArgument,
			Err:  err,
		},
	}
}

//PermissionDenied PermissionDenied
func PermissionDenied(err error, data interface{}) *APIResult {
	return &APIResult{
		Data: data,
		Error: &APIError{
			Code: ErrPermissionDenied,
			Err:  err,
		},
	}
}

//NotFound NotFound
func NotFound(err error, data interface{}) *APIResult {
	return &APIResult{
		Data: data,
		Error: &APIError{
			Code: ErrNotFound,
			Err:  err,
		},
	}
}
