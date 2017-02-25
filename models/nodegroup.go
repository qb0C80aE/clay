package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/mohae/deepcopy"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/utils/mapstruct"
	"reflect"
)

type NodeGroup struct {
	ID     int            `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name   string         `json:"name" gorm:"not null;unique"`
	Remark sql.NullString `json:"remark"`
	Nodes  []*Node        `json:"nodes" gorm:"many2many:node_group_association;"`
}

func extraceNodeGroupsFromDesign(db *gorm.DB, designContent map[string]interface{}) error {
	nodeGroups := []*NodeGroup{}
	if err := db.Preload("Nodes").Select("*").Find(&nodeGroups).Error; err != nil {
		return err
	}
	designContent["node_groups"] = nodeGroups
	return nil
}
func deleteNodeGroupsFromDesign(db *gorm.DB) error {
	if err := db.Exec("delete from node_groups;").Error; err != nil {
		return err
	}
	if err := db.Exec("delete from node_group_association;").Error; err != nil {
		return err
	}
	return nil
}

func loadNodeGroupsFromDesign(db *gorm.DB, data interface{}) error {
	container := []*NodeGroup{}
	design := data.(*Design)
	if value, exists := design.Content["node_groups"]; exists {
		if err := mapstruct.MapToStruct(value.([]interface{}), &container); err != nil {
			return err
		}
		original := deepcopy.Copy(container).([]*NodeGroup)
		for _, nodeGroup := range container {
			nodeGroup.Nodes = nil
			if err := db.Create(nodeGroup).Error; err != nil {
				return err
			}
		}
		nodeGroups := original
		for _, nodeGroup := range nodeGroups {
			for _, node := range nodeGroup.Nodes {
				relatedNode := &Node{}
				if err := db.Preload("NodeGroups").First(&relatedNode, node.ID).Association("NodeGroups").Append(nodeGroup).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func init() {
	extension.RegisterModelType(reflect.TypeOf(NodeGroup{}))
	extension.RegisterDesignExtractor(extraceNodeGroupsFromDesign)
	extension.RegisterDesignDeleter(deleteNodeGroupsFromDesign)
	extension.RegisterDesignLoader(loadNodeGroupsFromDesign)
}

var NodeGroupModel = &NodeGroup{}
