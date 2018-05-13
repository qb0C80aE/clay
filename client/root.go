package client

import (
	"fmt"
	"os"

	_ "github.com/qb0C80aE/clay/runtime" // Import runtime package to register Clay runtime

	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/spf13/cobra"
	"strings"
)

func getEnvironmentalVariable(name string, defaultValue string) string {
	value := strings.Trim(os.Getenv(name), " ")
	if len(value) == 0 {
		return defaultValue
	}

	return value
}

var clayConfigFilePath = getEnvironmentalVariable("CLAY_CONFIG_FILE_PATH", "(default)")
var clayHost = getEnvironmentalVariable("CLAY_HOST", "(default)")
var clayPort = getEnvironmentalVariable("CLAY_PORT", "(default)")
var clayDBMode = getEnvironmentalVariable("CLAY_DB_MODE", "(default)")
var clayDBFilePath = getEnvironmentalVariable("CLAY_DB_FILE_PATH", "(default)")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "clay",
	Short: "An abstract system model store to automate various kind of operations",
	Long: `An abstract system model store to automate various kind of operations
By default, clay boots the system model store server.
If you want to know what API endpoints clay has, send a GET request to the path '/' of the clay server.`,
	Run: func(cmd *cobra.Command, args []string) {
		logging.Logger().Debugf(`Environmental Variables:
  CLAY_CONFIG_FILE_PATH   [Default: "./clay_config.json", Current: "%s"]
  CLAY_HOST               [Default: "localhost"         , Current: "%s"]
  CLAY_PORT               [Default: "8080"              , Current: "%s"]
  CLAY_DB_MODE            [Default: "file"              , Current: "%s"]
  CLAY_DB_FILE_PATH       [Default: "./clay.db"         , Current: "%s"]`,
			clayConfigFilePath,
			clayHost,
			clayPort,
			clayDBMode,
			clayDBFilePath,
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

	rootCmd.SetUsageTemplate(fmt.Sprintf(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}

Environmental Variables:
  CLAY_CONFIG_FILE_PATH   [Default: "./clay_config.json", Current: "%s"]
  CLAY_HOST               [Default: "localhost"         , Current: "%s"]
  CLAY_PORT               [Default: "8080"              , Current: "%s"]
  CLAY_DB_MODE            [Default: "file"              , Current: "%s"]
  CLAY_DB_FILE_PATH       [Default: "./clay.db"         , Current: "%s"]{{if gt (len .Aliases) 0}}

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
		clayConfigFilePath,
		clayHost,
		clayPort,
		clayDBMode,
		clayDBFilePath,
	),
	)
}
