package extension

var registeredProgramInformation ProgramInformation

// ProgramInformation is the interface what handles the build information like build time, version, and commit hash
// * GetBuildTime returns its build time
// * GetBranch returns branch
// * GetVersion returns version tag
// * GetCommitHash returns commit hash
type ProgramInformation interface {
	GetBuildTime() string
	GetBranch() string
	GetVersion() string
	GetCommitHash() string
}

// RegisterProgramInformation registers the program information
func RegisterProgramInformation(programInformation ProgramInformation) {
	registeredProgramInformation = programInformation
}

// GetRegisteredProgramInformation returns the registered program information
func GetRegisteredProgramInformation() ProgramInformation {
	return registeredProgramInformation
}
