package extension

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/logging"
	"github.com/serenize/snaker"
	"net/url"
	"reflect"
	"strings"
)

var modelList = []Model{}
var typeNameModelMap = map[string]Model{}
var resourceNameModelMap = map[string]Model{}
var typeNameResourceNameMap = map[string]string{}
var resourceNameTypeNameMap = map[string]string{}
var typeNameDefaultModelKeyMap = map[string]ModelKey{}

var typeNameStructFieldListMap = map[string][]reflect.StructField{}
var typeNameJSONKeyStructFieldMapMap = map[string]map[string]reflect.StructField{}

type structTree map[string]structTree

var structFieldTypeProxyType = reflect.TypeOf(&StructFieldTypeProxy{})

// StructFieldTypeProxy is a struct used to represent a type of struct fields in container
// It proxies those types by returning custom information about types
// This cannot be used with reflect.New, just used to know its name or kind
type StructFieldTypeProxy struct {
	reflect.Type
	name string
	kind reflect.Kind
}

// Name overrides reflect.Type.Name()
func (receiver *StructFieldTypeProxy) Name() string {
	return receiver.name
}

// Kind overrides reflect.Type.Kind()
func (receiver *StructFieldTypeProxy) Kind() reflect.Kind {
	return receiver.kind
}

// NewStructFieldTypeProxy creates a new StructFieldTypeProxy instance
func NewStructFieldTypeProxy(name string, kind reflect.Kind) *StructFieldTypeProxy {
	structFieldTypePlaceHolder := &StructFieldTypeProxy{
		Type: structFieldTypeProxyType,
		name: name,
		kind: kind,
	}
	return structFieldTypePlaceHolder
}

// Model is the interface that all model struct inherit
// It has methods about container instance creation, and container CRUD operations
// * GetTypeName returns its struct type name
// * GetResourceName returns its resource/table name in URL/DB
// * GenerateTableName generates its resource/table name in URL/DB
// * GetStructFields returns its struct fields used to create containers
// * GetModelKey returns its key fields
// * GetContainerForMigration returns its container for migration
// * GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
// * GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
// * Create corresponds HTTP POST message and handles a request for multi resource to create a new information
// * Update corresponds HTTP PUT message and handles a request for a single resource to update the specific information
// * Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
// * Patch corresponds HTTP PATCH message and handles a request for a single resource to update partially the specific information
// * GetOptions corresponds HTTP OPTIONS message and handles a request for multi resources to retrieve its supported options
// * GetCount returns the record count under current db conditions
type Model interface {
	GetTypeName(model Model) string
	GetResourceName(model Model) (string, error)
	GenerateTableName(model Model, db *gorm.DB) string
	GetStructFields(model Model) []reflect.StructField
	GetModelKey(model Model, keyParameterSpecifier string) (ModelKey, error)
	GetContainerForMigration() (interface{}, error)
	GetSingle(model Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, queryString string) (interface{}, error)
	GetMulti(model Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, queryString string) (interface{}, error)
	Create(model Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, inputContainer interface{}) (interface{}, error)
	Update(model Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, inputContainer interface{}) (interface{}, error)
	Delete(model Model, db *gorm.DB, parameters gin.Params, urlValues url.Values) error
	Patch(model Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, inputContainer interface{}) (interface{}, error)
	GetOptions(model Model, db *gorm.DB, parameters gin.Params, urlValues url.Values) error
	GetCount(model Model, db *gorm.DB) (int, error)
}

// ModelKey is the type that defines the key parameter of types. It's used in various functions like "to_be_deleted" (delete specific children in update) logic or any other key replacement logic.
type ModelKey struct {
	KeyParameter string
	KeyField     string
}

// InspectActualElementType returns the actual type of given object
// It strips pointer. If the object is array or slice, this returns the type of element of the object
func InspectActualElementType(object interface{}) reflect.Type {
	objectType := reflect.TypeOf(object)
	for (objectType.Kind() == reflect.Ptr) || (objectType.Kind() == reflect.Interface) {
		objectType = objectType.Elem()
	}

	switch objectType.Kind() {
	case reflect.Array, reflect.Slice:
		temporarySlice := reflect.MakeSlice(objectType, 1, 1)
		objectType = reflect.TypeOf(temporarySlice.Index(0).Interface())

		for (objectType.Kind() == reflect.Ptr) || (objectType.Kind() == reflect.Interface) {
			objectType = objectType.Elem()
		}
	}

	return objectType
}

