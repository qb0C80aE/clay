package logics

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/models"
)

type designLogic struct {
	*BaseLogic
}

func newDesignLogic() *designLogic {
	logic := &designLogic{
		BaseLogic: SharedBaseLogic(),
	}
	return logic
}

// GetSingle returns all models to store into versioning repositories
func (logic *designLogic) GetSingle(db *gorm.DB, _ string, _ string) (interface{}, error) {
	db = db.New()

	design := &models.Design{
		Content: map[string]interface{}{},
	}

	designAccessors := extensions.RegisteredDesignAccessors()
	for _, accessor := range designAccessors {
		key, value, err := accessor.ExtractFromDesign(db)
		if err != nil {
			return nil, err
		}
		design.Content[key] = value
	}

	return design, nil
}

// Update deletes and updates all models bases on the given data
func (logic *designLogic) Update(db *gorm.DB, _ string, data interface{}) (interface{}, error) {
	design := data.(*models.Design)

	designAccessors := extensions.RegisteredDesignAccessors()
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

// Delete deletes all models
func (logic *designLogic) Delete(db *gorm.DB, _ string) error {
	designAccessors := extensions.RegisteredDesignAccessors()
	for _, accessor := range designAccessors {
		if err := accessor.DeleteFromDesign(db); err != nil {
			return err
		}
	}

	return nil
}

var uniqueDesignLogic = newDesignLogic()

// UniqueDesignLogic returns the unique design logic instance
func UniqueDesignLogic() extensions.Logic {
	return uniqueDesignLogic
}

func init() {
}
