package buildtime

type clayProgramInformation struct {
	buildTime  string
	branch     string
	version    string
	commitHash string
}

func (clayProgramInformation *clayProgramInformation) GetBuildTime() string {
	return clayProgramInformation.buildTime
}

func (clayProgramInformation *clayProgramInformation) GetBranch() string {
	return clayProgramInformation.branch
}

func (clayProgramInformation *clayProgramInformation) GetVersion() string {
	return clayProgramInformation.version
}

func (clayProgramInformation *clayProgramInformation) GetCommitHash() string {
	return clayProgramInformation.commitHash
}
