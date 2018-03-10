package buildtime

import "github.com/qb0C80aE/clay/extension"

type clayProgramInformation struct {
	buildTime                    string
	claySubModuleInformationList []*claySubModuleInformation
}

type claySubModuleInformation struct {
	name     string
	revision string
}

func (clayProgramInformation *clayProgramInformation) BuildTime() string {
	return clayProgramInformation.buildTime
}

func (clayProgramInformation *clayProgramInformation) SubModuleInformationList() []extension.SubModuleInformation {
	result := make([]extension.SubModuleInformation, len(clayProgramInformation.claySubModuleInformationList))
	for i, subModuleInformation := range clayProgramInformation.claySubModuleInformationList {
		result[i] = subModuleInformation
	}
	return result
}

func (claySubModuleInformation *claySubModuleInformation) Name() string {
	return claySubModuleInformation.name
}

func (claySubModuleInformation *claySubModuleInformation) Revision() string {
	return claySubModuleInformation.revision
}
