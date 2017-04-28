package router

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/controllers"
	"github.com/qb0C80aE/clay/extensions"
)

// Initialize initializes the router
func Initialize(r *gin.Engine) {
	routerInitializers := extensions.RegisteredRouterInitializers()

	for _, initializer := range routerInitializers {
		initializer.InitializeEarly(r)
	}

	r.GET("/", controllers.APIEndpoints)

	api := r.Group("")
	{
		methodFunctionMap := map[int]func(string, ...gin.HandlerFunc) gin.IRoutes{
			extensions.MethodGet:     api.GET,
			extensions.MethodPost:    api.POST,
			extensions.MethodPut:     api.PUT,
			extensions.MethodDelete:  api.DELETE,
			extensions.MethodPatch:   api.PATCH,
			extensions.MethodOptions: api.OPTIONS,
		}

		controllers := extensions.RegisteredControllers()
		for _, controller := range controllers {
			routeMap := controller.RouteMap()
			for method, routingFunction := range methodFunctionMap {
				routes := routeMap[method]
				for relativePath, handlerFunc := range routes {
					routingFunction(relativePath, handlerFunc)
				}
			}
		}
	}

	for _, initializer := range routerInitializers {
		initializer.InitializeLate(r)
	}
}
