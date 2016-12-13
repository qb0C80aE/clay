package models

type Design struct {
	Nodes      []*Node      `json:"nodes"`
	NodeTypes  []*NodeType  `json:"node_types"`
	NodePvs    []*NodePv    `json:"node_pvs"`
	Ports      []*Port      `json:"ports"`
	NodeGroups []*NodeGroup `json:"node_groups"`
}

var DesignModel = &Design{}
