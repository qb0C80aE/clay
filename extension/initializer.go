package extension

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var initializerList = []Initializer{}

// Initializer is the interface what adds router processes specifically before and after router registration
// * DoBeforeDBMigration execute initialization process before DB migration
// * DoAfterDBMigration execute initialization process after DB migration
// * DoBeforeRouterSetup execute initialization process before Router initialization
// * DoAfterRouterSetup execute initialization process after Router initialization
type Initializer interface {
	DoBeforeDBMigration(db *gorm.DB) error
	DoAfterDBMigration(tx *gorm.DB) error
	DoBeforeRouterSetup(r *gin.Engine) error
	DoAfterRouterSetup(r *gin.Engine) error
}

// RegisterInitializer registers an initializer used in the router logic
func RegisterInitializer(initializer Initializer) {
	initializerList = append(initializerList, initializer)
}

// GetRegisteredInitializerList returns the registered router initializers
func GetRegisteredInitializerList() []Initializer {
	result := []Initializer{}
	result = append(result, initializerList...)
	return result
}
