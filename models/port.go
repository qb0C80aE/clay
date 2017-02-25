package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/mohae/deepcopy"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/utils/mapstruct"
	"reflect"
)

type Port struct {
	ID                int            `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	NodeID            int            `json:"node_id" gorm:"index" sql:"type:integer references nodes(id) on delete cascade"`
	Node              *Node          `json:"node"`
	Name              string         `json:"name" gorm:"not null"`
	DestinationPortID sql.NullInt64  `json:"destination_port_id" sql:"type:integer references ports(id) on delete set null"`
	DestinationPort   *Port          `json:"destination_port"`
	MacAddress        sql.NullString `json:"mac_address"`
	Ipv4Address       sql.NullString `json:"ipv4_address"`
	Ipv4Prefix        sql.NullInt64  `json:"ipv4_prefix"`
	Remark            sql.NullString `json:"remark"`
}

func extracePortsFromDesign(db *gorm.DB, designContent map[string]interface{}) error {
	ports := []*Port{}
	if err := db.Select("*").Find(&ports).Error; err != nil {
		return err
	}
	designContent["ports"] = ports
	return nil
}

func deletePortsFromDesign(db *gorm.DB) error {
	if err := db.Exec("delete from ports;").Error; err != nil {
		return err
	}
	return nil
}

func loadPortsFromDesign(db *gorm.DB, data interface{}) error {
	container := []*Port{}
	design := data.(*Design)
	if value, exists := design.Content["ports"]; exists {
		if err := mapstruct.MapToStruct(value.([]interface{}), &container); err != nil {
			return err
		}
		original := deepcopy.Copy(container).([]*Port)
		for _, port := range container {
			port.DestinationPortID = sql.NullInt64{Int64: 0, Valid: false}
			if err := db.Create(port).Error; err != nil {
				return err
			}
		}
		ports := original
		for _, port := range ports {
			if err := db.Save(port).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func init() {
	extension.RegisterModelType(reflect.TypeOf(Port{}))
	extension.RegisterDesignExtractor(extracePortsFromDesign)
	extension.RegisterDesignDeleter(deletePortsFromDesign)
	extension.RegisterDesignLoader(loadPortsFromDesign)
}

var PortModel = &Port{}
