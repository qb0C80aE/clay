package models

import (
	"github.com/qb0C80aE/clay/extension"
	"reflect"
)

type NodeType struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"not null"`
}

func init() {
	extension.RegisterModelType(reflect.TypeOf(NodeType{}))
}

var NodeTypeModel = &NodeType{}
