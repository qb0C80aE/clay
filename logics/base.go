package logics

import (
	"github.com/jinzhu/gorm"
)

type BaseLogic struct {
}

func newBaseLogic() *BaseLogic {
	return &BaseLogic{}
}

func (logic *BaseLogic) GetSingle(db *gorm.DB, id string, queryFields string) (interface{}, error) {
	return nil, nil
}

func (logic *BaseLogic) GetMulti(db *gorm.DB, queryFields string) ([]interface{}, error) {
	return nil, nil
}

func (logic *BaseLogic) Create(db *gorm.DB, data interface{}) (interface{}, error) {
	return nil, nil
}

func (logic *BaseLogic) Update(db *gorm.DB, id string, data interface{}) (interface{}, error) {
	return nil, nil
}

func (logic *BaseLogic) Delete(db *gorm.DB, id string) error {
	return nil
}

func (logic *BaseLogic) Patch(db *gorm.DB, id string) (interface{}, error) {
	return nil, nil
}

func (logic *BaseLogic) Options(db *gorm.DB) error {
	return nil
}

var baseLogic = newBaseLogic()

func SharedBaseLogic() *BaseLogic {
	return baseLogic
}

func init() {
}
