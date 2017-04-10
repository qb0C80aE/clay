package extensions

import "github.com/jinzhu/gorm"

var designAccessors = []DesignAccessor{}

type DesignAccessor interface {
	ExtractFromDesign(db *gorm.DB) (string, interface{}, error)
	DeleteFromDesign(db *gorm.DB) error
	LoadToDesign(db *gorm.DB, model interface{}) error
}

func RegisterDesignAccessor(designAccessor DesignAccessor) {
	designAccessors = append(designAccessors, designAccessor)
}

func GetDesignAccessos() []DesignAccessor {
	result := []DesignAccessor{}
	result = append(result, designAccessors...)
	return result
}
