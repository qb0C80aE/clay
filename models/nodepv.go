package models

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/utils/mapstruct"
	"reflect"
)

type NodePv struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"not null"`
}

func extractNodePvsFromDesign(db *gorm.DB, designContent map[string]interface{}) error {
	nodePvs := []*NodePv{}
	if err := db.Select("*").Find(&nodePvs).Error; err != nil {
		return err
	}
	designContent["node_pvs"] = nodePvs
	return nil
}

func deleteNodePvsFromDesign(db *gorm.DB) error {
	return db.Exec("delete from node_pvs;").Error
}

func loadNodePvsFromDesign(db *gorm.DB, data interface{}) error {
	container := []*NodePv{}
	design := data.(*Design)
	if value, exists := design.Content["node_pvs"]; exists {
		if err := mapstruct.MapToStruct(value.([]interface{}), &container); err != nil {
			return err
		}
		for _, nodePv := range container {
			if err := db.Create(nodePv).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func init() {
	extension.RegisterModelType(reflect.TypeOf(NodePv{}))
	extension.RegisterDesignExtractor(extractNodePvsFromDesign)
	extension.RegisterDesignDeleter(deleteNodePvsFromDesign)
	extension.RegisterDesignLoader(loadNodePvsFromDesign)
}

var NodePvModel = &NodePv{}
