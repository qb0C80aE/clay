package models

type Segment struct {
	Ports []*Port `json:"ports"`
	Cidr  string  `json:"cidr"`
}

var SegmentModel = &Segment{}

func init() {
}
