package buildtime

import "github.com/qb0C80aE/clay/extensions"

func init() {
	var programInformation = &clayProgramInformation{
		buildTime: "20180121151801",
		claySubModuleInformationList: []*claySubModuleInformation{
			{
				name:     "clay",
				revision: "d3d2ce380ce4d23014e4045069d49285f2c948f4",
			},
			{
				name:     "github.com/qb0C80aE/loam",
				revision: "c366ccbccef307f6cd1d59313e8ba29a036911a6",
			},
			{
				name:     "github.com/qb0C80aE/pottery",
				revision: "95f480d4467cb6240b7b760b3000728a587a9608",
			},
		},
	}
	extensions.RegisterProgramInformation(programInformation)
}