// GetActualType returns the type of given object stripping pointer and interface
func GetActualType(object interface{}) reflect.Type {
	objectType := reflect.TypeOf(object)
	for (objectType.Kind() == reflect.Ptr) || (objectType.Kind() == reflect.Interface) {
		objectType = objectType.Elem()
	}
	return objectType
}

// GetActualValue returns the value of given object stripping pointer and interface
func GetActualValue(object interface{}) reflect.Value {
	objectValue := reflect.ValueOf(object)
	for (objectValue.Kind() == reflect.Ptr) || (objectValue.Kind() == reflect.Interface) {
		objectValue = objectValue.Elem()
	}
	return objectValue
}

// GetModelKey returns the model key of given model
// If key_parameter is not empty, this returns overwritten one
func GetModelKey(model Model, keyParameterSpecifier string) (ModelKey, error) {
	modelKey, err := GetRegisteredDefaultModelKey(model)
	if err != nil {
		return ModelKey{}, err
	}

	if len(keyParameterSpecifier) > 0 {
		container, err := CreateContainerByTypeName(model.GetTypeName(model))
		if err != nil {
			return ModelKey{}, err
		}

		containerType := GetActualType(container)

		jsonKey := ""
		fieldName := ""
		for i := 0; i < containerType.NumField(); i++ {
			field := containerType.Field(i)
			jsonTag, ok := field.Tag.Lookup("json")
			if ok {
				tagStatementList := strings.Split(jsonTag, ",")
				for _, tagStatement := range tagStatementList {
					switch tagStatement {
					case "omitempty", "-":
						continue
					default:
						if keyParameterSpecifier == tagStatement {
							jsonKey = tagStatement
							fieldName = field.Name
							break
						}
					}
				}
			}
		}

		if len(fieldName) == 0 {
			return ModelKey{}, fmt.Errorf("key_parameter %s does not exist in %s", keyParameterSpecifier, containerType.String())
		}

		modelKey = ModelKey{
			KeyParameter: jsonKey,
			KeyField:     fieldName,
		}
	}

	return modelKey, nil
}

// RegisterModel registers a model to migrate automatically, and to generate new instances in processing requests
func RegisterModel(model Model) {
	modelList = append(modelList, model)
	modelTypeName := model.GetTypeName(model)
	modelFieldList := model.GetStructFields(model)
	typeNameJSONKeyStructFieldMap := map[string]reflect.StructField{}
	typeNameStructFieldList := []reflect.StructField{}

	for _, field := range modelFieldList {
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
		typeNameJSONKeyStructFieldMap[jsonKey] = field

		typeNameStructFieldList = append(typeNameStructFieldList, field)

		tag, ok := field.Tag.Lookup("clay")
		if ok {
			tagStatementList := strings.Split(tag, ";")
			for _, tagStatement := range tagStatementList {
				switch tagStatement {
				case "key_parameter":
					typeNameDefaultModelKeyMap[modelTypeName] = ModelKey{
						KeyParameter: jsonKey,
						KeyField:     field.Name,
					}
					break
				}
			}
		}
	}

	if _, exists := typeNameDefaultModelKeyMap[modelTypeName]; !exists {
		typeNameDefaultModelKeyMap[modelTypeName] = ModelKey{
			KeyParameter: "id",
			KeyField:     "ID",
		}
	}

	typeNameStructFieldListMap[modelTypeName] = typeNameStructFieldList
	typeNameJSONKeyStructFieldMapMap[modelTypeName] = typeNameJSONKeyStructFieldMap
}

// GetRegisteredModelList returns the registered models
func GetRegisteredModelList() []Model {
	result := make([]Model, len(modelList))
	for i, model := range modelList {
		result[i] = model
	}
	return result
}

