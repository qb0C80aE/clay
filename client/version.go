package client

import (
	"fmt"
	"github.com/qb0C80aE/clay/extension"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows version",
	Long:  `Shows clay and all submodule versions.`,
	Run: func(cmd *cobra.Command, args []string) {
		programInformation := extension.GetRegisteredProgramInformation()
		fmt.Printf(
			"Clay %s\nbranch: %s\ncommit hash: %s\nbuild time: %s\n",
			programInformation.GetVersion(),
			programInformation.GetBranch(),
			programInformation.GetCommitHash(),
			programInformation.GetBuildTime(),
		)
		return
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
