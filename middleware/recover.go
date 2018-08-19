package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"net/http"
	"runtime/debug"
	"strings"
)

// Recover recovers from panic, and return error data to the client
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(c *gin.Context) {
			if e := recover(); e != nil {
				logging.Logger().Debug(e)
				stackLineList := strings.Split(string(debug.Stack()), "\n")
				for _, stackLine := range stackLineList {
					logging.Logger().Debug(strings.Replace(stackLine, "\t", "  ", -1))
				}

				acceptList := strings.Split(c.Request.Header.Get("Accept"), ",")

				err := struct {
					Error string `json:"error" yaml:"error"`
				}{
					Error: fmt.Sprintf("%v", e),
				}

				switch extension.DetermineResponseContentTypeFromAccept(acceptList) {
				case extension.AcceptXYAML, extension.AcceptTextYAML:
					c.YAML(http.StatusInternalServerError, err)
				default:
					c.JSON(http.StatusInternalServerError, err)
				}
			}
		}(c)
		c.Next()
	}
}
