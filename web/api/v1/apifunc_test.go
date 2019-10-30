package v1_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	v1 "github.com/phantom-atom/file-explorer/web/api/v1"
)

func TestCallbackCorrectType(t *testing.T) {
	defer func() {
		recover()
	}()
	v1.APIWithModel(func(c *gin.Context) {})
	t.Errorf("argument [func(c *gin.Context)] pass")
}

type testModel struct{}

func TestCallbackCorrectType2(t *testing.T) {
	defer func() {
		recover()
	}()
	v1.APIWithModel(func(c *gin.Context, _1 *testModel) {})
	t.Errorf("argument [func(c *gin.Context, _1 *testModel)] pass")
}

func TestCallbackCorrectType3(t *testing.T) {
	defer func() {
		recover()
	}()
	v1.APIWithModel(func(c *gin.Context) *v1.APIResult {
		return nil
	})
	t.Errorf("argument [func(c *gin.Context) *v1.APIResult] pass")
}

func TestCallbackCorrectType4(t *testing.T) {
	defer func() {
		recover()
	}()
	v1.APIWithModel(func(_1 *testModel) *v1.APIResult {
		return nil
	})
	t.Errorf("argument [func(_1 *testModel) *v1.APIResult] pass")
}

func TestCallbackType(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal(err)
		}
	}()
	v1.APIWithModel(func(c *gin.Context, _1 *testModel) *v1.APIResult {
		return nil
	})
}
