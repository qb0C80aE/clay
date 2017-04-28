package extensions

import "github.com/gin-gonic/gin"

// Outputter is the interface what handles outputs of the results from the logic to controllers
// * OutputError handles an error output
// * OutputGetSingle corresponds HTTP GET message and handles the output of a single result from logic classes
// * OutputGetMulti corresponds HTTP GET message and handles the output of multiple result from logic classes
// * OutputCreate corresponds HTTP POST message and handles the output of a single result from logic classes
// * OutputUpdate corresponds HTTP PUT message and handles the output of a single result from logic classes
// * OutputDelete corresponds HTTP DELETE message and handles the code result from logic classes
// * OutputPatch corresponds HTTP PATCH message and handles the output of a single result from logic classes
// * OutputOptions corresponds HTTP DELETE message and handles the code result from logic classes, as well as OutputDelete
type Outputter interface {
	OutputError(c *gin.Context, code int, err error)
	OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{})
	OutputGetMulti(c *gin.Context, code int, result interface{}, total int, fields map[string]interface{})
	OutputCreate(c *gin.Context, code int, result interface{})
	OutputUpdate(c *gin.Context, code int, result interface{})
	OutputDelete(c *gin.Context, code int)
	OutputPatch(c *gin.Context, code int, result interface{})
	OutputOptions(c *gin.Context, code int)
}
