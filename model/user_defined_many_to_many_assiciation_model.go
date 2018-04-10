package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"net/url"
	"reflect"
	"strconv"
)

// UserDefinedManyToManyAssociationModel is the model class what represents many to many association between models
type UserDefinedManyToManyAssociationModel struct {
	UserDefinedModel
	leftResourceName     string
	leftResourceKeyName  string
	rightResourceName    string
	rightResourceKeyName string
}

// NewUserDefinedManyToManyAssociationModel creates an user defined many to many association model instance
func NewUserDefinedManyToManyAssociationModel() *UserDefinedManyToManyAssociationModel {
	return &UserDefinedManyToManyAssociationModel{}
}

// GetLeftResourceName returns its parent resource name
func (receiver *UserDefinedManyToManyAssociationModel) GetLeftResourceName() string {
	return receiver.leftResourceName
}

// GetLeftResourceKeyName returns its parent resource key name
func (receiver *UserDefinedManyToManyAssociationModel) GetLeftResourceKeyName() string {
	return receiver.leftResourceKeyName
}

// GetRightResourceName returns its child resource name
func (receiver *UserDefinedManyToManyAssociationModel) GetRightResourceName() string {
	return receiver.rightResourceName
}

// GetRightResourceKeyName returns its child resource name
func (receiver *UserDefinedManyToManyAssociationModel) GetRightResourceKeyName() string {
	return receiver.rightResourceKeyName
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (receiver *UserDefinedManyToManyAssociationModel) GetMulti(model extension.Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, queryFields string) (interface{}, error) {
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

	if (len(parameters.ByName(receiver.leftResourceKeyName)) == 0) && (len(parameters.ByName(receiver.rightResourceKeyName)) == 0) {
		if err := db.Select(queryFields).Find(slicePointer.Interface()).Error; err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
	} else if len(parameters.ByName(receiver.leftResourceKeyName)) == 0 {
		if err := db.Select(queryFields).Find(slicePointer.Interface(), fmt.Sprintf("%s = ?", receiver.rightResourceKeyName), parameters.ByName(receiver.rightResourceKeyName)).Error; err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
	} else {
		if err := db.Select(queryFields).Find(slicePointer.Interface(), fmt.Sprintf("%s = ?", receiver.leftResourceKeyName), parameters.ByName(receiver.leftResourceKeyName)).Error; err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
	}

	return slicePointer.Elem().Interface(), nil
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (receiver *UserDefinedManyToManyAssociationModel) Delete(model extension.Model, db *gorm.DB, parameters gin.Params, _ url.Values) error {
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

	leftID, err := strconv.Atoi(parameters.ByName(receiver.leftResourceKeyName))
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	rightID, err := strconv.Atoi(parameters.ByName(receiver.rightResourceKeyName))
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	if err := db.First(target, fmt.Sprintf("%s = ? and %s = ?", receiver.leftResourceKeyName, receiver.rightResourceKeyName), leftID, rightID).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	return db.Delete(target, fmt.Sprintf("%s = ? and %s = ?", receiver.leftResourceKeyName, receiver.rightResourceKeyName), leftID, rightID).Error
}

func init() {
}
