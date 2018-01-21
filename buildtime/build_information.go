package buildtime

import "github.com/qb0C80aE/clay/extensions"

func init() {
	var programInformation = &clayProgramInformation{
		buildTime: "20180110155656",
		claySubModuleInformationList: []*claySubModuleInformation{
			{
				name:     "clay",
				revision: "6dcaf03f94f46822d97a9899bb6d20408879d225",
			},
			{
				name:     "github.com/qb0C80aE/loam",
				revision: "b03da4d83378bdb0204b06b2b676f3517edd9b3c",
			},
			{
				name:     "github.com/qb0C80aE/pottery",
				revision: "36dff00f9847c342882b4d7d38eae9ee99b4cac9",
			},
		},
	}
	extensions.RegisterProgramInformation(programInformation)
}
