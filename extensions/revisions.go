package extensions

// ProgramInformation is the interface what handles the build information like build time, and sub module revisions
// * BuildTime returns its build time
// * SubModuleInformationList returns its sub module information list
type ProgramInformation interface {
	BuildTime() string
	SubModuleInformationList() []SubModuleInformation
}

// SubModuleInformation is the interface what handles the build information of sub module revisions
// * Name returns its name
// * Revision returns its revision
type SubModuleInformation interface {
	Name() string
	Revision() string
}

var registeredProgramInformation ProgramInformation

// RegisterProgramInformation registers the program information
func RegisterProgramInformation(programInformation ProgramInformation) {
	registeredProgramInformation = programInformation
}

// RegisteredProgramInformation returns the registered program information
func RegisteredProgramInformation() ProgramInformation {
	return registeredProgramInformation
}
