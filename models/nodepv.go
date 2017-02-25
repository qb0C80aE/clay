package models

import (
	"github.com/qb0C80aE/clay/extension"
	"reflect"
)

type NodePv struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"not null"`
}

func init() {
	extension.RegisterModelType(reflect.TypeOf(NodePv{}))
}

var NodePvModel = &NodePv{}
