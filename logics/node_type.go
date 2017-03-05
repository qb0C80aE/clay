package logics

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/models"
	"strconv"
)

type NodeTypeLogic struct {
}

func NewNodeTypeLogic() *NodeTypeLogic {
	return &NodeTypeLogic{}
}

func (_ *NodeTypeLogic) GetSingle(db *gorm.DB, id string, queryFields string) (interface{}, error) {

	nodeType := &models.NodeType{}

	if err := db.Select(queryFields).First(nodeType, id).Error; err != nil {
		return nil, err
	}

	return nodeType, nil

}

func (_ *NodeTypeLogic) GetMulti(db *gorm.DB, queryFields string) ([]interface{}, error) {

	nodeTypes := []*models.NodeType{}

	if err := db.Select(queryFields).Find(&nodeTypes).Error; err != nil {
		return nil, err
	}

	result := make([]interface{}, len(nodeTypes))
	for i, data := range nodeTypes {
		result[i] = data
	}

	return result, nil

}

func (_ *NodeTypeLogic) Create(db *gorm.DB, data interface{}) (interface{}, error) {

	nodeType := data.(*models.NodeType)

	if err := db.Create(nodeType).Error; err != nil {
		return nil, err
	}

	return nodeType, nil

}

func (_ *NodeTypeLogic) Update(db *gorm.DB, id string, data interface{}) (interface{}, error) {

	nodeType := data.(*models.NodeType)
	nodeType.ID, _ = strconv.Atoi(id)

	if err := db.Save(&nodeType).Error; err != nil {
		return nil, err
	}

	return nodeType, nil

}

func (_ *NodeTypeLogic) Delete(db *gorm.DB, id string) error {

	nodeType := &models.NodeType{}

	if err := db.First(&nodeType, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&nodeType).Error; err != nil {
		return err
	}

	return nil

}

func (_ *NodeTypeLogic) Patch(_ *gorm.DB, _ string, _ string) (interface{}, error) {
	return nil, nil
}

func (_ *NodeTypeLogic) Options(db *gorm.DB) error {
	return nil
}
