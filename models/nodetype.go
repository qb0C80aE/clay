package models

import (
	"github.com/qb0C80aE/clay/extension"
)

type NodeType struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"not null"`
}

var NodeTypeModel = &NodeType{}

func init() {
	extension.RegisterModelType(NodeTypeModel)
}
