package buildtime

import "github.com/qb0C80aE/clay/extension"

func init() {
	var programInformation = &clayProgramInformation{
		buildTime: "unknown",
		claySubModuleInformationList: []*claySubModuleInformation{
			{
				name:     "clay",
				revision: "built-manually",
				version:  "built-manually",
			},
		},
	}
	extension.RegisterProgramInformation(programInformation)
}