// GetRegisteredDefaultModelKey returns the registered default model key
func GetRegisteredDefaultModelKey(model Model) (ModelKey, error) {
	modelTypeName := model.GetTypeName(model)
	result, exists := typeNameDefaultModelKeyMap[modelTypeName]
	if !exists {
		logging.Logger().Debugf("the model key of '%s' has not been registered yet", modelTypeName)
		return ModelKey{}, fmt.Errorf("the model key of '%s' has not been registered yet", modelTypeName)
	}
	return result, nil
}

// AssociateResourceNameWithModel registers a name of given model
func AssociateResourceNameWithModel(resourceName string, model Model) error {
	modelTypeName := model.GetTypeName(model)

	if registeredModelTypeName, exists := typeNameResourceNameMap[modelTypeName]; exists {
		logging.Logger().Debugf("the resource name '%s' has already associated with %s", resourceName, registeredModelTypeName)
		return fmt.Errorf("the resource name '%s' has already associated with %s", resourceName, registeredModelTypeName)
	}

	typeNameResourceNameMap[modelTypeName] = resourceName
	resourceNameTypeNameMap[resourceName] = modelTypeName
	typeNameModelMap[modelTypeName] = model
	resourceNameModelMap[resourceName] = model

	return nil
}

// GetAssociatedResourceNameWithModel returns the registered resource name of given model
func GetAssociatedResourceNameWithModel(model Model) (string, error) {
	modelTypeName := model.GetTypeName(model)

	result, exists := typeNameResourceNameMap[modelTypeName]
	if !exists {
		logging.Logger().Debugf("the resource name '%s' has not been registered yet", modelTypeName)
		return "", fmt.Errorf("the resource name '%s' has not been registered yet", modelTypeName)
	}
	return result, nil
}

// GetAssociatedModelWithResourceName returns the registered model associated with given resource name
func GetAssociatedModelWithResourceName(resourceName string) (Model, error) {
	result, exists := resourceNameModelMap[resourceName]
	if !exists {
		logging.Logger().Debugf("the resource name '%s' has not been associated yet", resourceName)
		return nil, fmt.Errorf("the resource name '%s' has not been associated yet", resourceName)
	}
	return result, nil
}

// GetAssociatedModelWithTypeName returns the registered model associated with given type name
func GetAssociatedModelWithTypeName(typeName string) (Model, error) {
	result, exists := typeNameModelMap[typeName]
	if !exists {
		logging.Logger().Debugf("the type name '%s' has not been associated yet", typeName)
		return nil, fmt.Errorf("the type name '%s' has not been associated yet", typeName)
	}
	return result, nil
}

func createPreloadTree(preloadQuery string) structTree {
	if len(preloadQuery) == 0 {
		return structTree{}
	}

	tree := structTree{}
	fullFieldNameList := strings.Split(preloadQuery, ",")
	for _, fullFieldName := range fullFieldNameList {
		subTree := tree
		fieldNameList := strings.Split(fullFieldName, ".")
		for _, fieldName := range fieldNameList {
			cameledFieldName := snaker.SnakeToCamel(fieldName)
			temp, exists := subTree[cameledFieldName]
			if exists {
				subTree = temp
			} else {
				subTree[cameledFieldName] = structTree{}
				subTree = subTree[cameledFieldName]
			}
		}
	}

	return tree
}

func pruneInputTree(input interface{}, tree structTree) error {
	inputValue := reflect.ValueOf(input)
	switch inputValue.Type().Kind() {
	case reflect.Map:
		mapKeyList := inputValue.MapKeys()
		for _, mapKey := range mapKeyList {
			mapValue := inputValue.MapIndex(mapKey)

			for (mapValue.Kind() == reflect.Ptr) || (mapValue.Kind() == reflect.Interface) {
				mapValue = mapValue.Elem()
			}

			switch mapValue.Kind() {
			case reflect.Map:
				subTree, exists := tree[mapKey.Interface().(string)]
				if !exists {
					subTree = structTree{}
					tree[mapKey.Interface().(string)] = subTree
				}

				if err := pruneInputTree(mapValue.Interface(), subTree); err != nil {
					logging.Logger().Debug(err.Error())
					return err
				}
			case reflect.Slice, reflect.Array:
				subTree, exists := tree[mapKey.Interface().(string)]
				if !exists {
					subTree = structTree{}
					tree[mapKey.Interface().(string)] = subTree
				}

				for i := 0; i < mapValue.Len(); i++ {
					inputValueElementValue := mapValue.Index(i)
					if err := pruneInputTree(inputValueElementValue.Interface(), subTree); err != nil {
						logging.Logger().Debug(err.Error())
						return err
					}
				}
			default:
				// just ignore primitives
			}
		}
	default:
		// just ignore primitive slice/array elements, top-level slice/array/primitives
	}

	return nil
}

