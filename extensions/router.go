package extensions

import "github.com/gin-gonic/gin"

var routerInitializers = []RouterInitializer{}

// RouterInitializer is the interface what adds router processes specifically before and after router registration
// * InitializeEarly adds router processes before router registration
// * InitializeLate  adds router processes after router registration
type RouterInitializer interface {
	InitializeEarly(r *gin.Engine) error
	InitializeLate(r *gin.Engine) error
}

// RegisterRouterInitializer registers an initializer used in the router logic
func RegisterRouterInitializer(initializer RouterInitializer) {
	routerInitializers = append(routerInitializers, initializer)
}

// RegisteredRouterInitializers returns the registered router initializers
func RegisteredRouterInitializers() []RouterInitializer {
	result := []RouterInitializer{}
	result = append(result, routerInitializers...)
	return result
}
