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
		fmt.Printf("Clay build-%s\n", programInformation.BuildTime())
		subModuleInformationList := programInformation.SubModuleInformationList()
		for _, subModuleInformation := range subModuleInformationList {
			fmt.Printf("  module %s\n    revision: %s\n", subModuleInformation.Name(), subModuleInformation.Revision())
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
