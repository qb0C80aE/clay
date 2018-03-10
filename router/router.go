package router

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/controller"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"os"
)

// Initialize initializes the router
func Initialize(r *gin.Engine) {
	initializerList := extension.GetRegisteredInitializerList()

	for _, initializer := range initializerList {
		initializer.DoBeforeRouterSetup(r)
	}

	r.GET("/", controller.GetAPIEndpoints)

	api := r.Group("")
	{
		methodFunctionMap := map[int]func(string, ...gin.HandlerFunc) gin.IRoutes{
			extension.MethodGet:     api.GET,
			extension.MethodPost:    api.POST,
			extension.MethodPut:     api.PUT,
			extension.MethodDelete:  api.DELETE,
			extension.MethodPatch:   api.PATCH,
			extension.MethodOptions: api.OPTIONS,
		}

		controllerList := extension.GetRegisteredControllerList()
		for _, controller := range controllerList {
			routeMap := controller.GetRouteMap()
			for method, routingFunction := range methodFunctionMap {
				routeList := routeMap[method]
				for pathType, handlerFunc := range routeList {
					switch pathType {
					case extension.URLSingle:
						routingFunction(controller.GetResourceSingleURL(), handlerFunc)
						extension.AssociateControllerWithPath(controller.GetResourceSingleURL(), controller)
					case extension.URLMulti:
						routingFunction(controller.GetResourceMultiURL(), handlerFunc)
						extension.AssociateControllerWithPath(controller.GetResourceMultiURL(), controller)
					default:
						logging.Logger().Criticalf("invalid url type: %d", pathType)
						os.Exit(1)
					}
				}
			}
		}
	}

	for _, initializer := range initializerList {
		initializer.DoAfterRouterSetup(r)
	}
}
