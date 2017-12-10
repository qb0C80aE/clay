package logics

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extensions"
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

	if err := db.Select(queryFields).First(result, parameters.ByName("id")).Error; err != nil {
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
	id, err := strconv.Atoi(parameters.ByName("id"))
	if err != nil {
		return nil, err
	}

	value := reflect.ValueOf(data)
	value.Elem().FieldByName("ID").SetInt(int64(id))

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

	if err := db.First(model, parameters.ByName("id")).Error; err != nil {
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

func init() {
}
