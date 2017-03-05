package logics

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/models"
	"github.com/qb0C80aE/clay/utils/mapstruct"
	"strconv"
)

type NodePvLogic struct {
}

func (_ *NodePvLogic) GetSingle(db *gorm.DB, id string, queryFields string) (interface{}, error) {

	nodePv := &models.NodePv{}

	if err := db.Select(queryFields).First(nodePv, id).Error; err != nil {
		return nil, err
	}

	return nodePv, nil

}

func (_ *NodePvLogic) GetMulti(db *gorm.DB, queryFields string) ([]interface{}, error) {

	nodePvs := []*models.NodePv{}

	if err := db.Select(queryFields).Find(&nodePvs).Error; err != nil {
		return nil, err
	}

	result := make([]interface{}, len(nodePvs))
	for i, data := range nodePvs {
		result[i] = data
	}

	return result, nil

}

func (_ *NodePvLogic) Create(db *gorm.DB, data interface{}) (interface{}, error) {

	nodePv := data.(*models.NodePv)

	if err := db.Create(nodePv).Error; err != nil {
		return nil, err
	}

	return nodePv, nil

}

func (_ *NodePvLogic) Update(db *gorm.DB, id string, data interface{}) (interface{}, error) {

	nodePv := data.(*models.NodePv)
	nodePv.ID, _ = strconv.Atoi(id)

	if err := db.Save(&nodePv).Error; err != nil {
		return nil, err
	}

	return nodePv, nil

}

func (_ *NodePvLogic) Delete(db *gorm.DB, id string) error {

	nodePv := &models.NodePv{}

	if err := db.First(&nodePv, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&nodePv).Error; err != nil {
		return err
	}

	return nil

}

func (_ *NodePvLogic) Patch(_ *gorm.DB, _ string, _ string) (interface{}, error) {
	return nil, nil
}

func (_ *NodePvLogic) Options(db *gorm.DB) error {
	return nil
}

func (_ *NodePvLogic) ExtractFromDesign(db *gorm.DB, designContent map[string]interface{}) error {
	nodePvs := []*models.NodePv{}
	if err := db.Select("*").Find(&nodePvs).Error; err != nil {
		return err
	}
	designContent["node_pvs"] = nodePvs
	return nil
}

func (_ *NodePvLogic) DeleteFromDesign(db *gorm.DB) error {
	return db.Exec("delete from node_pvs;").Error
}

func (_ *NodePvLogic) LoadToDesign(db *gorm.DB, data interface{}) error {
	container := []*models.NodePv{}
	design := data.(*models.Design)
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

var NodePvLogicInstance = &NodePvLogic{}

func init() {
	extension.RegisterDesignAccessor(NodePvLogicInstance)
}
