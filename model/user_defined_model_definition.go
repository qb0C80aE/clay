package model

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/util/mapstruct"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"sync"
)

var typeNameUserDefinedModelDefinitionMap = map[string]*UserDefinedModelDefinition{}
var typeNameUserDefinedModelDefinitionMapMutex = new(sync.Mutex)

var typeNameTypeMap = map[string]reflect.Type{
	"int":     reflect.TypeOf(int(0)),
	"int8":    reflect.TypeOf(int(0)),
	"int16":   reflect.TypeOf(int(0)),
	"int32":   reflect.TypeOf(int(0)),
	"int64":   reflect.TypeOf(int64(0)),
	"uint":    reflect.TypeOf(uint(0)),
	"uint8":   reflect.TypeOf(uint(0)),
	"uint16":  reflect.TypeOf(uint(0)),
	"uint32":  reflect.TypeOf(uint(0)),
	"uint64":  reflect.TypeOf(uint64(0)),
	"float":   reflect.TypeOf(float32(0)),
	"float64": reflect.TypeOf(float64(0)),
	"bool":    reflect.TypeOf(bool(false)),
	"string":  reflect.TypeOf(string("")),
}

// UserDefinedModelDefinition is the model class what represents raw template
type UserDefinedModelDefinition struct {
	Base
	TypeName                    string                             `json:"type_name" clay:"key_parameter" validate:"required"`
	ResourceName                string                             `json:"resource_name" validate:"required"`
	ToBeMigrated                bool                               `json:"to_be_migrated"`
	IsControllerEnabled         bool                               `json:"is_controller_enabled"`
	SQLBeforeMigration          []string                           `json:"sql_before_migration"`
	SQLAfterMigration           []string                           `json:"sql_after_migration"`
	SQLWhereForDesignExtraction []string                           `json:"sql_where_for_design_extraction"`
	SQLWhereForDesignDeletion   []string                           `json:"sql_where_for_design_deletion"`
	IsManyToManyAssociation     bool                               `json:"is_many_to_many_association"`
	IsDesignAccessDisabled      bool                               `json:"is_design_access_disabled"`
	Fields                      []*UserDefinedModelFieldDefinition `json:"fields" validate:"gt=0"`
}

