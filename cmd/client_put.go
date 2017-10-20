package cmd

import (
	"fmt"

	"github.com/qb0C80aE/clay/extensions"
	"github.com/spf13/cobra"
	"net/url"
)

var clientPutCmd = &cobra.Command{
	Use:   "put <url>",
	Short: "Sends a POST request",
	Long:  `Sends a POST request to a specified server.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("invalid argument")
		}
		if _, err := url.ParseRequestURI(args[0]); err != nil {
			return err
		}
		contentTypeFlagValue := cmd.Flag("content-type").Value.String()
		jsonFlagValue := cmd.Flag("json").Value.String()
		formFlagValue := cmd.Flag("form").Value.String()
		if (contentTypeFlagValue != "json") && (contentTypeFlagValue != "form") {
			return fmt.Errorf("invalid content-type")
		}
		switch contentTypeFlagValue {
		case "json":
			if len(formFlagValue) > 0 {
				return fmt.Errorf("form must not be specified when content-type is json")
			}
		case "form":
			if len(jsonFlagValue) > 0 {
				return fmt.Errorf("json must not be specified when content-type is form")
			}
		default:
			return fmt.Errorf("invalid content-type")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		contentTypeFlagValue := cmd.Flag("content-type").Value.String()
		jsonFlagValue := cmd.Flag("json").Value.String()
		formFlagValue := cmd.Flag("form").Value.String()
		urlArg := args[0]

		switch contentTypeFlagValue {
		case "json":
			content, err := createJSONBody(jsonFlagValue)
			if err != nil {
				return err
			}

			response, err := sendRequest(extensions.LookUpMethodName(extensions.MethodPut),
				urlArg,
				"application/json",
				content)

			if err != nil {
				return err
			}

			if err := outputClientResponse(cmd, response); err != nil {
				return err
			}
		case "form":
			contentType, content, err := createMultipartFormData(formFlagValue)
			if err != nil {
				return err
			}

			response, err := sendRequest(extensions.LookUpMethodName(extensions.MethodPut),
				urlArg,
				contentType,
				content)

			if err != nil {
				return err
			}

			if err := outputClientResponse(cmd, response); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	clientCmd.AddCommand(clientPutCmd)
	clientPutCmd.Flags().StringP("content-type", "c", "json", "Content-Type value. 'json' and 'form' are valid")
	clientPutCmd.Flags().StringP("json", "j", "", `JSON data like '{"name": "clay", "remark": "tool"}'. Also @filename is available`)
	clientPutCmd.Flags().StringP("form", "f", "", "Form value like name=clay,remark=tool. Also key=@filename is available")
}
