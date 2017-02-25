package router

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
)

func Initialize(r *gin.Engine) {

	r.GET("/", extension.APIEndpoints)

	api := r.Group("/v1")
	{
		routes := extension.GetRoutes(extension.MethodGet)
		for relativePath, handlerFunc := range routes {
			api.GET(relativePath, handlerFunc)
		}
		routes = extension.GetRoutes(extension.MethodPost)
		for relativePath, handlerFunc := range routes {
			api.POST(relativePath, handlerFunc)
		}
		routes = extension.GetRoutes(extension.MethodPut)
		for relativePath, handlerFunc := range routes {
			api.PUT(relativePath, handlerFunc)
		}
		routes = extension.GetRoutes(extension.MethodDelete)
		for relativePath, handlerFunc := range routes {
			api.DELETE(relativePath, handlerFunc)
		}
		routes = extension.GetRoutes(extension.MethodPatch)
		for relativePath, handlerFunc := range routes {
			api.PATCH(relativePath, handlerFunc)
		}
	}

	routePostInitializers := extension.GetRoutePostInitializers()
	for _, initializer := range routePostInitializers {
		initializer(r)
	}

}
