package extensions

import "github.com/jinzhu/gorm"

type Logic interface {
	GetMulti(db *gorm.DB, queryString string) ([]interface{}, error)
	GetSingle(db *gorm.DB, id string, queryString string) (interface{}, error)
	Create(db *gorm.DB, model interface{}) (interface{}, error)
	Update(db *gorm.DB, id string, model interface{}) (interface{}, error)
	Delete(db *gorm.DB, id string) error
	Patch(db *gorm.DB, id string) (interface{}, error)
	Options(db *gorm.DB) error
}