func createOutputContainerByTypeName(typeName string, preloadTree structTree) (reflect.Type, error) {
	structFieldList, exists := typeNameStructFieldListMap[typeName]
	if !exists {
		logging.Logger().Debugf("the type name '%s' has not been registered yet", typeName)
		return nil, fmt.Errorf("the type name '%s' has not been registered yet", typeName)
	}

	for key := range preloadTree {
		found := false
		for _, structField := range structFieldList {
			if structField.Name == key {
				found = true
				break
			}
		}
		if !found {
			logging.Logger().Debugf("could not preload, struct '%s' does not have the field named %s", typeName, key)
			return nil, fmt.Errorf("could not preload, struct '%s' does not have the field named %s", typeName, key)
		}
	}

	newStructFieldList := make([]reflect.StructField, len(structFieldList), len(structFieldList))

	for i, structField := range structFieldList {
		switch structField.Type.Kind() {
		case reflect.Struct, reflect.Slice, reflect.Array:
			preloadSubTree, exists := preloadTree[structField.Name]
			if exists {
				childTypeName := structField.Type.Name()
				childStructType, err := createOutputContainerByTypeName(childTypeName, preloadSubTree)
				if err != nil {
					logging.Logger().Debug(err.Error())
					return nil, err
				}

				switch structField.Type.Kind() {
				case reflect.Struct:
					childStructType = reflect.PtrTo(childStructType)
				case reflect.Slice, reflect.Array:
					childStructType = reflect.SliceOf(reflect.PtrTo(childStructType))
				}

				newStructField := structField
				newStructField.Type = childStructType

				newStructFieldList[i] = newStructField
			} else {
				newStructFieldList[i] = structField
			}
		default:
			newStructFieldList[i] = structField
		}
	}

	tableNameStructField := reflect.StructField{
		Name: "StructMetaInformation",
		Type: structFieldTypeProxyType,
		Tag:  reflect.StructTag(fmt.Sprintf("json:\"-,omitempty\" yaml:\"-,omitempty\" sql:\"-\" clay:\"type_name=%s\"", typeName)),
	}
	newStructFieldList = append(newStructFieldList, tableNameStructField)

	return reflect.StructOf(newStructFieldList), nil
}

