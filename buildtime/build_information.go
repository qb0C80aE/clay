package buildtime

import "github.com/qb0C80aE/clay/extensions"

func init() {
	var programInformation = &clayProgramInformation{
		buildTime: "20180710153023",
		claySubModuleInformationList: []*claySubModuleInformation{
			{
				name:     "clay",
				revision: "85c2443a97b029f785e0fa19302bee6ff306bdef",
			},
			{
				name:     "github.com/qb0C80aE/loam",
				revision: "c6458eef282948d3cd93e52e5e0b90b4028eb451",
			},
			{
				name:     "github.com/qb0C80aE/pottery",
				revision: "9d825494b545f2771d85cfac94f9a98e598a8f9b",
			},
		},
	}
	extensions.RegisterProgramInformation(programInformation)
}
