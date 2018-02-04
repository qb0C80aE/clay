package logics

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/models"
	"net/url"
	"time"
)

type designLogic struct {
	*BaseLogic
}

func newDesignLogic() *designLogic {
	logic := &designLogic{
		BaseLogic: NewBaseLogic(
			models.SharedDesignModel(),
		),
	}
	return logic
}

// GetSingle returns all models to store into versioning repositories
func (logic *designLogic) GetSingle(db *gorm.DB, _ gin.Params, urlValues url.Values, _ string) (interface{}, error) {
	// Reset previous conditions
	db = db.New()

	programInformation := extensions.RegisteredProgramInformation()

	design := &models.Design{
		ClayVersion:   programInformation.BuildTime(),
		GeneratedDate: "",
		Content:       map[string]interface{}{},
	}

	if _, exists := urlValues["timestamp"]; exists {
		design.GeneratedDate = time.Now().String()
	}

	designAccessors := extensions.RegisteredDesignAccessors()
	for _, accessor := range designAccessors {
		key, value, err := accessor.ExtractFromDesign(db)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		design.Content[key] = value
	}

	return design, nil
}

// Update deletes and updates all models bases on the given data
func (logic *designLogic) Update(db *gorm.DB, _ gin.Params, _ url.Values, data interface{}) (interface{}, error) {
	design := data.(*models.Design)

	designAccessors := extensions.RegisteredDesignAccessors()
	for _, accessor := range designAccessors {
		if err := accessor.DeleteFromDesign(db); err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
	}
	for _, accessor := range designAccessors {
		if err := accessor.LoadToDesign(db, design); err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
	}

	return design, nil
}

// Delete deletes all models
func (logic *designLogic) Delete(db *gorm.DB, _ gin.Params, _ url.Values) error {
	designAccessors := extensions.RegisteredDesignAccessors()
	for _, accessor := range designAccessors {
		if err := accessor.DeleteFromDesign(db); err != nil {
			logging.Logger().Debug(err.Error())
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
	extensions.RegisterLogic(models.SharedDesignModel(), UniqueDesignLogic())
}
