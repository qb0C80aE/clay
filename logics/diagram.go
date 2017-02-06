package logics

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/models"
)

var physicalNodeIconPaths = map[int]string{
	1: "/ui/images/diagram/l2switch.png",
	2: "/ui/images/diagram/l3switch.png",
	3: "/ui/images/diagram/firewall.png",
	4: "/ui/images/diagram/router.png",
	5: "/ui/images/diagram/loadbalancer.png",
	6: "/ui/images/diagram/server.png",
	7: "/ui/images/diagram/network.png",
}

var virtualNodeIconPaths = map[int]string{
	1: "/ui/images/diagram/l2switch_v.png",
	2: "/ui/images/diagram/l3switch_v.png",
	3: "/ui/images/diagram/firewall_v.png",
	4: "/ui/images/diagram/router_v.png",
	5: "/ui/images/diagram/loadbalancer_v.png",
	6: "/ui/images/diagram/server_v.png",
	7: "/ui/images/diagram/network.png",
}

const segmentIconPath = "/ui/images/diagram/segment.png"

func GetPhysicalDiagram(db *gorm.DB, _ string, queryFields string) (interface{}, error) {

	diagram := &models.Diagram{}

	nodes := []*models.Node{}
	if err := db.Preload("Ports").Select(queryFields).Find(&nodes).Error; err != nil {
		return nil, err
	}

	nodeMap := make(map[int]*models.Node)
	for _, node := range nodes {
		nodeMap[node.ID] = node
	}

	ports := []*models.Port{}
	if err := db.Select(queryFields).Find(&ports).Error; err != nil {
		return nil, err
	}

	portMap := make(map[int]*models.Port)
	for _, port := range ports {
		portMap[port.ID] = port
	}

	for _, node := range nodes {
		var iconPathMap map[int]string = nil
		if node.NodePvID == 1 {
			iconPathMap = physicalNodeIconPaths
		} else {
			iconPathMap = virtualNodeIconPaths
		}
		diagramNode := &models.DiagramNode{
			Name: node.Name,
			Icon: iconPathMap[node.NodeTypeID],
		}
		diagram.Nodes = append(diagram.Nodes, diagramNode)
	}

	registerdPortMap := make(map[int]int)
	for _, port := range ports {
		_, exists := registerdPortMap[int(port.DestinationPortID.Int64)]
		if (port.DestinationPortID.Valid) && (!exists) {
			sourceNode := nodeMap[port.NodeID]
			destinationPort := portMap[int(port.DestinationPortID.Int64)]
			destinationNode := nodeMap[destinationPort.NodeID]

			diagramInterface := &models.DiagramInterface{
				Source: port.Name,
				Target: destinationPort.Name,
			}
			diagramMeta := &models.DiagramMeta{
				Interface: diagramInterface,
			}
			diagramLink := &models.DiagramLink{
				Source: sourceNode.Name,
				Target: destinationNode.Name,
				Meta:   diagramMeta,
			}

			diagram.Links = append(diagram.Links, diagramLink)

			registerdPortMap[port.ID] = port.ID
		}
	}

	return diagram, nil

}

func GetLogicalDiagram(db *gorm.DB, _ string, queryFields string) (interface{}, error) {

	nodePvs := []*models.NodePv{}
	if err := db.Select(queryFields).Find(&nodePvs).Error; err != nil {
		return nil, err
	}

	nodeTypes := []*models.NodeType{}
	if err := db.Select(queryFields).Find(&nodeTypes).Error; err != nil {
		return nil, err
	}

	nodes := []*models.Node{}
	if err := db.Preload("Ports").Select(queryFields).Find(&nodes).Error; err != nil {
		return nil, err
	}

	ports := []*models.Port{}
	if err := db.Select(queryFields).Find(&ports).Error; err != nil {
		return nil, err
	}

	nodeMap := make(map[int]*models.Node)
	portMap := make(map[int]*models.Port)
	consumedPortMap := make(map[int]*models.Port)

	for _, node := range nodes {
		nodeMap[node.ID] = node
	}
	for _, port := range ports {
		portMap[port.ID] = port
	}

	segments := createSegments(nodeMap, portMap, consumedPortMap)

	diagram := &models.Diagram{}

	for _, node := range nodes {
		if node.NodeTypeID != 1 {
			var iconPathMap map[int]string = nil
			if node.NodePvID == 1 {
				iconPathMap = physicalNodeIconPaths
			} else {
				iconPathMap = virtualNodeIconPaths
			}
			diagramNode := &models.DiagramNode{
				node.Name,
				iconPathMap[node.NodeTypeID],
			}
			diagram.Nodes = append(diagram.Nodes, diagramNode)
		}
	}

	for i, segment := range segments {

		diagramNode := &models.DiagramNode{
			fmt.Sprintf("[%d]%s", i, segment.Cidr),
			segmentIconPath,
		}
		diagram.Nodes = append(diagram.Nodes, diagramNode)

		for _, port := range segment.Ports {
			diagramInterface := &models.DiagramInterface{
				Source: "",
				Target: fmt.Sprintf("%s[%s](%s/%d)",
					port.Name,
					port.MacAddress.String,
					port.Ipv4Address.String,
					port.Ipv4Prefix.Int64,
				),
			}
			diagramMeta := &models.DiagramMeta{
				Interface: diagramInterface,
			}
			diagramLink := &models.DiagramLink{
				Source: fmt.Sprintf("[%d]%s", i, segment.Cidr),
				Target: nodeMap[port.NodeID].Name,
				Meta:   diagramMeta,
			}
			diagram.Links = append(diagram.Links, diagramLink)
		}

	}

	return diagram, nil

}
