package router

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/controller"
	"github.com/qb0C80aE/clay/extension"
)

// Setup initializes the router
func Setup(engine *gin.Engine) error {
	engine.GET("/", controller.GetAPIEndpoints)

	apiPath := ""
	apiGroup := engine.Group(apiPath)
	extension.RegisterRouterGroup(apiPath, apiGroup)

	initializerList := extension.GetRegisteredInitializerList()
	for _, initializer := range initializerList {
		if err := initializer.DoBeforeRouterSetup(engine); err != nil {
			return nil
		}
	}

	controllerList := extension.GetRegisteredControllerList()
	if err := extension.SetupController(apiPath, controllerList); err != nil {
		return err
	}

	for _, initializer := range initializerList {
		if err := initializer.DoAfterRouterSetup(engine); err != nil {
			return nil
		}
	}

	return nil
}
