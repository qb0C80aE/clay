package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/util/mapstruct"
	"net/url"
	"reflect"
	"time"
)

type master struct {
	Name string
	SQL  string
	Type string
}

// Design is the model class what represents the whole object model store
type Design struct {
	Base
	ClayVersion   string                 `json:"clay_version,omitempty" yaml:"clay_version,omitempty"`
	GeneratedDate string                 `json:"generated_date,omitempty" yaml:"generated_date,omitempty"`
	Content       map[string]interface{} `json:"content" yaml:"content"`
}

// NewDesign creates a design model instance
func NewDesign() *Design {
	return &Design{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *Design) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle returns all models to store into versioning repositories
func (receiver *Design) GetSingle(_ extension.Model, db *gorm.DB, _ gin.Params, urlValues url.Values, _ string) (interface{}, error) {
	// Reset previous conditions
	db = db.New()

	programInformation := extension.GetRegisteredProgramInformation()

	design := NewDesign()
	design.ClayVersion = programInformation.GetBuildTime()
	design.GeneratedDate = ""
	design.Content = map[string]interface{}{}

	if _, exists := urlValues["timestamp"]; exists {
		design.GeneratedDate = time.Now().String()
	}

	designAccessors := extension.GetRegisteredDesignAccessorList()

	for _, accessor := range designAccessors {
		key, value, err := accessor.ExtractFromDesign(accessor, db)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}

		if len(key) > 0 {
			design.Content[key] = value
		}
	}

	return design, nil
}

// Update deletes and updates all models bases on the given data
func (receiver *Design) Update(_ extension.Model, db *gorm.DB, _ gin.Params, _ url.Values, inputContainer interface{}) (interface{}, error) {
	triggerList := []*master{}

	if err := db.Table("sqlite_master").Where("type = ?", "trigger").Find(&triggerList).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	for _, trigger := range triggerList {
		if err := db.Exec(fmt.Sprintf("drop trigger %s", trigger.Name)).Error; err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
	}

	design := NewDesign()
	if err := mapstruct.RemapToStruct(inputContainer, design); err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	mapKeys := reflect.ValueOf(inputContainer).Elem().FieldByName("Content").MapKeys()
	for _, mapKey := range mapKeys {
		if _, err := extension.GetAssociatedModelWithResourceName(mapKey.Interface().(string)); err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
	}

	designAccessors := extension.GetRegisteredDesignAccessorList()
	for _, accessor := range designAccessors {
		if err := accessor.DeleteFromDesign(accessor, db); err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
	}
	for _, accessor := range designAccessors {
		if err := accessor.LoadToDesign(accessor, db, design); err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
	}

	for _, trigger := range triggerList {
		if err := db.Exec(trigger.SQL).Error; err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
	}

	return design, nil
}

// Delete deletes all models
func (receiver *Design) Delete(_ extension.Model, db *gorm.DB, _ gin.Params, _ url.Values) error {
	triggerList := []*master{}

	if err := db.Table("sqlite_master").Where("type = ?", "trigger").Find(&triggerList).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	for _, trigger := range triggerList {
		if err := db.Exec(fmt.Sprintf("drop trigger %s", trigger.Name)).Error; err != nil {
			logging.Logger().Debug(err.Error())
			return err
		}
	}

	designAccessors := extension.GetRegisteredDesignAccessorList()
	for _, accessor := range designAccessors {
		if err := accessor.DeleteFromDesign(accessor, db); err != nil {
			logging.Logger().Debug(err.Error())
			return err
		}
	}

	for _, trigger := range triggerList {
		if err := db.Exec(trigger.SQL).Error; err != nil {
			logging.Logger().Debug(err.Error())
			return err
		}
	}

	return nil
}

func init() {
	extension.RegisterModel(NewDesign())
}
