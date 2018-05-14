package modelstore

import (
	"errors"
	"github.com/jinzhu/gorm"
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/helper"
	"github.com/qb0C80aE/clay/logging"
	"net/url"
	"reflect"
)

// ModelStore handles model retrieving operation
type ModelStore struct {
	db *gorm.DB
}

// NewModelStore creates and returns a new ModelStore instance
func NewModelStore(db *gorm.DB) *ModelStore {
	return &ModelStore{db: db}
}

// Single retrieves a single record from DB
func (receiver *ModelStore) Single(pathInterface interface{}, queryInterface interface{}) (interface{}, error) {
	path := pathInterface.(string)

	controller, err := extension.GetAssociatedControllerWithPath(path)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	singleURL, err := controller.GetResourceSingleURL()
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	parameters, err := extension.CreateParametersFromPathAntRoute(path, singleURL)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	urlValues, err := url.ParseQuery(queryInterface.(string))
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	model := controller.GetModel()

	parameter, err := dbpkg.NewParameter(urlValues)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	// single resets db conditions like preloads, so you should use this method in GetSingle or GetMulti only,
	// and note that all conditions go away after this method.
	db := receiver.db.New()
	db, err = parameter.Paginate(db)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	db = parameter.SetPreloads(db)
	db = parameter.FilterFields(db)
	fields := helper.ParseFields(parameter.DefaultQuery(urlValues, "fields", "*"))
	queryFields := helper.QueryFields(model, fields)

	result, err := model.GetSingle(model, db, parameters, urlValues, queryFields)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	return result, nil
}

// Multi retrieves multiple record from DB
func (receiver *ModelStore) Multi(pathInterface interface{}, queryInterface interface{}) (interface{}, error) {
	path := pathInterface.(string)

	controller, err := extension.GetAssociatedControllerWithPath(path)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	multiURL, err := controller.GetResourceMultiURL()
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	parameters, err := extension.CreateParametersFromPathAntRoute(path, multiURL)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	urlValues, err := url.ParseQuery(queryInterface.(string))
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	model := controller.GetModel()

	parameter, err := dbpkg.NewParameter(urlValues)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	// multi resets db conditions like preloads, so you should use this method in GetSingle or GetMulti only,
	// and note that all conditions go away after this method.
	db := receiver.db.New()
	db, err = parameter.Paginate(db)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	db = parameter.SetPreloads(db)
	db = parameter.SortRecords(db)
	db = parameter.FilterFields(db)
	fields := helper.ParseFields(parameter.DefaultQuery(urlValues, "fields", "*"))
	queryFields := helper.QueryFields(model, fields)
	result, err := model.GetMulti(model, db, parameters, urlValues, queryFields)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	// reset all conditions in order to get the total number of records
	db = db.New()
	total, err := model.GetCount(model, db)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	// reset conditions except for limit and offset in order to get the record count before limitation
	db = db.New()
	db = parameter.SetPreloads(db)
	db = parameter.SortRecords(db)
	db = parameter.FilterFields(db)
	countBeforePagination, err := model.GetCount(model, db)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	type multiResult struct {
		Records               interface{}
		Total                 interface{}
		CountBeforePagination interface{}
	}

	multiResultObject := &multiResult{
		Records: result,
		Total:   total,
		CountBeforePagination: countBeforePagination,
	}

	return multiResultObject, nil
}

// First retrieves the first record under specific conditions from DB
func (receiver *ModelStore) First(pathInterface interface{}, queryInterface interface{}) (interface{}, error) {
	path := pathInterface.(string)

	controller, err := extension.GetAssociatedControllerWithPath(path)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	multiURL, err := controller.GetResourceMultiURL()
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	parameters, err := extension.CreateParametersFromPathAntRoute(path, multiURL)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	urlValues, err := url.ParseQuery(queryInterface.(string))
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	model := controller.GetModel()

	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	parameter, err := dbpkg.NewParameter(urlValues)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	// first resets db conditions like preloads, so you should use this method in GetSingle or GetMulti only,
	// and note that all conditions go away after this method.
	db := receiver.db.New()
	db, err = parameter.Paginate(db)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	db = parameter.SetPreloads(db)
	db = parameter.SortRecords(db)
	db = parameter.FilterFields(db)
	fields := helper.ParseFields(parameter.DefaultQuery(urlValues, "fields", "*"))
	queryFields := helper.QueryFields(model, fields)
	result, err := model.GetMulti(model, db, parameters, urlValues, queryFields)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	resultValue := reflect.ValueOf(result)
	if resultValue.Len() == 0 {
		logging.Logger().Debug("no record selected")
		return nil, errors.New("no record selected")
	}

	return resultValue.Index(0).Interface(), nil
}
