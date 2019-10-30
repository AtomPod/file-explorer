package v1

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

var (
	ginContextTyp = reflect.TypeOf(&gin.Context{})
	apiResultTyp  = reflect.TypeOf(&APIResult{})
)

//APIWithModel 转换(*gin.Context , model) *apiResult到GinFunc
func APIWithModel(f interface{}) GinFunc {
	val := reflect.ValueOf(f)
	typ := val.Type()

	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = val.Type()
	}

	if typ.Kind() != reflect.Func {
		panic("api_with_model: argument must be a function")
	}

	if typ.NumIn() != 2 || typ.NumOut() != 1 {
		panic("api_with_model: arguments must be a function <func(gin.Context,model) *apiResult>")
	}

	inTyp1 := typ.In(0)
	outTyp := typ.Out(0)

	if inTyp1 != ginContextTyp || outTyp != apiResultTyp {
		panic("api_with_model: arguments must be a function <func(gin.Context,model) *apiResult>")
	}
	return withModelFunc(val, typ.In(1))
}

func withModelFunc(funcVal reflect.Value, modelTyp reflect.Type) GinFunc {
	var isPtr = false

	if !funcVal.IsValid() {
		panic("api_with_model: invalid function")
	}

	if modelTyp.Kind() == reflect.Ptr {
		modelTyp = modelTyp.Elem()
		isPtr = true
	}

	return GinFunc(func(c *gin.Context) *APIResult {
		modelVal := reflect.New(modelTyp)

		c.ShouldBindUri(modelVal.Interface())
		if err := c.ShouldBind(modelVal.Interface()); err != nil {
			return InvalidArgument(err, nil)
		}

		if !isPtr {
			modelVal = modelVal.Elem()
		}

		arguments := make([]reflect.Value, 2)
		arguments[0] = reflect.ValueOf(c)
		arguments[1] = modelVal
		results := funcVal.Call(arguments)
		apiResultVal := results[0]
		apiResult := apiResultVal.Interface()
		return apiResult.(*APIResult)
	})
}
