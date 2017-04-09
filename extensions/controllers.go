package extensions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var controllers = []Controller{}

const (
	MethodGet     = 1
	MethodPost    = 2
	MethodPut     = 3
	MethodDelete  = 4
	MethodPatch   = 5
	MethodOptions = 6
)

var methodNameMap = map[int]string{
	MethodGet:     http.MethodGet,
	MethodPost:    http.MethodPost,
	MethodPut:     http.MethodPut,
	MethodDelete:  http.MethodDelete,
	MethodPatch:   http.MethodPatch,
	MethodOptions: http.MethodOptions,
}

type Controller interface {
	ResourceName() string
	RouteMap() map[int]map[string]gin.HandlerFunc
}

func GetMethodName(method int) string {
	return methodNameMap[method]
}

func GetResourceSingleUrl(resourceName string) string {
	return fmt.Sprintf("/%ss/:id", resourceName)
}

func GetResourceMultiUrl(resourceName string) string {
	return fmt.Sprintf("/%ss", resourceName)
}

func RegisterController(controller Controller) {
	controllers = append(controllers, controller)
}

func GetControllers() []Controller {
	result := []Controller{}
	result = append(result, controllers...)
	return result
}
