package models

import (
	"github.com/qb0C80aE/clay/extension"
)

type NodePv struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"not null"`
}

var NodePvModel = &NodePv{}

func init() {
	extension.RegisterModelType(NodePvModel)
}
