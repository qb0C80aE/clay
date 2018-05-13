package client

import (
	"fmt"
	"os"

	_ "github.com/qb0C80aE/clay/runtime" // Import runtime package to register Clay runtime

	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "clay",
	Short: "An abstract system model store to automate various kind of operations",
	Long: `An abstract system model store to automate various kind of operations
By default, clay boots the system model store server.
If you want to know what API endpoints clay has, send a GET request to the path '/' of the clay server.`,
	Run: func(cmd *cobra.Command, args []string) {
		defaultEnvironmentalVariableSet := extension.GetDefaultEnvironmentalVariableSet()
		environmentalVariableSet := extension.GetCurrentEnvironmentalVariableSet()

		logging.Logger().Debugf(`Environmental Variables:
  CLAY_CONFIG_FILE_PATH   [Default: %-21v Current: %v]
  CLAY_HOST               [Default: %-21v Current: %v]
  CLAY_PORT               [Default: %-21v Current: %v]
  CLAY_DB_MODE            [Default: %-21v Current: %v]
  CLAY_DB_FILE_PATH       [Default: %-21v Current: %v]`,
			fmt.Sprintf(`"%s",`, defaultEnvironmentalVariableSet.GetClayConfigFilePath()),
			fmt.Sprintf(`"%s"`, environmentalVariableSet.GetClayConfigFilePath()),
			fmt.Sprintf(`"%s",`, defaultEnvironmentalVariableSet.GetClayHost()),
			fmt.Sprintf(`"%s"`, environmentalVariableSet.GetClayHost()),
			fmt.Sprintf(`"%s",`, defaultEnvironmentalVariableSet.GetClayPort()),
			fmt.Sprintf(`"%s"`, environmentalVariableSet.GetClayPort()),
			fmt.Sprintf(`"%s",`, defaultEnvironmentalVariableSet.GetClayDBMode()),
			fmt.Sprintf(`"%s"`, environmentalVariableSet.GetClayDBMode()),
			fmt.Sprintf(`"%s",`, defaultEnvironmentalVariableSet.GetClayDBFilePath()),
			fmt.Sprintf(`"%s"`, environmentalVariableSet.GetClayDBFilePath()),
		)
		extension.GetRegisteredRuntime().Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()

	defaultEnvironmentalVariableSet := extension.GetDefaultEnvironmentalVariableSet()
	environmentalVariableSet := extension.GetCurrentEnvironmentalVariableSet()

	rootCmd.SetUsageTemplate(fmt.Sprintf(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}

Environmental Variables:
  CLAY_CONFIG_FILE_PATH   [Default: %-21v Current: %v]
  CLAY_HOST               [Default: %-21v Current: %v]
  CLAY_PORT               [Default: %-21v Current: %v]
  CLAY_DB_MODE            [Default: %-21v Current: %v]
  CLAY_DB_FILE_PATH       [Default: %-21v Current: %v]{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`,
		fmt.Sprintf(`"%s",`, defaultEnvironmentalVariableSet.GetClayConfigFilePath()),
		fmt.Sprintf(`"%s"`, environmentalVariableSet.GetClayConfigFilePath()),
		fmt.Sprintf(`"%s",`, defaultEnvironmentalVariableSet.GetClayHost()),
		fmt.Sprintf(`"%s"`, environmentalVariableSet.GetClayHost()),
		fmt.Sprintf(`"%s",`, defaultEnvironmentalVariableSet.GetClayPort()),
		fmt.Sprintf(`"%s"`, environmentalVariableSet.GetClayPort()),
		fmt.Sprintf(`"%s",`, defaultEnvironmentalVariableSet.GetClayDBMode()),
		fmt.Sprintf(`"%s"`, environmentalVariableSet.GetClayDBMode()),
		fmt.Sprintf(`"%s",`, defaultEnvironmentalVariableSet.GetClayDBFilePath()),
		fmt.Sprintf(`"%s"`, environmentalVariableSet.GetClayDBFilePath()),
	),
	)
}
