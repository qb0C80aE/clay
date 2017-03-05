package logics

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/models"
	"github.com/qb0C80aE/clay/utils/mapstruct"
	"strconv"
)

type NodeTypeLogic struct {
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

func (_ *NodeTypeLogic) ExtractFromDesign(db *gorm.DB, designContent map[string]interface{}) error {
	nodeTypes := []*models.NodeType{}
	if err := db.Select("*").Find(&nodeTypes).Error; err != nil {
		return err
	}
	designContent["node_types"] = nodeTypes
	return nil
}

func (_ *NodeTypeLogic) DeleteFromDesign(db *gorm.DB) error {
	return db.Exec("delete from node_types;").Error
}

func (_ *NodeTypeLogic) LoadToDesign(db *gorm.DB, data interface{}) error {
	container := []*models.NodeType{}
	design := data.(*models.Design)
	if value, exists := design.Content["node_types"]; exists {
		if err := mapstruct.MapToStruct(value.([]interface{}), &container); err != nil {
			return err
		}
		for _, nodeType := range container {
			if err := db.Create(nodeType).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

var NodeTypeLogicInstance = &NodeTypeLogic{}

func init() {
	extension.RegisterDesignAccessor(NodeTypeLogicInstance)
}
