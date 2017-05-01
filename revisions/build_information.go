package revisions

import "github.com/qb0C80aE/clay/extensions"

func init() {
	var programInformation = &clayProgramInformation{
		buildTime: "unknown",
		claySubModuleInformationList: []*claySubModuleInformation{
			{
				name:     "clay",
				revision: "built-manually",
			},
		},
	}
	extensions.RegisterProgramInformation(programInformation)
}
