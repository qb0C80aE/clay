package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extensions"
	"net/http"
	"sort"
)

// APIEndpoints returns endpoints of Clay API
func APIEndpoints(c *gin.Context) {
	reqScheme := "http"

	if c.Request.TLS != nil {
		reqScheme = "https"
	}

	reqHost := c.Request.Host
	baseURL := fmt.Sprintf("%s://%s", reqScheme, reqHost)
	resources := []string{}

	controllers := extensions.RegisteredControllers()
	for _, controller := range controllers {
		routeMap := controller.RouteMap()
		for method, routes := range routeMap {
			title := fmt.Sprintf("%s_url [%s]", controller.ResourceName(), extensions.LookUpMethodName(method))
			for relativePath := range routes {
				resources = append(resources, fmt.Sprintf("%s %s/%s", title, baseURL, relativePath))
			}
		}
	}
	sort.Strings(resources)
	c.IndentedJSON(http.StatusOK, resources)
}
