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

// CreateModel creates a concrete model with Base
func CreateModel(actualModel extension.Model) extension.Model {
	actualModelValue := reflect.ValueOf(actualModel).Elem()
	base := &Base{
		actualModel: actualModel,
	}
	actualModelValue.FieldByName("Base").Set(reflect.ValueOf(base))
	return actualModel
}

// NewModelContainer creates a container, which means a model without Base
func (receiver *Base) NewModelContainer() extension.Model {
	if receiver == nil {
		logging.Logger().Criticalf("the model is a container which does not have *Base")
		panic("the model is a container which does not have *Base")
	}

	actualModelType := extension.GetActualType(receiver.actualModel)
	newModel := reflect.New(actualModelType)
	return newModel.Interface().(extension.Model)
}

// ExecuteActualGetSingle executes GetSingle with container check
func (receiver *Base) ExecuteActualGetSingle(db *gorm.DB, parameters gin.Params, urlValues url.Values, queryFields string) (interface{}, error) {
	if (receiver == nil) || (receiver.actualModel == nil) {
		logging.Logger().Criticalf("the model is a container which does not have *Base")
		return nil, errors.New("the model is a container which does not have *Base")
	}

	return receiver.actualModel.GetSingle(db, parameters, urlValues, queryFields)
}

// ExecuteActualGetMulti executes GetMulti with container check
func (receiver *Base) ExecuteActualGetMulti(db *gorm.DB, parameters gin.Params, urlValues url.Values, queryFields string) (interface{}, error) {
	if (receiver == nil) || (receiver.actualModel == nil) {
		logging.Logger().Criticalf("the model is a container which does not have *Base")
		return nil, errors.New("the model is a container which does not have *Base")
	}

	return receiver.actualModel.GetMulti(db, parameters, urlValues, queryFields)
}

// ExecuteActualCreate executes Create with container check
func (receiver *Base) ExecuteActualCreate(db *gorm.DB, parameters gin.Params, urlValues url.Values, input extension.Model) (interface{}, error) {
	if (receiver == nil) || (receiver.actualModel == nil) {
		logging.Logger().Criticalf("the model is a container which does not have *Base")
		return nil, errors.New("the model is a container which does not have *Base")
	}

	return receiver.actualModel.Create(db, parameters, urlValues, input)
}

// ExecuteActualUpdate executes Update with container check
func (receiver *Base) ExecuteActualUpdate(db *gorm.DB, parameters gin.Params, urlValues url.Values, input extension.Model) (interface{}, error) {
	if (receiver == nil) || (receiver.actualModel == nil) {
		logging.Logger().Criticalf("the model is a container which does not have *Base")
		return nil, errors.New("the model is a container which does not have *Base")
	}

	return receiver.actualModel.Update(db, parameters, urlValues, input)
}

// ExecuteActualDelete executes Delete with container check
func (receiver *Base) ExecuteActualDelete(db *gorm.DB, parameters gin.Params, urlValues url.Values) error {
	if (receiver == nil) || (receiver.actualModel == nil) {
		logging.Logger().Criticalf("the model is a container which does not have *Base")
		return errors.New("the model is a container which does not have *Base")
	}

	return receiver.actualModel.Delete(db, parameters, urlValues)
}

// ExecuteActualPatch executes Patch with container check
func (receiver *Base) ExecuteActualPatch(db *gorm.DB, parameters gin.Params, urlValues url.Values, input extension.Model) (interface{}, error) {
	if (receiver == nil) || (receiver.actualModel == nil) {
		logging.Logger().Criticalf("the model is a container which does not have *Base")
		return nil, errors.New("the model is a container which does not have *Base")
	}

	return receiver.actualModel.Patch(db, parameters, urlValues, input)
}

// ExecuteActualGetOptions executes GetOptions with container check
func (receiver *Base) ExecuteActualGetOptions(db *gorm.DB, parameters gin.Params, urlValues url.Values) error {
	if (receiver == nil) || (receiver.actualModel == nil) {
		logging.Logger().Criticalf("the model is a container which does not have *Base")
		return errors.New("the model is a container which does not have *Base")
	}

	return receiver.actualModel.GetOptions(db, parameters, urlValues)
}

// ExecuteActualGetTotal executes GetTotal with container check
func (receiver *Base) ExecuteActualGetTotal(db *gorm.DB) (int, error) {
	if (receiver == nil) || (receiver.actualModel == nil) {
		logging.Logger().Criticalf("the model is a container which does not have *Base")
		return 0, errors.New("the model is a container which does not have *Base")
	}

	return receiver.actualModel.GetTotal(db)
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *Base) GetSingle(db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	result := receiver.NewModelContainer()

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
	elementOfResult := receiver.NewModelContainer()

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

	target := receiver.NewModelContainer()
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

	model := receiver.NewModelContainer()

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

	model := receiver.NewModelContainer()

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

	model := receiver.NewModelContainer()

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
