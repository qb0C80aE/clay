package extensions

import "github.com/jinzhu/gorm"

var designAccessors = []DesignAccessor{}

// DesignAccessor is the interface what extracts, deletes, and load the design model from, or to db
// The methods methods disables foreign key while they are processing to avoid insertion order problems
// * ExtractFromDesign extracts the model related to this logic from db
// * DeleteFromDesign deletes the model related to this logic in db
// * LoadToDesign loads the model related to this logic into db
type DesignAccessor interface {
	GetSequenceNumber() int
	ExtractFromDesign(db *gorm.DB) (string, interface{}, error)
	DeleteFromDesign(db *gorm.DB) error
	LoadToDesign(db *gorm.DB, model interface{}) error
}

// RegisterDesignAccessor registers a design accessor used in the design logic
func RegisterDesignAccessor(designAccessor DesignAccessor) {
	designAccessors = append(designAccessors, designAccessor)
}

// RegisteredDesignAccessors returns the registered design accessors
func RegisteredDesignAccessors() []DesignAccessor {
	result := []DesignAccessor{}
	result = append(result, designAccessors...)
	return result
}