func createInputContainerByTypeName(typeName string, inputMap interface{}) (reflect.Type, error) {
	inputMapValue := reflect.ValueOf(inputMap)

	switch inputMapValue.Type().Kind() {
	case reflect.Slice, reflect.Array:
		// zero-length slice/array are not allowed
		if inputMapValue.Len() == 0 {
			logging.Logger().Debug("zero-length slice/array are not allowed")
			return nil, errors.New("zero-length slice/array are not allowed")
		}
	case reflect.Map:
	default:
		// top-level non map/slice/array are not allowed
		logging.Logger().Debug("top-level non map/slice/array are not allowed")
		return nil, errors.New("top-level non map/slice/array are not allowed")
	}

	jsonKeyStructFieldMap, exists := typeNameJSONKeyStructFieldMapMap[typeName]
	if !exists {
		logging.Logger().Debugf("the type name '%s' has not been registered yet", typeName)
		return nil, fmt.Errorf("the type name '%s' has not been registered yet", typeName)
	}

	newStructFieldList := make([]reflect.StructField, 0, len(jsonKeyStructFieldMap))

	for jsonKey, jsonKeyStructField := range jsonKeyStructFieldMap {
		inputMapKeyValueValue := inputMapValue.MapIndex(reflect.ValueOf(jsonKey))
		if inputMapKeyValueValue.IsValid() {
			switch jsonKeyStructField.Type.Kind() {
			case reflect.Struct:
				childTypeName := jsonKeyStructField.Type.Name()
				childStructType, err := createInputContainerByTypeName(childTypeName, inputMapKeyValueValue.Interface())
				if err != nil {
					logging.Logger().Debug(err.Error())
					return nil, err
				}

				childStructType = reflect.PtrTo(childStructType)

				newStructField := jsonKeyStructField
				newStructField.Type = childStructType

				newStructFieldList = append(newStructFieldList, newStructField)
			case reflect.Slice, reflect.Array:
				childTypeName := jsonKeyStructField.Type.Name()
				//if len(jsonKeyStructField.Type.Name()) > 0 {
				if _, isStructFieldTypeProxy := jsonKeyStructField.Type.(*StructFieldTypeProxy); isStructFieldTypeProxy {
					// if the type of element is a registered struct type
					childStructType, err := createInputContainerByTypeName(childTypeName, inputMapKeyValueValue.Interface())
					if err != nil {
						logging.Logger().Debug(err.Error())
						return nil, err
					}

					childStructType = reflect.SliceOf(reflect.PtrTo(childStructType))

					newStructField := jsonKeyStructField
					newStructField.Type = childStructType

					newStructFieldList = append(newStructFieldList, newStructField)
				} else {
					// if the type of element is a primitive type
					newStructFieldList = append(newStructFieldList, jsonKeyStructField)
				}
			default: //primitive
				newStructFieldList = append(newStructFieldList, jsonKeyStructField)
			}
		} else {
			newStructFieldList = append(newStructFieldList, jsonKeyStructField)
		}
	}

	tableNameStructField := reflect.StructField{
		Name: "StructMetaInformation",
		Type: structFieldTypeProxyType,
		Tag:  reflect.StructTag(fmt.Sprintf("json:\"-,omitempty\" sql:\"-\" clay:\"type_name=%s\"", typeName)),
	}
	newStructFieldList = append(newStructFieldList, tableNameStructField)

	return reflect.StructOf(newStructFieldList), nil
}

// CreateInputContainerByTypeName creates a container related to given typeName
// The struct fields of target container are determined by inputMap
func CreateInputContainerByTypeName(typeName string, inputMap interface{}) (interface{}, error) {
	prunedInputTree := structTree{}
	if err := pruneInputTree(inputMap, prunedInputTree); err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	_, exists := typeNameJSONKeyStructFieldMapMap[typeName]
	if !exists {
		logging.Logger().Debugf("the type name '%s' has not been registered yet", typeName)
		return nil, fmt.Errorf("the type name '%s' has not been registered yet", typeName)
	}

	containerType, err := createInputContainerByTypeName(typeName, prunedInputTree)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	return reflect.New(containerType).Interface(), nil
}

// CreateInputContainerByResourceName creates a container related to given resourceName
// The struct fields of target container are determined by inputMap
func CreateInputContainerByResourceName(resourceName string, inputMap interface{}) (interface{}, error) {
	modelTypeName, exists := resourceNameTypeNameMap[resourceName]
	if !exists {
		logging.Logger().Debugf("the resource name '%s' has not been associated yet", resourceName)
		return nil, fmt.Errorf("the resource name '%s' has not been associated yet", resourceName)
	}

	return CreateInputContainerByTypeName(modelTypeName, inputMap)
}

// CreateContainerByTypeName creates a container related to given typeName
// The struct fields of target container are only primitives, Struct or Array, Slice will be omitted
func CreateContainerByTypeName(typeName string) (interface{}, error) {
	return CreateOutputContainerByTypeName(typeName, "")
}

// CreateContainerByResourceName creates a container related to given resourceName
// The struct fields of target container are only primitives, Struct or Array, Slice will be omitted
func CreateContainerByResourceName(resourceName string) (interface{}, error) {
	return CreateOutputContainerByResourceName(resourceName, "")
}

// CreateOutputContainerByTypeName creates a container related to given typeName
// The struct fields of target container are determined by preloadQuery
func CreateOutputContainerByTypeName(typeName string, preloadQuery string) (interface{}, error) {
	_, exists := typeNameStructFieldListMap[typeName]
	if !exists {
		logging.Logger().Debugf("the type name '%s' has not been registered yet", typeName)
		return nil, fmt.Errorf("the type name '%s' has not been registered yet", typeName)
	}

	preloadTree := createPreloadTree(preloadQuery)
	containerType, err := createOutputContainerByTypeName(typeName, preloadTree)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	return reflect.New(containerType).Interface(), nil
}

