package logics

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/models"
)

type DesignLogic struct {
}

func NewDesignLogic() *DesignLogic {
	return &DesignLogic{}
}

func (_ *DesignLogic) GetSingle(db *gorm.DB, _ string, _ string) (interface{}, error) {

	design := &models.Design{
		Content: map[string]interface{}{},
	}

	designAccessors := extension.GetDesignAccessos()
	for _, accessor := range designAccessors {
		key, value, err := accessor.ExtractFromDesign(db)
		if err != nil {
			return nil, err
		}
		design.Content[key] = value
	}

	return design, nil
}

func (_ *DesignLogic) GetMulti(_ *gorm.DB, _ string) ([]interface{}, error) {
	return nil, nil
}

func (this *DesignLogic) Create(db *gorm.DB, data interface{}) (interface{}, error) {
	return nil, nil
}

func (_ *DesignLogic) Update(db *gorm.DB, _ string, data interface{}) (interface{}, error) {
	design := data.(*models.Design)

	designAccessors := extension.GetDesignAccessos()
	for _, accessor := range designAccessors {
		if err := accessor.DeleteFromDesign(db); err != nil {
			return nil, err
		}
	}
	for _, accessor := range designAccessors {
		if err := accessor.LoadToDesign(db, design); err != nil {
			return nil, err
		}
	}

	return design, nil
}

func (_ *DesignLogic) Delete(db *gorm.DB, _ string) error {
	designAccessors := extension.GetDesignAccessos()
	for _, accessor := range designAccessors {
		if err := accessor.DeleteFromDesign(db); err != nil {
			return err
		}
	}

	return nil
}

func (_ *DesignLogic) Patch(_ *gorm.DB, _ string, _ string) (interface{}, error) {
	return nil, nil
}

func (_ *DesignLogic) Options(db *gorm.DB) error {
	return nil
}

var DesignLogicInstance = &DesignLogic{}

func init() {
}
