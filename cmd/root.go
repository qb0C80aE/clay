package cmd

import (
	"fmt"
	"os"

	_ "github.com/qb0C80aE/clay/runtime" // Import runtime package to register Clay runtime

	"github.com/qb0C80aE/clay/extensions"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "clay",
	Short: "An abstract system model store to automate various kind of operations",
	Long: `An abstract system model store to automate various kind of operations
By default, clay boots the system model store server.
If you want to know what API endpoints clay has, send a GET request to the path '/' of the clay server.`,
	Run: func(cmd *cobra.Command, args []string) {
		extensions.RegisteredRuntime().Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
}
