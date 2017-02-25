package models

type Segment struct {
	Ports []*Port `json:"ports"`
	Cidr  string  `json:"cidr"`
}

func init() {
}

var SegmentModel = &Segment{}
