package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/utils/mapstruct"
)

type Node struct {
	ID         int            `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name       string         `json:"name" gorm:"not null;unique"`
	NodeTypeID int            `json:"node_type_id" gorm:"not null" sql:"type:integer references node_types(id) on delete set null"`
	NodeType   *NodeType      `json:"node_type"`
	NodePvID   int            `json:"node_pv_id" gorm:"not null" sql:"type:integer references node_pvs(id) on delete set null"`
	NodePv     *NodePv        `json:"node_pv"`
	Remark     sql.NullString `json:"remark"`
	Ports      []*Port        `json:"ports"`
	NodeGroups []*NodeGroup   `json:"node_groups" gorm:"many2many:node_group_association;"`
}

func (_ *Node) ExtractFromDesign(db *gorm.DB, designContent map[string]interface{}) error {
	nodes := []*Node{}
	if err := db.Preload("Ports").Preload("NodeGroups").Select("*").Find(&nodes).Error; err != nil {
		return err
	}
	designContent["nodes"] = nodes
	return nil
}

func (_ *Node) DeleteFromDesign(db *gorm.DB) error {
	return db.Exec("delete from nodes;").Error
}

func (_ *Node) LoadToDesign(db *gorm.DB, data interface{}) error {
	container := []*Node{}
	design := data.(*Design)
	if value, exists := design.Content["nodes"]; exists {
		if err := mapstruct.MapToStruct(value.([]interface{}), &container); err != nil {
			return err
		}
		for _, node := range container {
			node.Ports = nil
			node.NodeGroups = nil
			if err := db.Create(node).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

var NodeModel = &Node{}

func init() {
	extension.RegisterModelType(NodeModel)
	extension.RegisterDesignAccessor(NodeModel)
}
