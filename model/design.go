package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"net/url"
	"time"
)

// Design is the model class what represents the whole object model store
type Design struct {
	*Base         `json:"base,omitempty"`
	ClayVersion   string                 `json:"clay_version,omitempty"`
	GeneratedDate string                 `json:"generated_date,omitempty"`
	Content       map[string]interface{} `json:"content"`
}

// NewDesign creates a design model instance
func NewDesign() *Design {
	return CreateModel(&Design{}).(*Design)
}

// GetSingle returns all models to store into versioning repositories
func (receiver *Design) GetSingle(db *gorm.DB, _ gin.Params, urlValues url.Values, _ string) (interface{}, error) {
	// Reset previous conditions
	db = db.New()

	programInformation := extension.GetRegisteredProgramInformation()

	design := receiver.NewModelContainer().(*Design)
	design.ClayVersion = programInformation.BuildTime()
	design.GeneratedDate = ""
	design.Content = map[string]interface{}{}

	if _, exists := urlValues["timestamp"]; exists {
		design.GeneratedDate = time.Now().String()
	}

	designAccessors := extension.GetRegisteredDesignAccessorList()
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
func (receiver *Design) Update(db *gorm.DB, _ gin.Params, _ url.Values, input extension.Model) (interface{}, error) {
	design := input.(*Design)

	designAccessors := extension.GetRegisteredDesignAccessorList()
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
func (receiver *Design) Delete(db *gorm.DB, _ gin.Params, _ url.Values) error {
	designAccessors := extension.GetRegisteredDesignAccessorList()
	for _, accessor := range designAccessors {
		if err := accessor.DeleteFromDesign(db); err != nil {
			logging.Logger().Debug(err.Error())
			return err
		}
	}

	return nil
}

func init() {
	extension.RegisterModel(NewDesign(), false)
}