// CreateOutputContainerByResourceName creates a container related to given resourceName
// The struct fields of target container are determined by preloadQuery
func CreateOutputContainerByResourceName(resourceName string, preloadQuery string) (interface{}, error) {
	modelTypeName, exists := resourceNameTypeNameMap[resourceName]
	if !exists {
		logging.Logger().Debugf("the resource name '%s' has not been associated yet", resourceName)
		return nil, fmt.Errorf("the resource name '%s' has not been associated yet", resourceName)
	}

	return CreateOutputContainerByTypeName(modelTypeName, preloadQuery)
}

// GetRegisteredModelByContainer returns the model related to given container
func GetRegisteredModelByContainer(container interface{}) (Model, error) {
	fieldElementType := InspectActualElementType(container)

	structMetaInformationField, exists := fieldElementType.FieldByName("StructMetaInformation")
	if !exists {
		// Caution: comment out in order to suppress at boot time
		// logging.Logger().Debug("the container does not have StructMetaInformation field, it might not be a container")
		return nil, errors.New("the container does not have StructMetaInformation field")
	}

	tag, ok := structMetaInformationField.Tag.Lookup("clay")
	if !ok {
		logging.Logger().Debug("the container does not have proper StructMetaInformation field, it might not be a container")
		return nil, errors.New("the container does not have proper StructMetaInformation field, it might not be a container")
	}

	typeName := ""

	tagStatementList := strings.Split(tag, ";")
	for _, tagStatement := range tagStatementList {
		tagKeyValue := strings.Split(tagStatement, "=")
		switch tagKeyValue[0] {
		case "type_name":
			typeName = tagKeyValue[1]
			break
		}
	}

	if len(typeName) == 0 {
		logging.Logger().Debug("the container does not have proper StructMetaInformation field, it might not be a container")
		return nil, errors.New("the container does not have proper StructMetaInformation field, it might not be a container")
	}

	model, err := GetAssociatedModelWithTypeName(typeName)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	return model, nil
}

// SetupModel setups models
// It's responsible in DB table creation, and resource/table name registration
func SetupModel(db *gorm.DB, initializerList []Initializer, modelList []Model) (*gorm.DB, error) {
	// Caution: Even if you input the inconsistent data like foreign keys do not exist,
	//          it will be registered, and never be checked this time.
	//          Todo: It requires order resolution logic like "depends on" between models.

	db.Exec("pragma foreign_keys = on")

	for _, initializer := range initializerList {
		err := initializer.DoBeforeDBMigration(db)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return db, err
		}
	}

	containerToBeMigratedList := []interface{}{}
	for _, model := range modelList {
		resourceName := model.GenerateTableName(model, db)

		if err := AssociateResourceNameWithModel(resourceName, model); err != nil {
			logging.Logger().Critical(err.Error())
			return db, err
		}

		container, err := model.GetContainerForMigration()
		if err != nil {
			logging.Logger().Critical(err.Error())
			return db, err
		}

		if container != nil {
			containerType := InspectActualElementType(container)
			_, exists := containerType.FieldByName("StructMetaInformation")
			if !exists {
				logging.Logger().Critical("the container does not have StructMetaInformation field, it might not be a container")
				return db, errors.New("the container does not have StructMetaInformation field, it might not be a container")
			}
			containerToBeMigratedList = append(containerToBeMigratedList, container)
		}
	}

	if err := db.AutoMigrate(containerToBeMigratedList...).Error; err != nil {
		logging.Logger().Critical(err.Error())
		return db, err
	}

	db.Exec("pragma foreign_keys = off;")

	tx := db.Begin()
	for _, initializer := range initializerList {
		err := initializer.DoAfterDBMigration(tx)
		if err != nil {
			tx.Rollback()
			logging.Logger().Critical(err.Error())
			return db, err
		}
	}
	tx.Commit()

	db.Exec("pragma foreign_keys = on;")

	return db, nil
}
