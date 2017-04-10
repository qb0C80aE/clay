package extensions

import "github.com/gin-gonic/gin"

type Outputter interface {
	OutputError(c *gin.Context, code int, err error)
	OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{})
	OutputGetMulti(c *gin.Context, code int, result []interface{}, fields map[string]interface{})
	OutputCreate(c *gin.Context, code int, result interface{})
	OutputUpdate(c *gin.Context, code int, result interface{})
	OutputDelete(c *gin.Context, code int)
	OutputPatch(c *gin.Context, code int, result interface{})
	OutputOptions(c *gin.Context, code int)
}
