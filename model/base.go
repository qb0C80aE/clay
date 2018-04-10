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
	"strconv"
)

// Base is the base class that all model classes inherit
type Base struct {
	actualModel extension.Model
}

func convertContainerToModel(container interface{}) {
	typeOfContainer := extension.GetActualType(container)
	actualValueOfContainer := extension.GetActualValue(container)
	valueOfContainer := reflect.ValueOf(container)

	switch actualValueOfContainer.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < valueOfContainer.Len(); i++ {
			valueOfElement := valueOfContainer.Index(i)
			convertContainerToModel(valueOfElement.Interface())
		}
	case reflect.Struct:
		if _, ok := valueOfContainer.Interface().(extension.Model); !ok {
			panic(fmt.Errorf("%v is not extension.Model", container))
		}

		if !actualValueOfContainer.FieldByName("Base").IsValid() {
			panic(fmt.Errorf("%v is not extension.Model", container))
		}

		for i := 0; i < typeOfContainer.NumField(); i++ {
			field := typeOfContainer.Field(i)
			if field.Name == "Base" {
				continue
			}

			valueOfField := valueOfContainer.Elem().FieldByName(field.Name)
			convertContainerToModel(valueOfField.Interface())
		}

		base := Base{
			actualModel: valueOfContainer.Interface().(extension.Model),
		}
		actualValueOfContainer.FieldByName("Base").Set(reflect.ValueOf(base))
	}
}

// ConvertContainerToModel converts from container, which means, what does not have the relationship with Base, to model, what have the one recursively
func ConvertContainerToModel(container interface{}) extension.Model {
	convertContainerToModel(container)
	return container.(extension.Model)
}

// New creates a new model
func (receiver *Base) New() extension.Model {
	actualModelType := extension.GetActualType(receiver.actualModel)
	newModel := reflect.New(actualModelType)
	return ConvertContainerToModel(newModel.Interface()).(extension.Model)
}

