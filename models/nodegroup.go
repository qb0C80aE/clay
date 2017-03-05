package models

import (
	"database/sql"
	"github.com/qb0C80aE/clay/extension"
)

type NodeGroup struct {
	ID     int            `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name   string         `json:"name" gorm:"not null;unique"`
	Remark sql.NullString `json:"remark"`
	Nodes  []*Node        `json:"nodes" gorm:"many2many:node_group_association;"`
}

var NodeGroupModel = &NodeGroup{}

func init() {
	extension.RegisterModelType(NodeGroupModel)
}
