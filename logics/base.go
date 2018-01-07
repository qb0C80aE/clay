package logics

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/models"
	"github.com/qb0C80aE/clay/utils/mapstruct"
	"net/url"
	"reflect"
	"strconv"
)

// BaseLogic is the base class that all logic classes inherit
type BaseLogic struct {
	model interface{}
}

// NewBaseLogic creates a new instance of BaseLogic
func NewBaseLogic(model interface{}) *BaseLogic {
	logic := &BaseLogic{
		model: model,
	}
	return logic
}

// ResourceName returns its target resource name
func (logic *BaseLogic) ResourceName() string {
	return extensions.RegisteredResourceName(logic.model)
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (logic *BaseLogic) GetSingle(db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	result, err := extensions.CreateModel(logic.ResourceName())
	if err != nil {
		return nil, err
	}

	modelKey, err := extensions.RegisteredModelKey(result)
	if err != nil {
		return nil, err
	}

	if err := db.Select(queryFields).First(result, fmt.Sprintf("%s = ?", modelKey.KeyParameter), parameters.ByName(modelKey.KeyParameter)).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (logic *BaseLogic) GetMulti(db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	model, err := extensions.CreateModel(logic.ResourceName())
	if err != nil {
		return nil, err
	}

	modelType := extensions.ModelType(model)
	modelPointerType := reflect.PtrTo(modelType)
	sliceType := reflect.SliceOf(modelPointerType)
	slice := reflect.MakeSlice(sliceType, 0, 0)

	slicePointer := reflect.New(sliceType)
	slicePointer.Elem().Set(slice)

	if err := db.Select(queryFields).Find(slicePointer.Interface()).Error; err != nil {
		return nil, err
	}

	return slicePointer.Elem().Interface(), nil
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (logic *BaseLogic) Create(db *gorm.DB, parameters gin.Params, _ url.Values, data interface{}) (interface{}, error) {
	if err := db.Create(data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// Update corresponds HTTP PUT message and handles a request for a single resource to update the specific information
func (logic *BaseLogic) Update(db *gorm.DB, parameters gin.Params, _ url.Values, data interface{}) (interface{}, error) {
	value := reflect.ValueOf(data)

	modelKey, err := extensions.RegisteredModelKey(data)
	if err != nil {
		return nil, err
	}

	switch value.Elem().FieldByName(modelKey.KeyField).Kind() {
	case reflect.String:
		value.Elem().FieldByName(modelKey.KeyField).SetString(parameters.ByName(modelKey.KeyParameter))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		id, err := strconv.Atoi(parameters.ByName(modelKey.KeyParameter))
		if err != nil {
			return nil, err
		}
		value.Elem().FieldByName(modelKey.KeyField).SetInt(int64(id))
	default:
		return nil, fmt.Errorf("the field %s does not exist, or is neither int nor string", modelKey.KeyField)
	}

	if err := db.Save(value.Interface()).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (logic *BaseLogic) Delete(db *gorm.DB, parameters gin.Params, _ url.Values) error {
	model, err := extensions.CreateModel(logic.ResourceName())
	if err != nil {
		return err
	}

	modelKey, err := extensions.RegisteredModelKey(model)
	if err != nil {
		return err
	}

	if err := db.First(model, fmt.Sprintf("%s = ?", modelKey.KeyParameter), parameters.ByName(modelKey.KeyParameter)).Error; err != nil {
		return err
	}

	return db.Delete(model).Error
}

// Patch corresponds HTTP PATCH message and handles a request for a single resource to update partially the specific information
func (logic *BaseLogic) Patch(_ *gorm.DB, _ gin.Params, _ url.Values) (interface{}, error) {
	return nil, nil
}

// Options corresponds HTTP OPTIONS message and handles a request for multi resources to retrieve its supported options
func (logic *BaseLogic) Options(_ *gorm.DB, _ gin.Params, _ url.Values) error {
	return nil
}

// Total returns the count of for multi resource
func (logic *BaseLogic) Total(db *gorm.DB, model interface{}) (int, error) {
	var total int
	if err := db.Model(model).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// ExtractFromDesign extracts the model related to this logic from db
func (logic *BaseLogic) ExtractFromDesign(db *gorm.DB) (string, interface{}, error) {
	model, err := extensions.CreateModel(logic.ResourceName())
	if err != nil {
		return "", nil, err
	}

	modelType := extensions.ModelType(model)
	modelPointerType := reflect.PtrTo(modelType)
	sliceType := reflect.SliceOf(modelPointerType)
	slice := reflect.MakeSlice(sliceType, 0, 0)

	slicePointer := reflect.New(sliceType)
	slicePointer.Elem().Set(slice)

	if err := db.Select("*").Find(slicePointer.Interface()).Error; err != nil {
		return "", nil, err
	}
	return logic.ResourceName(), slicePointer.Elem().Interface(), nil
}

// DeleteFromDesign deletes the model related to this logic in db
func (logic *BaseLogic) DeleteFromDesign(db *gorm.DB) error {
	model, err := extensions.CreateModel(logic.ResourceName())
	if err != nil {
		return err
	}

	return db.Delete(model).Error
}

// LoadToDesign loads the model related to this logic into db
func (logic *BaseLogic) LoadToDesign(db *gorm.DB, data interface{}) error {
	model, err := extensions.CreateModel(logic.ResourceName())
	if err != nil {
		return err
	}

	modelType := extensions.ModelType(model)
	modelPointerType := reflect.PtrTo(modelType)
	sliceType := reflect.SliceOf(modelPointerType)
	slice := reflect.MakeSlice(sliceType, 0, 0)

	slicePointer := reflect.New(sliceType)
	slicePointer.Elem().Set(slice)

	design := data.(*models.Design)
	if value, exists := design.Content[logic.ResourceName()]; exists {
		if err := mapstruct.MapToStruct(value.([]interface{}), slicePointer.Interface()); err != nil {
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
				return err
			}
		}
	}
	return nil
}

func init() {
}
