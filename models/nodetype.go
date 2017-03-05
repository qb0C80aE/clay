package models

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/utils/mapstruct"
)

type NodeType struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"not null"`
}

func (_ *NodeType) ExtractFromDesign(db *gorm.DB, designContent map[string]interface{}) error {
	nodeTypes := []*NodeType{}
	if err := db.Select("*").Find(&nodeTypes).Error; err != nil {
		return err
	}
	designContent["node_types"] = nodeTypes
	return nil
}

func (_ *NodeType) DeleteFromDesign(db *gorm.DB) error {
	return db.Exec("delete from node_types;").Error
}

func (_ *NodeType) LoadToDesign(db *gorm.DB, data interface{}) error {
	container := []*NodeType{}
	design := data.(*Design)
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

var NodeTypeModel = &NodeType{}

func init() {
	extension.RegisterModelType(NodeTypeModel)
	extension.RegisterDesignAccessor(NodeTypeModel)
}