// IsContainer tells if the instance is container, that is the model instance without Base-ActualModel relationship
func (receiver *Base) IsContainer() bool {
	return receiver.actualModel == nil
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *Base) GetSingle(db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	result := receiver.New()

	modelKey, err := extension.GetRegisteredModelKey(result)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	if err := db.Select(queryFields).First(result, fmt.Sprintf("%s = ?", modelKey.KeyParameter), parameters.ByName(modelKey.KeyParameter)).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	return result, nil
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (receiver *Base) GetMulti(db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	elementOfResult := receiver.New()

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
func (receiver *Base) Create(db *gorm.DB, _ gin.Params, _ url.Values, input extension.Model) (interface{}, error) {
	if err := db.Create(input).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	return input, nil
}

// Update corresponds HTTP PUT message and handles a request for a single resource to update the specific information
func (receiver *Base) Update(db *gorm.DB, parameters gin.Params, _ url.Values, input extension.Model) (interface{}, error) {
	value := reflect.ValueOf(input)

	modelKey, err := extension.GetRegisteredModelKey(input)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	switch value.Elem().FieldByName(modelKey.KeyField).Kind() {
	case reflect.String:
		value.Elem().FieldByName(modelKey.KeyField).SetString(parameters.ByName(modelKey.KeyParameter))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		id, err := strconv.Atoi(parameters.ByName(modelKey.KeyParameter))
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		value.Elem().FieldByName(modelKey.KeyField).SetInt(int64(id))
	default:
		logging.Logger().Debug(err.Error())
		return nil, fmt.Errorf("the field %s does not exist, or is neither int nor string", modelKey.KeyField)
	}

	if err := db.Save(value.Interface()).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	return input, nil
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (receiver *Base) Delete(db *gorm.DB, parameters gin.Params, _ url.Values) error {
	modelKey, err := extension.GetRegisteredModelKey(receiver.actualModel)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	target := receiver.New()
	if err := db.First(target, fmt.Sprintf("%s = ?", modelKey.KeyParameter), parameters.ByName(modelKey.KeyParameter)).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	return db.Delete(target).Error
}

// Patch corresponds HTTP PATCH message and handles a request for a single resource to update partially the specific information
func (receiver *Base) Patch(_ *gorm.DB, _ gin.Params, _ url.Values, _ extension.Model) (interface{}, error) {
	return nil, nil
}

// GetOptions corresponds HTTP OPTIONS message and handles a request for multi resources to retrieve its supported options
func (receiver *Base) GetOptions(_ *gorm.DB, _ gin.Params, _ url.Values) error {
	return nil
}

// GetTotal returns the count of for multi resource
func (receiver *Base) GetTotal(db *gorm.DB) (int, error) {
	var total int
	if err := db.Model(receiver.actualModel).Count(&total).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return 0, err
	}

	return total, nil
}

// ExtractFromDesign extracts the model related to this model from db
func (receiver *Base) ExtractFromDesign(db *gorm.DB) (string, interface{}, error) {
	if receiver == nil {
		logging.Logger().Criticalf("the model is a container which does not have *Base")
		return "", nil, errors.New("the model is a container which does not have *Base")
	}

	model := receiver.New()

	modelType := extension.GetActualType(model)
	modelPointerType := reflect.PtrTo(modelType)
	sliceType := reflect.SliceOf(modelPointerType)
	slice := reflect.MakeSlice(sliceType, 0, 0)

	slicePointer := reflect.New(sliceType)
	slicePointer.Elem().Set(slice)

	if err := db.Select("*").Find(slicePointer.Interface()).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return "", nil, err
	}

	resourceName := extension.GetAssociateResourceNameWithModel(model)

	return resourceName, slicePointer.Elem().Interface(), nil
}

// DeleteFromDesign deletes the model related to this model in db
func (receiver *Base) DeleteFromDesign(db *gorm.DB) error {
	if receiver == nil {
		logging.Logger().Criticalf("the model is a container which does not have *Base")
		return errors.New("the model is a container which does not have *Base")
	}

	model := receiver.New()

	if err := db.Delete(model).Error; err != nil {
		logging.Logger().Debug(err.Error())
	}

	return nil
}

// LoadToDesign loads the model related to this model into db
func (receiver *Base) LoadToDesign(db *gorm.DB, data interface{}) error {
	if receiver == nil {
		logging.Logger().Criticalf("the model is a container which does not have *Base")
		return errors.New("the model is a container which does not have *Base")
	}

	model := receiver.New()

	modelType := extension.GetActualType(model)
	modelPointerType := reflect.PtrTo(modelType)
	sliceType := reflect.SliceOf(modelPointerType)
	slice := reflect.MakeSlice(sliceType, 0, 0)

	slicePointer := reflect.New(sliceType)
	slicePointer.Elem().Set(slice)

	resourceName := extension.GetAssociateResourceNameWithModel(model)

	design := data.(*Design)
	if value, exists := design.Content[resourceName]; exists {
		if err := mapstruct.MapToStruct(value.([]interface{}), slicePointer.Interface()); err != nil {
			logging.Logger().Debug(err.Error())
			return err
		}

		slice = slicePointer.Elem()
		size := slice.Len()
		for i := 0; i < size; i++ {
			modelValue := slice.Index(i).Elem()
			fieldCount := modelValue.NumField()
			for j := 0; j < fieldCount; j++ {
				modelFieldValue := modelValue.Field(j)
				switch modelFieldValue.Kind() {
				case reflect.Array, reflect.Slice, reflect.Ptr:
					modelFieldValue.Set(reflect.Zero(modelFieldValue.Type()))
				}
			}
			model := modelValue.Interface()
			if err := db.Create(model).Error; err != nil {
				logging.Logger().Debug(err.Error())
				return err
			}
		}
	}
	return nil
}

// DoAfterDBMigration execute initialization process after DB migration
func (receiver *Base) DoAfterDBMigration(db *gorm.DB) error {
	if receiver == nil {
		logging.Logger().Criticalf("the model is a container which does not have *Base")
		return errors.New("the model is a container which does not have *Base")
	}

	return nil
}

// DoBeforeRouterSetup execute initialization process before Router initialization
func (receiver *Base) DoBeforeRouterSetup(r *gin.Engine) error {
	if receiver == nil {
		logging.Logger().Criticalf("the model is a container which does not have *Base")
		return errors.New("the model is a container which does not have *Base")
	}

	return nil
}

// DoAfterRouterSetup execute initialization process after Router initialization
func (receiver *Base) DoAfterRouterSetup(r *gin.Engine) error {
	if receiver == nil {
		logging.Logger().Criticalf("the model is a container which does not have *Base")
		return errors.New("the model is a container which does not have *Base")
	}

	return nil
}

// BeforeCreate is executed before db.Create with the model
func (receiver *Base) BeforeCreate(tx *gorm.DB) error {
	return receiver.deleteMarkedItemsInSlices(tx, receiver.actualModel)
}

// BeforeSave is executed before db.Save with the model
func (receiver *Base) BeforeSave(tx *gorm.DB) error {
	return receiver.deleteMarkedItemsInSlices(tx, receiver.actualModel)
}

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
					model := ConvertContainerToModel(itemValue.Addr().Interface()).(extension.Model)

					modelKey, err := extension.GetRegisteredModelKey(model)
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
							Key:   modelKey.KeyParameter,
							Value: keyParameterValue,
						},
					}

					if err := model.Delete(db, parameters, nil); err != nil {
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
		}
	}
	return nil
}
