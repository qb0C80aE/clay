package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"net/http"
)

func getAPIEndpoints(c *gin.Context) {
	reqScheme := "http"

	if c.Request.TLS != nil {
		reqScheme = "https"
	}

	reqHost := c.Request.Host
	baseURL := fmt.Sprintf("%s://%s/%s", reqScheme, reqHost, "v1")
	resources := map[string]string{}

	controllers := extension.GetControllers()
	for _, controller := range controllers {
		routeMap := controller.GetRouteMap()
		for method, routes := range routeMap {
			title := fmt.Sprintf("%s_url [%s]", controller.GetResourceName(), extension.GetMethodName(method))
			for relativePath := range routes {
				resources[title] = fmt.Sprintf("%s/%s", baseURL, relativePath)
			}
		}
	}
	c.IndentedJSON(http.StatusOK, resources)
}

func Initialize(r *gin.Engine) {
	routerInitializers := extension.GetRouterInitializers()

	for _, initializer := range routerInitializers {
		initializer.InitializeEarly(r)
	}

	r.GET("/", getAPIEndpoints)

	api := r.Group("/v1")
	{
		var methodFunctionMap map[int]func(string, ...gin.HandlerFunc) gin.IRoutes = map[int]func(string, ...gin.HandlerFunc) gin.IRoutes{
			extension.MethodGet:     api.GET,
			extension.MethodPost:    api.POST,
			extension.MethodPut:     api.PUT,
			extension.MethodDelete:  api.DELETE,
			extension.MethodPatch:   api.PATCH,
			extension.MethodOptions: api.OPTIONS,
		}

		controllers := extension.GetControllers()
		for _, controller := range controllers {
			routeMap := controller.GetRouteMap()
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
