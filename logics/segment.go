package logics

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/models"
	"github.com/qb0C80aE/clay/utils/net"
	"sort"
)

type SegmentLogic struct {
}

func sortNodeIDs(nodeMap map[int]*models.Node) []int {
	var sortedNodeIDs []int
	for nodeID := range nodeMap {
		sortedNodeIDs = append(sortedNodeIDs, nodeID)
	}
	sort.Ints(sortedNodeIDs)
	return sortedNodeIDs
}

type sortablePorts []*models.Port

func (this sortablePorts) Len() int {
	return len(this)
}

func (this sortablePorts) Less(i, j int) bool {
	return this[i].ID < this[j].ID
}

func (this sortablePorts) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func getStartPort(nodeMap map[int]*models.Node, sortedNodeIDs []int, consumedPortMap map[int]*models.Port) (*models.Port, string) {
	for _, nodeID := range sortedNodeIDs {
		var ports sortablePorts = nodeMap[nodeID].Ports
		sort.Sort(ports)
		for _, port := range ports {
			if port.Ipv4Address.Valid {
				if _, consumed := consumedPortMap[port.ID]; consumed {
					continue
				}
				portIpv4AddressString := fmt.Sprintf("%s/%d", port.Ipv4Address.String, port.Ipv4Prefix.Int64)
				portIpv4Address, _ := netutil.ParseCIDR(portIpv4AddressString)
				portIpv4AddressString = fmt.Sprintf("%s/%d", netutil.String(netutil.Network(portIpv4Address)), port.Ipv4Prefix.Int64)
				return port, portIpv4AddressString
			}
		}
	}

	return nil, ""

}

func tracePort(nodeMap map[int]*models.Node, portMap map[int]*models.Port, consumedPortMap map[int]*models.Port, port *models.Port, segmentIpv4Address string) []*models.Port {
	segmentPorts := make([]*models.Port, 0)
	trace := true

	if _, consumed := consumedPortMap[port.ID]; consumed {
		return segmentPorts
	}

	if !port.Ipv4Address.Valid {
		consumedPortMap[port.ID] = port

	} else {
		portIpv4AddressString := fmt.Sprintf("%s/%d", port.Ipv4Address.String, port.Ipv4Prefix.Int64)
		portIpv4Address, _ := netutil.ParseCIDR(portIpv4AddressString)
		portIpv4AddressString = fmt.Sprintf("%s/%d", netutil.String(netutil.Network(portIpv4Address)), port.Ipv4Prefix.Int64)

		if segmentIpv4Address == portIpv4AddressString {
			segmentPorts = append(segmentPorts, port)
			consumedPortMap[port.ID] = port
		} else {
			trace = false
		}
	}

	if trace {
		if port.DestinationPortID.Valid {
			destinationPort := portMap[int(port.DestinationPortID.Int64)]
			segmentPorts = append(segmentPorts, tracePort(nodeMap, portMap, consumedPortMap, destinationPort, segmentIpv4Address)...)
		}

		node := nodeMap[port.NodeID]
		for _, attachedPort := range node.Ports {
			targetPort := portMap[attachedPort.ID]
			segmentPorts = append(segmentPorts, tracePort(nodeMap, portMap, consumedPortMap, targetPort, segmentIpv4Address)...)
		}
	}

	return segmentPorts

}

func GenerateSegments(nodeMap map[int]*models.Node, portMap map[int]*models.Port, consumedPortMap map[int]*models.Port) []*models.Segment {

	sortedNodeIDs := sortNodeIDs(nodeMap)

	segments := make([]*models.Segment, 0, 10)
	for {
		startPort, segmentIp4Address := getStartPort(nodeMap, sortedNodeIDs, consumedPortMap)
		if startPort == nil {
			break
		}
		startPort = portMap[startPort.ID]
		segment := &models.Segment{
			Ports: tracePort(nodeMap, portMap, consumedPortMap, startPort, segmentIp4Address),
			Cidr:  segmentIp4Address,
		}
		segments = append(segments, segment)
	}
	return segments
}

func (_ *SegmentLogic) GetSingle(db *gorm.DB, id string, queryFields string) (interface{}, error) {
	return nil, nil
}

func (this *SegmentLogic) GetMulti(db *gorm.DB, queryFields string) ([]interface{}, error) {

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
	if err := db.Preload("Node").Select(queryFields).Find(&ports).Error; err != nil {
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

	segments := GenerateSegments(nodeMap, portMap, consumedPortMap)

	result := make([]interface{}, len(segments))
	for i, data := range segments {
		result[i] = data
	}

	return result, nil

}

func (_ *SegmentLogic) Create(_ *gorm.DB, _ interface{}) (interface{}, error) {
	return nil, nil
}

func (_ *SegmentLogic) Update(_ *gorm.DB, _ string, _ interface{}) (interface{}, error) {
	return nil, nil
}

func (_ *SegmentLogic) Delete(_ *gorm.DB, _ string) error {
	return nil
}

func (_ *SegmentLogic) Patch(_ *gorm.DB, _ string, _ string) (interface{}, error) {
	return nil, nil
}

func (_ *SegmentLogic) Options(db *gorm.DB) error {
	return nil
}

func (this *SegmentLogic) ExtractFromDesign(db *gorm.DB, designContent map[string]interface{}) error {
	segments, err := this.GetMulti(db, "*")
	if err != nil {
		return err
	}
	designContent["segments"] = segments
	return nil
}

func (_ *SegmentLogic) DeleteFromDesign(db *gorm.DB) error {
	return nil
}

func (_ *SegmentLogic) LoadToDesign(db *gorm.DB, data interface{}) error {
	return nil
}

var SegmentLogicInstance = &SegmentLogic{}

func init() {
	extension.RegisterDesignAccessor(SegmentLogicInstance)
}
