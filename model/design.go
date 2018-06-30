package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	mapstructutilpkg "github.com/qb0C80aE/clay/util/mapstruct"
	"net/url"
	"reflect"
	"sort"
	"time"
)

type sqliteMaster struct {
	Name     string `json:"name" yaml:"name"`
	SQL      string `json:"sql" yaml:"sql"`
	Type     string `json:"type" yaml:"type"`
	Rootpage int    `json:"rootpage" yaml:"rootpage"`
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

func sortDesignAccessorList(designAccessorList []extension.DesignAccessor, tableMap map[string]*sqliteMaster, rootPageOrderAsc bool) error {
	var err error
	defer func() {
		recoverResult := recover()
		if recoverResult != nil {
			err = fmt.Errorf("%v", recoverResult)
		}
	}()

	sort.Slice(designAccessorList, func(i, j int) bool {
		left, err := designAccessorList[i].GetResourceName(designAccessorList[i])
		if err != nil {
			logging.Logger().Debug(err.Error())
			panic(err)
		}
		right, err := designAccessorList[j].GetResourceName(designAccessorList[j])
		if err != nil {
			logging.Logger().Debug(err.Error())
			panic(err)
		}
		if _, exists := tableMap[left]; !exists {
			return rootPageOrderAsc
		} else if _, exists := tableMap[right]; !exists {
			return !rootPageOrderAsc
		}
		if rootPageOrderAsc {
			return tableMap[left].Rootpage < tableMap[right].Rootpage
		}
		return tableMap[left].Rootpage > tableMap[right].Rootpage
	})

	return err
}

// Update deletes and updates all models bases on the given data
func (receiver *Design) Update(_ extension.Model, db *gorm.DB, _ gin.Params, _ url.Values, inputContainer interface{}) (interface{}, error) {
	triggerList := []*sqliteMaster{}
	tableList := []*sqliteMaster{}

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

	if err := db.Table("sqlite_master").Where("type = ?", "table").Order("rootpage asc").Find(&tableList).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	tableMap := make(map[string]*sqliteMaster, len(tableList))
	for _, table := range tableList {
		tableMap[table.Name] = table
	}

	design := NewDesign()
	if err := mapstructutilpkg.GetUtility().MapToStruct(inputContainer, design); err != nil {
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

	designAccessorList := extension.GetRegisteredDesignAccessorList()

	// delete by following the rootpage index
	sortDesignAccessorList(designAccessorList, tableMap, false)
	for _, accessor := range designAccessorList {
		if err := accessor.DeleteFromDesign(accessor, db); err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
	}

	// insert by following the rootpage index
	sortDesignAccessorList(designAccessorList, tableMap, true)
	for _, accessor := range designAccessorList {
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
	triggerList := []*sqliteMaster{}
	tableList := []*sqliteMaster{}

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

	if err := db.Table("sqlite_master").Where("type = ?", "table").Order("rootpage asc").Find(&tableList).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	tableMap := make(map[string]*sqliteMaster, len(tableList))
	for _, table := range tableList {
		tableMap[table.Name] = table
	}

	designAccessorList := extension.GetRegisteredDesignAccessorList()

	// delete by following the rootpage index
	sortDesignAccessorList(designAccessorList, tableMap, false)
	for _, accessor := range designAccessorList {
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
