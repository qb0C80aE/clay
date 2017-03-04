package logics

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/models"
)

func GetDesign(db *gorm.DB, _ string, queryFields string) (interface{}, error) {

	design := &models.Design{
		Content: map[string]interface{}{},
	}

	designAccessors := extension.GetDesignAccessos()
	for _, accessor := range designAccessors {
		if err := accessor.ExtractFromDesign(db, design.Content); err != nil {
			return nil, err
		}
	}

	return design, nil
}

func UpdateDesign(db *gorm.DB, _ string, data interface{}) (interface{}, error) {
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

func DeleteDesign(db *gorm.DB, _ string) error {
	designAccessors := extension.GetDesignAccessos()
	for _, accessor := range designAccessors {
		if err := accessor.DeleteFromDesign(db); err != nil {
			return err
		}
	}

	return nil
}
