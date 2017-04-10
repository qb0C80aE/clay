package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extensions"
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

	controllers := extensions.GetControllers()
	for _, controller := range controllers {
		routeMap := controller.RouteMap()
		for method, routes := range routeMap {
			title := fmt.Sprintf("%s_url [%s]", controller.ResourceName(), extensions.GetMethodName(method))
			for relativePath := range routes {
				resources[title] = fmt.Sprintf("%s/%s", baseURL, relativePath)
			}
		}
	}
	c.IndentedJSON(http.StatusOK, resources)
}

func Initialize(r *gin.Engine) {
	routerInitializers := extensions.GetRouterInitializers()

	for _, initializer := range routerInitializers {
		initializer.InitializeEarly(r)
	}

	r.GET("/", getAPIEndpoints)

	api := r.Group("/v1")
	{
		var methodFunctionMap map[int]func(string, ...gin.HandlerFunc) gin.IRoutes = map[int]func(string, ...gin.HandlerFunc) gin.IRoutes{
			extensions.MethodGet:     api.GET,
			extensions.MethodPost:    api.POST,
			extensions.MethodPut:     api.PUT,
			extensions.MethodDelete:  api.DELETE,
			extensions.MethodPatch:   api.PATCH,
			extensions.MethodOptions: api.OPTIONS,
		}

		controllers := extensions.GetControllers()
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
