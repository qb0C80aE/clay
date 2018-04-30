package model

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"reflect"
)

// UserDefinedModel is the model class what handles user defined models
// Todo: Note that current Go implementation, reflect.StructOf does not generate wrapper methods to embedded structure
// Todo: therefore, methods like BeforeSave do not work.
type UserDefinedModel struct {
	Base
	resourceName                string
	typeName                    string
	keyParameter                string
	toBeMigrated                bool
	isControllerEnabled         bool
	isDesignAccessDisabled      bool
	sqlBeforeMigration          string
	sqlAfterMigration           string
	sqlWhereForDesignExtraction string
	sqlWhereForDesignDeletion   string
	structFieldList             []reflect.StructField
}

// NewUserDefinedModel creates an user defined model instance
func NewUserDefinedModel() *UserDefinedModel {
	return &UserDefinedModel{}
}

// GetTypeName returns its struct type name
func (receiver *UserDefinedModel) GetTypeName(_ extension.Model) string {
	return receiver.typeName
}

// GetResourceName returns its resource/table name in URL/DB
func (receiver *UserDefinedModel) GetResourceName(_ extension.Model) (string, error) {
	return receiver.resourceName, nil
}

// GenerateTableName generates its resource/table name in URL/DB
func (receiver *UserDefinedModel) GenerateTableName(_ extension.Model, _ *gorm.DB) string {
	return receiver.resourceName
}

// GetStructFields returns its struct fields used to create containers
func (receiver *UserDefinedModel) GetStructFields(model extension.Model) []reflect.StructField {
	structFieldList := make([]reflect.StructField, len(receiver.structFieldList), len(receiver.structFieldList))
	for i := 0; i < len(receiver.structFieldList); i++ {
		structFieldList[i] = receiver.structFieldList[i]
	}

	return structFieldList
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *UserDefinedModel) GetContainerForMigration() (interface{}, error) {
	if receiver.toBeMigrated {
		return extension.CreateContainerByTypeName(receiver.GetTypeName(receiver))
	}

	return nil, nil
}

// DoBeforeDBMigration execute initialization process before DB migration
func (receiver *UserDefinedModel) DoBeforeDBMigration(db *gorm.DB) error {
	if len(receiver.sqlBeforeMigration) > 0 {
		return db.Exec(receiver.sqlBeforeMigration).Error
	}

	return nil
}

// DoAfterDBMigration execute initialization process after DB migration
func (receiver *UserDefinedModel) DoAfterDBMigration(db *gorm.DB) error {
	if len(receiver.sqlAfterMigration) > 0 {
		return db.Exec(receiver.sqlAfterMigration).Error
	}

	return nil
}

// ToBeMigrated returns if it should be migrated or not
func (receiver *UserDefinedModel) ToBeMigrated() bool {
	return receiver.toBeMigrated
}

// IsControllerEnabled returns if the correspond controller should be generated or registered
func (receiver *UserDefinedModel) IsControllerEnabled() bool {
	return receiver.isControllerEnabled
}

// ExtractFromDesign extracts the model related to this model from db
func (receiver *UserDefinedModel) ExtractFromDesign(model extension.Model, db *gorm.DB) (string, interface{}, error) {
	resourceName, err := model.GetResourceName(model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", nil, err
	}

	if receiver.isDesignAccessDisabled {
		logging.Logger().Debugf("design access for %s is disabled, skipped", resourceName)
		return "", nil, nil
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

	if len(receiver.sqlWhereForDesignExtraction) == 0 {
		if err := db.Select("*").Find(slicePointer.Interface()).Error; err != nil {
			logging.Logger().Debug(err.Error())
			return "", nil, err
		}
	} else {
		if err := db.Select("*").Where(receiver.sqlWhereForDesignExtraction).Find(slicePointer.Interface()).Error; err != nil {
			logging.Logger().Debug(err.Error())
			return "", nil, err
		}
	}

	return resourceName, slicePointer.Elem().Interface(), nil
}

// DeleteFromDesign deletes the model related to this model in db
func (receiver *UserDefinedModel) DeleteFromDesign(model extension.Model, db *gorm.DB) error {
	resourceName, err := model.GetResourceName(model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	if receiver.isDesignAccessDisabled {
		logging.Logger().Debugf("design access for %s is disabled, skipped", resourceName)
		return nil
	}

	inputContainer, err := extension.CreateOutputContainerByResourceName(resourceName, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	if len(receiver.sqlWhereForDesignDeletion) == 0 {
		if err := db.Delete(inputContainer).Error; err != nil {
			logging.Logger().Debug(err.Error())
			return err
		}
	} else {
		if err := db.Where(receiver.sqlWhereForDesignDeletion).Delete(inputContainer).Error; err != nil {
			logging.Logger().Debug(err.Error())
			return err
		}
	}

	return nil
}

func init() {
}
