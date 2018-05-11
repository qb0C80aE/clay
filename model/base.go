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
	"regexp"
	"strconv"
)

var regrxpExportedSymbol, _ = regexp.Compile(`^[0-9A-Z].*$`)

// Base is the base class that all model classes inherit
type Base struct {
}

// GetTypeName returns its struct type name
func (receiver *Base) GetTypeName(model extension.Model) string {
	return extension.GetActualType(model).Name()
}

// GetResourceName returns its resource/table name in URL/DB
func (receiver *Base) GetResourceName(model extension.Model) (string, error) {
	return extension.GetAssociatedResourceNameWithModel(model)
}

// GenerateTableName generates its resource/table name in URL/DB
func (receiver *Base) GenerateTableName(model extension.Model, db *gorm.DB) string {
	return db.NewScope(model).TableName()
}

// GetModelKey returns its key fields
func (receiver *Base) GetModelKey(model extension.Model, keyParameterSpecifier string) (extension.ModelKey, error) {
	return extension.GetModelKey(model, keyParameterSpecifier)
}

// GetStructFields returns its struct fields used to create containers
func (receiver *Base) GetStructFields(model extension.Model) []reflect.StructField {
	actualType := extension.GetActualType(model)
	actualValue := extension.GetActualValue(model)

	structFieldList := []reflect.StructField{}
	for i := 0; i < actualType.NumField(); i++ {
		field := actualType.Field(i)

		if !actualType.Field(i).Anonymous && regrxpExportedSymbol.MatchString(field.Name) {
			actualFieldType := field.Type
			for actualFieldType.Kind() == reflect.Ptr {
				actualFieldType = actualFieldType.Elem()
			}

			switch actualFieldType.Kind() {
			case reflect.Struct:
				actualElementType := extension.InspectActualElementType(actualValue.FieldByIndex(field.Index).Interface())
				newStructField := field
				newStructField.Type = extension.NewStructFieldTypeProxy(actualElementType.Name(), actualFieldType.Kind())
				structFieldList = append(structFieldList, newStructField)
			case reflect.Slice, reflect.Array:
				actualElementType := extension.InspectActualElementType(actualValue.FieldByIndex(field.Index).Interface())
				newStructField := field
				if actualElementType.Kind() == reflect.Struct {
					newStructField.Type = extension.NewStructFieldTypeProxy(actualElementType.Name(), actualFieldType.Kind())
					structFieldList = append(structFieldList, newStructField)
				} else {
					// primitive array/slice
					structFieldList = append(structFieldList, field)
				}
			default:
				structFieldList = append(structFieldList, field)
			}
		}
	}

	return structFieldList
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *Base) GetSingle(model extension.Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, queryFields string) (interface{}, error) {
	preloadString := urlValues.Get("preloads")

	resourceName, err := model.GetResourceName(model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	result, err := extension.CreateOutputContainerByResourceName(resourceName, preloadString)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	modelKey, err := model.GetModelKey(model, urlValues.Get("key_parameter"))
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	if err := db.Select(queryFields).First(result, fmt.Sprintf("%s = ?", modelKey.KeyParameter), parameters.ByName("key_parameter")).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	return result, nil
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (receiver *Base) GetMulti(model extension.Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, queryFields string) (interface{}, error) {
	preloadString := urlValues.Get("preloads")

	resourceName, err := model.GetResourceName(model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	elementOfResult, err := extension.CreateOutputContainerByResourceName(resourceName, preloadString)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	modelType := extension.GetActualType(elementOfResult)
	modelPointerType := reflect.PtrTo(modelType)
	sliceType := reflect.SliceOf(modelPointerType)
	slice := reflect.MakeSlice(sliceType, 0, 0)

	slicePointer := reflect.New(sliceType)
	slicePointer.Elem().Set(slice)

	if err := db.Select(queryFields).Find(slicePointer.Interface()).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	return slicePointer.Elem().Interface(), nil
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (receiver *Base) Create(_ extension.Model, db *gorm.DB, _ gin.Params, _ url.Values, inputContainer interface{}) (interface{}, error) {
	if err := db.Create(inputContainer).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	return inputContainer, nil
}

// Update corresponds HTTP PUT message and handles a request for a single resource to update the specific information
func (receiver *Base) Update(model extension.Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, inputContainer interface{}) (interface{}, error) {
	value := reflect.ValueOf(inputContainer)

	modelKey, err := model.GetModelKey(model, urlValues.Get("key_parameter"))
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	switch value.Elem().FieldByName(modelKey.KeyField).Kind() {
	case reflect.String:
		value.Elem().FieldByName(modelKey.KeyField).SetString(parameters.ByName("key_parameter"))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		id, err := strconv.Atoi(parameters.ByName("key_parameter"))
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		value.Elem().FieldByName(modelKey.KeyField).SetInt(int64(id))
	default:
		logging.Logger().Debug(err.Error())
		return nil, fmt.Errorf("the field %s does not exist, or is neither int nor string", modelKey.KeyField)
	}

	// NOTE: due to reflect.StructOf() restriction, BeforeSave does not work on containers,
	//       so delete children before save explicitly.
	if err := receiver.deleteMarkedItemsInSlices(db, value.Interface()); err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	if err := db.Save(value.Interface()).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	return inputContainer, nil
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (receiver *Base) Delete(model extension.Model, db *gorm.DB, parameters gin.Params, urlValues url.Values) error {
	modelKey, err := model.GetModelKey(model, urlValues.Get("key_parameter"))
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	resourceName, err := model.GetResourceName(model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	target, err := extension.CreateContainerByResourceName(resourceName)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	if err := db.First(target, fmt.Sprintf("%s = ?", modelKey.KeyParameter), parameters.ByName("key_parameter")).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	return db.Delete(target).Error
}

// Patch corresponds HTTP PATCH message and handles a request for a single resource to update partially the specific information
func (receiver *Base) Patch(_ extension.Model, _ *gorm.DB, _ gin.Params, _ url.Values, _ interface{}) (interface{}, error) {
	return nil, nil
}

// GetOptions corresponds HTTP OPTIONS message and handles a request for multi resources to retrieve its supported options
func (receiver *Base) GetOptions(_ extension.Model, _ *gorm.DB, _ gin.Params, _ url.Values) error {
	return nil
}

// GetCount returns the record count under current db conditions
func (receiver *Base) GetCount(model extension.Model, db *gorm.DB) (int, error) {
	resourceName, err := model.GetResourceName(model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return 0, err
	}

	var total int
	if err := db.Table(resourceName).Count(&total).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return 0, err
	}

	return total, nil
}

// ExtractFromDesign extracts the model related to this model from db
func (receiver *Base) ExtractFromDesign(model extension.Model, db *gorm.DB) (string, interface{}, error) {
	resourceName, err := model.GetResourceName(model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", nil, err
	}

	outputContainer, err := extension.CreateOutputContainerByResourceName(resourceName, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", nil, err
	}

	outputContainerType := extension.GetActualType(outputContainer)
	outputContainerTypePointerType := reflect.PtrTo(outputContainerType)
	sliceType := reflect.SliceOf(outputContainerTypePointerType)
	slice := reflect.MakeSlice(sliceType, 0, 0)

	slicePointer := reflect.New(sliceType)
	slicePointer.Elem().Set(slice)

	if err := db.Select("*").Find(slicePointer.Interface()).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return "", nil, err
	}

	return resourceName, slicePointer.Elem().Interface(), nil
}

// DeleteFromDesign deletes the model related to this model in db
func (receiver *Base) DeleteFromDesign(model extension.Model, db *gorm.DB) error {
	resourceName, err := model.GetResourceName(model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	inputContainer, err := extension.CreateOutputContainerByResourceName(resourceName, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	if err := db.Delete(inputContainer).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	return nil
}

// LoadToDesign loads the model related to this model into db
func (receiver *Base) LoadToDesign(model extension.Model, db *gorm.DB, data interface{}) error {
	resourceName, err := model.GetResourceName(model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	outputContainer, err := extension.CreateOutputContainerByResourceName(resourceName, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	outputContainerType := extension.GetActualType(outputContainer)
	outputContainerTypePointerType := reflect.PtrTo(outputContainerType)
	sliceType := reflect.SliceOf(outputContainerTypePointerType)
	slice := reflect.MakeSlice(sliceType, 0, 0)

	slicePointer := reflect.New(sliceType)
	slicePointer.Elem().Set(slice)

	design := data.(*Design)
	if value, exists := design.Content[resourceName]; exists {
		if err := mapstruct.MapToStruct(value.([]interface{}), slicePointer.Interface()); err != nil {
			logging.Logger().Debug(err.Error())
			return err
		}

		slice = slicePointer.Elem()
		size := slice.Len()
		for i := 0; i < size; i++ {
			outputContainer := slice.Index(i).Interface()
			if err := db.Set("gorm:save_associations", false).Save(outputContainer).Error; err != nil {
				logging.Logger().Debug(err.Error())
				return err
			}
		}
	}

	return nil
}

// DoBeforeDBMigration execute initialization process before DB migration
func (receiver *Base) DoBeforeDBMigration(_ *gorm.DB) error {
	return nil
}

// DoAfterDBMigration execute initialization process after DB migration
func (receiver *Base) DoAfterDBMigration(_ *gorm.DB) error {
	return nil
}

// DoBeforeRouterSetup execute initialization process before Router initialization
func (receiver *Base) DoBeforeRouterSetup(_ *gin.Engine) error {
	return nil
}

// DoAfterRouterSetup execute initialization process after Router initialization
func (receiver *Base) DoAfterRouterSetup(_ *gin.Engine) error {
	return nil
}

/*
// This logic can not specify m2m records.
// By using multiple gorm:primary_key for m2m association model, db.Delete deletes a record correctly using those keys,
// So this logic should me simply replaced with db.Delete after getting to use multiple primary keys.
func (receiver *Base) deleteMarkedItem(db *gorm.DB, itemValue reflect.Value) error {
	model, err := extension.GetRegisteredModelByContainer(itemValue.Addr().Interface())
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	keyFieldValue := itemValue.FieldByName(modelKey.KeyField)
	keyParameterValue := ""

	switch keyFieldValue.Kind() {
	case reflect.String:
		keyParameterValue = keyFieldValue.Interface().(string)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		keyParameterValue = strconv.Itoa(int(keyFieldValue.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		keyParameterValue = strconv.Itoa(int(keyFieldValue.Int()))
	default:
		logging.Logger().Debugf("the field %s does not exist, or is neither int nor string", modelKey.KeyField)
		return fmt.Errorf("the field %s does not exist, or is neither int nor string", modelKey.KeyField)
	}

	parameters := gin.Params{
		{
			Key:   "key_parameter",
			Value: keyParameterValue,
		},
	}

	if err := model.Delete(model, db, parameters, nil); err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	return nil
}
*/

func (receiver *Base) deleteMarkedItemsInSlices(db *gorm.DB, data interface{}) error {
	valueOfData := reflect.ValueOf(data)
	for valueOfData.Kind() == reflect.Ptr {
		valueOfData = valueOfData.Elem()
	}

	typeOfData := valueOfData.Type()

	for i := 0; i < typeOfData.NumField(); i++ {
		structField := typeOfData.Field(i)
		fieldValue := valueOfData.FieldByName(structField.Name)

		for fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		if fieldValue.Kind() == reflect.Slice {
			processed := reflect.New(fieldValue.Type()).Elem()
			for j := 0; j < fieldValue.Len(); j++ {
				itemValue := fieldValue.Index(j)

				itemValueToCheckIfStruct := itemValue
				for itemValueToCheckIfStruct.Kind() == reflect.Ptr {
					itemValueToCheckIfStruct = itemValueToCheckIfStruct.Elem()
				}
				if itemValueToCheckIfStruct.Kind() != reflect.Struct {
					return nil
				}

				if err := receiver.deleteMarkedItemsInSlices(db, itemValue.Interface()); err != nil {
					logging.Logger().Debug(err.Error())
					return err
				}

				for itemValue.Kind() == reflect.Ptr {
					itemValue = itemValue.Elem()
				}

				toBeDeleted := false
				toBeDeletedFieldValue := itemValue.FieldByName("ToBeDeleted")
				if toBeDeletedFieldValue.IsValid() {
					toBeDeleted = toBeDeletedFieldValue.Bool()
				}

				if toBeDeleted {
					//if err := receiver.deleteMarkedItem(db, itemValue); err != nil {
					if err := db.Delete(itemValue.Addr().Interface()).Error; err != nil {
						logging.Logger().Debug(err.Error())
						return err
					}
				} else {
					processed = reflect.Append(processed, itemValue.Addr())
				}
			}
			fieldValue.Set(processed)
		} else if fieldValue.Kind() == reflect.Struct {
			if err := receiver.deleteMarkedItemsInSlices(db, valueOfData.FieldByName(structField.Name).Interface()); err != nil {
				logging.Logger().Debug(err.Error())
				return err
			}

			toBeDeleted := false
			toBeDeletedFieldValue := fieldValue.FieldByName("ToBeDeleted")
			if toBeDeletedFieldValue.IsValid() {
				toBeDeleted = toBeDeletedFieldValue.Bool()
			}

			if toBeDeleted {
				// if err := receiver.deleteMarkedItem(db, fieldValue); err != nil {
				if err := db.Delete(fieldValue.Addr().Interface()).Error; err != nil {
					logging.Logger().Debug(err.Error())
					return err
				}
				valueOfData.FieldByName(structField.Name).Set(reflect.Zero(valueOfData.FieldByName(structField.Name).Type()))
			}
		}
	}
	return nil
}
