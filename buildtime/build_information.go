package buildtime

import "github.com/qb0C80aE/clay/extension"

func init() {
	var programInformation = &clayProgramInformation{
		buildTime:  "unknown",
		branch:     "unknown",
		version:    "unknown",
		commitHash: "unknown",
	}
	extension.RegisterProgramInformation(programInformation)
}
