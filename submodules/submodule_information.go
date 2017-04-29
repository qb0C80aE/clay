package submodules

// ProgramInformation represents the time when this program was build and module information
type ProgramInformation struct {
	BuildTime                string
	SubModuleInformationList []*SubModuleInformation
}

// SubModuleInformation represents name and revision of installed modules
type SubModuleInformation struct {
	Name     string
	Revision string
}

// BuildInformation returns the program information
func BuildInformation() *ProgramInformation {
	return programInformation
}
