package models

import (
	"database/sql"
	"github.com/qb0C80aE/clay/extension"
	"reflect"
)

type NodeGroup struct {
	ID     int            `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name   string         `json:"name" gorm:"not null;unique"`
	Remark sql.NullString `json:"remark"`
	Nodes  []*Node        `json:"nodes" gorm:"many2many:node_group_association;"`
}

func init() {
	extension.RegisterModelType(reflect.TypeOf(NodeGroup{}))
}

var NodeGroupModel = &NodeGroup{}
