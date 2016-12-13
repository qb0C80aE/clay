package logics

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/models"
	"strconv"
)

func GetNodeGroups(db *gorm.DB, queryFields string) ([]interface{}, error) {

	nodeGroups := []*models.NodeGroup{}

	if err := db.Select(queryFields).Find(&nodeGroups).Error; err != nil {
		return nil, err
	}

	result := make([]interface{}, len(nodeGroups))
	for i, data := range nodeGroups {
		result[i] = data
	}

	return result, nil

}

func GetNodeGroup(db *gorm.DB, id string, queryFields string) (interface{}, error) {

	nodeGroup := &models.NodeGroup{}

	if err := db.Select(queryFields).First(nodeGroup, id).Error; err != nil {
		return nil, err
	}

	return nodeGroup, nil

}

func CreateNodeGroup(db *gorm.DB, data interface{}) (interface{}, error) {

	nodeGroup := data.(*models.NodeGroup)

	inputNodes := nodeGroup.Nodes
	nodeGroup.Nodes = nil

	if err := db.Create(nodeGroup).Error; err != nil {
		return nil, err
	}

	nodeGroup.Nodes = inputNodes

	if err := db.Model(nodeGroup).Association("Nodes").Append(nodeGroup.Nodes).Error; err != nil {
		return nil, err
	}

	return nodeGroup, nil

}

func UpdateNodeGroup(db *gorm.DB, id string, data interface{}) (interface{}, error) {

	nodeGroup := data.(*models.NodeGroup)
	nodeGroup.ID, _ = strconv.Atoi(id)

	inputNodes := nodeGroup.Nodes
	nodeGroup.Nodes = nil

	if err := db.Save(nodeGroup).Error; err != nil {
		return nil, err
	}

	if err := db.Model(nodeGroup).Association("Nodes").Clear().Error; err != nil {
		return nil, err
	}

	nodeGroup.Nodes = inputNodes

	if err := db.Model(nodeGroup).Association("Nodes").Append(nodeGroup.Nodes).Error; err != nil {
		return nil, err
	}

	return nodeGroup, nil

}

func DeleteNodeGroup(db *gorm.DB, id string) error {

	nodeGroup := &models.NodeGroup{}

	if err := db.First(nodeGroup, id).Error; err != nil {
		return err
	}

	if err := db.Model(&nodeGroup).Association("Nodes").Clear().Error; err != nil {
		return err
	}

	nodeGroup.Nodes = nil

	if err := db.Delete(&nodeGroup).Error; err != nil {
		return err
	}

	return nil

}
