package cmd

import (
	"fmt"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/spf13/cobra"
	"net/url"
)

var clientDeleteCmd = &cobra.Command{
	Use:   "delete <url>",
	Short: "Sends a DELETE request",
	Long:  `Sends a DELETE request to a specified server.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("invalid argument")
		}
		if _, err := url.ParseRequestURI(args[0]); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		urlArg := args[0]
		response, err := sendRequest(extensions.LookUpMethodName(extensions.MethodDelete), urlArg, "", nil)
		if err != nil {
			return err
		}
		return outputClientResponse(cmd, response)
	},
}

func init() {
	clientCmd.AddCommand(clientDeleteCmd)
}
