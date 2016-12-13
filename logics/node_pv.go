package logics

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/models"
	"strconv"
)

func GetNodePvs(db *gorm.DB, queryFields string) ([]interface{}, error) {

	nodePvs := []*models.NodePv{}

	if err := db.Select(queryFields).Find(&nodePvs).Error; err != nil {
		return nil, err
	}

	result := make([]interface{}, len(nodePvs))
	for i, data := range nodePvs {
		result[i] = data
	}

	return result, nil

}

func GetNodePv(db *gorm.DB, id string, queryFields string) (interface{}, error) {

	nodePv := &models.NodePv{}

	if err := db.Select(queryFields).First(nodePv, id).Error; err != nil {
		return nil, err
	}

	return nodePv, nil

}

func CreateNodePv(db *gorm.DB, data interface{}) (interface{}, error) {

	nodePv := data.(*models.NodePv)

	if err := db.Create(nodePv).Error; err != nil {
		return nil, err
	}

	return nodePv, nil

}

func UpdateNodePv(db *gorm.DB, id string, data interface{}) (interface{}, error) {

	nodePv := data.(*models.NodePv)
	nodePv.ID, _ = strconv.Atoi(id)

	if err := db.Save(&nodePv).Error; err != nil {
		return nil, err
	}

	return nodePv, nil

}

func DeleteNodePv(db *gorm.DB, id string) error {

	nodePv := &models.NodePv{}

	if err := db.First(&nodePv, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&nodePv).Error; err != nil {
		return err
	}

	return nil

}
