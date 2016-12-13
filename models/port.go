package models

import "database/sql"

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

var PortModel = &Port{}