// NewUserDefinedModelDefinition creates a template raw model instance
func NewUserDefinedModelDefinition() *UserDefinedModelDefinition {
	return &UserDefinedModelDefinition{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *UserDefinedModelDefinition) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *UserDefinedModelDefinition) GetSingle(model extension.Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, queryFields string) (interface{}, error) {
	result, exists := typeNameUserDefinedModelDefinitionMap[parameters.ByName("type_name")]

	if !exists {
		return nil, errors.New("record not found")
	}

	return result, nil
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (receiver *UserDefinedModelDefinition) GetMulti(model extension.Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, queryFields string) (interface{}, error) {
	typeNameUserDefinedModelDefinitionMapMutex.Lock()
	defer typeNameUserDefinedModelDefinitionMapMutex.Unlock()

	keyList := make([]string, 0, len(typeNameUserDefinedModelDefinitionMap))
	for key := range typeNameUserDefinedModelDefinitionMap {
		keyList = append(keyList, key)
	}

	sort.Strings(keyList)

	userDefinedModelDefinitionList := make([]*UserDefinedModelDefinition, len(typeNameUserDefinedModelDefinitionMap))

	for i, key := range keyList {
		userDefinedModelDefinitionList[i] = typeNameUserDefinedModelDefinitionMap[key]
	}

	return userDefinedModelDefinitionList, nil
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (receiver *UserDefinedModelDefinition) Create(model extension.Model, db *gorm.DB, _ gin.Params, _ url.Values, inputContainer interface{}) (interface{}, error) {
	typeNameUserDefinedModelDefinitionMapMutex.Lock()
	defer typeNameUserDefinedModelDefinitionMapMutex.Unlock()

	userDefinedModelDefinition := NewUserDefinedModelDefinition()
	if err := mapstruct.RemapToStruct(inputContainer, userDefinedModelDefinition); err != nil {
		return nil, err
	}

	structFieldNameMapForCheckDuplication := map[string]bool{}
	structFieldJSONKeyMapForCheckDuplication := map[string]bool{}
	structFieldList := []reflect.StructField{}
	for _, field := range userDefinedModelDefinition.Fields {
		if _, exists := structFieldNameMapForCheckDuplication[field.Name]; exists {
			logging.Logger().Debug("the struct %s has the same name fields %s", userDefinedModelDefinition.TypeName, field.Name)
			return nil, fmt.Errorf("the struct %s has the same name fields %s", userDefinedModelDefinition.TypeName, field.Name)
		}
		structFieldNameMapForCheckDuplication[field.Name] = true

		reflectType, exists := typeNameTypeMap[field.TypeName]
		if !exists {
			if field.IsSlice {
				reflectType = extension.NewStructFieldTypeProxy(field.TypeName, reflect.Slice)
			} else {
				reflectType = extension.NewStructFieldTypeProxy(field.TypeName, reflect.Struct)
			}
		}

		structField := reflect.StructField{
			Name: field.Name,
			Tag:  reflect.StructTag(field.Tag),
			Type: reflectType,
		}

		jsonKey := field.Name
		jsonTag, ok := structField.Tag.Lookup("json")
		if ok {
			tagStatementList := strings.Split(jsonTag, ",")
			for _, tagStatement := range tagStatementList {
				switch tagStatement {
				case "omitempty", "-":
					continue
				default:
					jsonKey = tagStatement
					break
				}
			}
		}

		if _, exists := structFieldJSONKeyMapForCheckDuplication[jsonKey]; exists {
			logging.Logger().Debug("the struct %s has the same json key fields %s", userDefinedModelDefinition.TypeName, jsonKey)
			return nil, fmt.Errorf("the struct %s has the same json key fields %s", userDefinedModelDefinition.TypeName, jsonKey)
		}
		structFieldJSONKeyMapForCheckDuplication[jsonKey] = true

		structFieldList = append(structFieldList, structField)
	}

	var newModel extension.Model
	var newInitializer extension.Initializer
	var newDesignAccessor extension.DesignAccessor

	if userDefinedModelDefinition.IsManyToManyAssociation {
		manyToManyLeftResourceName := ""
		manyToManyLeftResourceKeyParameterName := ""
		manyToManyRightResourceName := ""
		manyToManyRightResourceKeyParameterName := ""

		for _, field := range structFieldList {
			jsonKey := field.Name
			jsonTag, ok := field.Tag.Lookup("json")
			if ok {
				tagStatementList := strings.Split(jsonTag, ",")
				for _, tagStatement := range tagStatementList {
					switch tagStatement {
					case "omitempty", "-":
						continue
					default:
						jsonKey = tagStatement
						break
					}
				}
			}

			tag, ok := field.Tag.Lookup("clay")
			if ok {
				tagStatementList := strings.Split(tag, ";")
				for _, tagStatement := range tagStatementList {
					tagKeyValue := strings.Split(tagStatement, "=")
					switch tagKeyValue[0] {
					case "many_to_many_left_resource_name":
						manyToManyLeftResourceName = tagKeyValue[1]
						manyToManyLeftResourceKeyParameterName = jsonKey
					case "many_to_many_right_resource_name":
						manyToManyRightResourceName = tagKeyValue[1]
						manyToManyRightResourceKeyParameterName = jsonKey
					}
				}
			}
		}

		if len(manyToManyLeftResourceName) == 0 {
			logging.Logger().Debug("many_to_many_left_resource_name is not defined")
			return nil, errors.New("many_to_many_left_resource_name is not defined")
		}
		if len(manyToManyRightResourceName) == 0 {
			logging.Logger().Debug("many_to_many_right_resource_name is not defined")
			return nil, errors.New("many_to_many_right_resource_name is not defined")
		}

		newUserDefinedManyToManyAssociationModel := NewUserDefinedManyToManyAssociationModel()
		newUserDefinedManyToManyAssociationModel.leftResourceName = manyToManyLeftResourceName
		newUserDefinedManyToManyAssociationModel.leftResourceKeyName = manyToManyLeftResourceKeyParameterName
		newUserDefinedManyToManyAssociationModel.rightResourceName = manyToManyRightResourceName
		newUserDefinedManyToManyAssociationModel.rightResourceKeyName = manyToManyRightResourceKeyParameterName

		newUserDefinedManyToManyAssociationModel.typeName = userDefinedModelDefinition.TypeName
		newUserDefinedManyToManyAssociationModel.resourceName = userDefinedModelDefinition.ResourceName
		newUserDefinedManyToManyAssociationModel.toBeMigrated = userDefinedModelDefinition.ToBeMigrated
		newUserDefinedManyToManyAssociationModel.isControllerEnabled = userDefinedModelDefinition.IsControllerEnabled
		newUserDefinedManyToManyAssociationModel.isDesignAccessDisabled = userDefinedModelDefinition.IsDesignAccessDisabled
		newUserDefinedManyToManyAssociationModel.sqlBeforeMigration = strings.Join(userDefinedModelDefinition.SQLBeforeMigration, "\n")
		newUserDefinedManyToManyAssociationModel.sqlAfterMigration = strings.Join(userDefinedModelDefinition.SQLAfterMigration, "\n")
		newUserDefinedManyToManyAssociationModel.sqlWhereForDesignExtraction = strings.Join(userDefinedModelDefinition.SQLWhereForDesignExtraction, "\n")
		newUserDefinedManyToManyAssociationModel.sqlWhereForDesignDeletion = strings.Join(userDefinedModelDefinition.SQLWhereForDesignDeletion, "\n")
		newUserDefinedManyToManyAssociationModel.structFieldList = structFieldList

		newModel = newUserDefinedManyToManyAssociationModel
		newInitializer = newUserDefinedManyToManyAssociationModel
		newDesignAccessor = newUserDefinedManyToManyAssociationModel
	} else {
		newUserDefinedModel := NewUserDefinedModel()

		newUserDefinedModel.typeName = userDefinedModelDefinition.TypeName
		newUserDefinedModel.resourceName = userDefinedModelDefinition.ResourceName
		newUserDefinedModel.toBeMigrated = userDefinedModelDefinition.ToBeMigrated
		newUserDefinedModel.isControllerEnabled = userDefinedModelDefinition.IsControllerEnabled
		newUserDefinedModel.isDesignAccessDisabled = userDefinedModelDefinition.IsDesignAccessDisabled
		newUserDefinedModel.sqlBeforeMigration = strings.Join(userDefinedModelDefinition.SQLBeforeMigration, "\n")
		newUserDefinedModel.sqlAfterMigration = strings.Join(userDefinedModelDefinition.SQLAfterMigration, "\n")
		newUserDefinedModel.sqlWhereForDesignExtraction = strings.Join(userDefinedModelDefinition.SQLWhereForDesignExtraction, "\n")
		newUserDefinedModel.sqlWhereForDesignDeletion = strings.Join(userDefinedModelDefinition.SQLWhereForDesignDeletion, "\n")
		newUserDefinedModel.structFieldList = structFieldList

		newModel = newUserDefinedModel
		newInitializer = newUserDefinedModel
		newDesignAccessor = newUserDefinedModel
	}

	extension.RegisterModel(newModel)
	extension.RegisterDesignAccessor(newDesignAccessor)

	initializerList := []extension.Initializer{newInitializer}
	modelList := []extension.Model{newModel}

	if _, err := extension.SetupModel(db, initializerList, modelList); err != nil {
		return nil, err
	}

	typeNameUserDefinedModelDefinitionMap[userDefinedModelDefinition.TypeName] = userDefinedModelDefinition

	return userDefinedModelDefinition, nil
}

// GetTotal returns the count of for multi resource
func (receiver *UserDefinedModelDefinition) GetTotal(_ extension.Model, _ *gorm.DB) (int, error) {
	return len(typeNameUserDefinedModelDefinitionMap), nil
}

func init() {
	extension.RegisterModel(NewUserDefinedModelDefinition())
}
