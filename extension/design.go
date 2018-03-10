package extension

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/logging"
	"reflect"
)

var designAccessorList = []DesignAccessor{}

// DesignAccessor is the interface what extracts, deletes, and load the design model from, or to db
// The methods methods disables foreign key while they are processing to avoid insertion order problems
// * ExtractFromDesign extracts the model related to this logic from db
// * DeleteFromDesign deletes the model related to this logic in db
// * LoadToDesign loads the model related to this logic into db
type DesignAccessor interface {
	ExtractFromDesign(db *gorm.DB) (string, interface{}, error)
	DeleteFromDesign(db *gorm.DB) error
	LoadToDesign(db *gorm.DB, model interface{}) error
}

// RegisterDesignAccessor registers a design accessor used in the design logic
func RegisterDesignAccessor(designAccessor DesignAccessor) {
	if reflect.ValueOf(designAccessor).Elem().FieldByName("Base").IsNil() {
		logging.Logger().Criticalf("the designAccessor is a container which does not have *Base")
		panic("the designAccessor is a container which does not have *Base")
	}

	designAccessorList = append(designAccessorList, designAccessor)
}

// GetRegisteredDesignAccessorList returns the registered design accessors
func GetRegisteredDesignAccessorList() []DesignAccessor {
	result := []DesignAccessor{}
	result = append(result, designAccessorList...)
	return result
}
