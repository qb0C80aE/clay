package logics

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/models"
)

func GetDesign(db *gorm.DB, _ string, queryFields string) (interface{}, error) {

	design := &models.Design{
		Content: map[string]interface{}{},
	}

	designExtractors := extension.GetDesignExtractors()
	for _, extractor := range designExtractors {
		if err := extractor(db, design.Content); err != nil {
			return nil, err
		}
	}

	return design, nil
}

func UpdateDesign(db *gorm.DB, _ string, data interface{}) (interface{}, error) {
	design := data.(*models.Design)

	designDeleters := extension.GetDesignDeleters()
	for _, deleter := range designDeleters {
		if err := deleter(db); err != nil {
			return nil, err
		}
	}

	designLoaders := extension.GetDesignLoaders()
	for _, loader := range designLoaders {
		if err := loader(db, design); err != nil {
			return nil, err
		}
	}

	return design, nil
}

func DeleteDesign(db *gorm.DB, _ string) error {
	designDeleters := extension.GetDesignDeleters()
	for _, deleter := range designDeleters {
		if err := deleter(db); err != nil {
			return err
		}
	}

	return nil
}
