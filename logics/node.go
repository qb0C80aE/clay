package logics

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/models"
	"strconv"
)

func updateNodeLink(db *gorm.DB, node *models.Node) error {

	for _, inputPort := range node.Ports {
		if err := updateLink(db, inputPort); err != nil {
			return err
		}
	}

	return nil

}

func GetNodes(db *gorm.DB, queryFields string) ([]interface{}, error) {

	nodes := []*models.Node{}

	if err := db.Select(queryFields).Find(&nodes).Error; err != nil {
		return nil, err
	}

	result := make([]interface{}, len(nodes))
	for i, data := range nodes {
		result[i] = data
	}

	return result, nil

}

func GetNode(db *gorm.DB, id string, queryFields string) (interface{}, error) {

	node := &models.Node{}

	if err := db.Select(queryFields).First(node, id).Error; err != nil {
		return nil, err
	}

	return node, nil

}

func CreateNode(db *gorm.DB, data interface{}) (interface{}, error) {

	node := data.(*models.Node)

	inputNodeGroups := node.NodeGroups
	node.NodeGroups = nil

	if err := db.Create(node).Error; err != nil {
		return nil, err
	}

	node.NodeGroups = inputNodeGroups

	if err := db.Model(node).Association("NodeGroups").Append(node.NodeGroups).Error; err != nil {
		return nil, err
	}


	if err := updateNodeLink(db, node); err != nil {
		return nil, err
	}

	return node, nil

}

func UpdateNode(db *gorm.DB, id string, data interface{}) (interface{}, error) {

	node := data.(*models.Node)
	node.ID, _ = strconv.Atoi(id)

	inputNodeGroups := node.NodeGroups
	node.NodeGroups = nil

	if err := db.Save(&node).Error; err != nil {
		return nil, err
	}

	if err := db.Model(node).Association("NodeGroups").Clear().Error; err != nil {
		return nil, err
	}

	node.NodeGroups = inputNodeGroups

	if err := db.Model(node).Association("NodeGroups").Append(node.NodeGroups).Error; err != nil {
		return nil, err
	}

	if err := updateNodeLink(db, node); err != nil {
		return err, nil
	}

	return node, nil

}

func DeleteNode(db *gorm.DB, id string) error {

	node := &models.Node{}

	if err := db.First(&node, id).Error; err != nil {
		return err
	}

	if err := db.Model(&node).Association("NodeGroups").Clear().Error; err != nil {
		return err
	}

	node.NodeGroups = nil

	if err := db.Delete(&node).Error; err != nil {
		return err
	}

	return nil

}
