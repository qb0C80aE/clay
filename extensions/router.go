package extensions

import "github.com/gin-gonic/gin"

var routerInitializers = []RouterInitializer{}

type RouterInitializer interface {
	InitializeEarly(r *gin.Engine) error
	InitializeLate(r *gin.Engine) error
}

func RegisterRouterInitializer(initializer RouterInitializer) {
	routerInitializers = append(routerInitializers, initializer)
}

func GetRouterInitializers() []RouterInitializer {
	result := []RouterInitializer{}
	result = append(result, routerInitializers...)
	return result
}
