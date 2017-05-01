package revisions

import "github.com/qb0C80aE/clay/extensions"

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

func (clayProgramInformation *clayProgramInformation) SubModuleInformationList() []extensions.SubModuleInformation {
	result := make([]extensions.SubModuleInformation, len(clayProgramInformation.claySubModuleInformationList))
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
