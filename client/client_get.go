package client

import (
	"fmt"
	"github.com/qb0C80aE/clay/extension"
	"github.com/spf13/cobra"
	"net/url"
)

var clientGetCmd = &cobra.Command{
	Use:   "get <url>",
	Short: "Sends a GET request",
	Long:  `Sends a GET request to a specified server.`,
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
		response, err := sendRequest(extension.LookUpMethodName(extension.MethodGet), urlArg, "", nil)
		if err != nil {
			return err
		}
		return outputClientResponse(cmd, response)
	},
}

func init() {
	clientCmd.AddCommand(clientGetCmd)
}
