package models

import (
	"database/sql"
	"github.com/qb0C80aE/clay/extension"
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

var NodeModel = &Node{}

func init() {
	extension.RegisterModelType(NodeModel)
}
