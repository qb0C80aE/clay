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

	var environmentalVariableSet = &defaultEnvironmentalVariableSet{
		clayConfigFilePath: "./clay_config.json",
		clayHost:           "localhost",
		clayPort:           "8080",
		clayDBMode:         "file",
		clayDBFilePath:     "./clay.db",
		clayAssetMode:      "external",
	}
	extension.RegisterDefaultEnvironmentalVariableSet(environmentalVariableSet)
}
