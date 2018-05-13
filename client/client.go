package client

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

var httpClient = &http.Client{}

func createJSONBody(value string) (io.Reader, error) {
	if len(value) == 0 {
		return bytes.NewReader([]byte("{}")), nil
	}

	if value[0:1] == "@" {
		fileName := value[1:]
		file, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(data), nil
	}

	return bytes.NewReader([]byte(value)), nil
}

func createMultipartFormData(formFlagValue string) (string, io.Reader, error) {
	var bytesBuffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&bytesBuffer)
	formValues := strings.Split(formFlagValue, ",")
	for _, formValue := range formValues {
		splittedFormValue := strings.Split(formValue, "=")
		if len(splittedFormValue) != 2 {
			return "", nil, fmt.Errorf("invalid value %s", formValue)
		}
		key := splittedFormValue[0]
		value := splittedFormValue[1]

		if len(value) == 0 {
			formWriter, err := multipartWriter.CreateFormField(key)
			if err != nil {
				return "", nil, err
			}
			if _, err = formWriter.Write([]byte{}); err != nil {
				return "", nil, err
			}
			continue
		}

		if value[0:1] == "@" {
			fileName := value[1:]
			file, err := os.Open(fileName)
			if err != nil {
				return "", nil, err
			}
			defer file.Close()

			formWriter, err := multipartWriter.CreateFormFile(key, fileName)
			if err != nil {
				return "", nil, err
			}
			if _, err = io.Copy(formWriter, file); err != nil {
				return "", nil, err
			}
		} else {
			formWriter, err := multipartWriter.CreateFormField(key)
			if err != nil {
				return "", nil, err
			}
			if _, err = formWriter.Write([]byte(value)); err != nil {
				return "", nil, err
			}
		}
	}

	if err := multipartWriter.Close(); err != nil {
		return "", nil, err
	}

	return multipartWriter.FormDataContentType(), &bytesBuffer, nil
}

func sendRequest(method string, url string, contentType string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if len(contentType) > 0 {
		request.Header.Set("Content-Type", contentType)
	}
	return httpClient.Do(request)
}

func outputClientResponse(cmd *cobra.Command, response *http.Response) error {
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	outputFlag := cmd.Flag("output")
	if outputFlag.Changed {
		file, err := os.OpenFile(outputFlag.Value.String(), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := file.Truncate(0); err != nil {
			return err
		}
		if _, err = file.Write(content); err != nil {
			return err
		}
	} else {
		fmt.Print(string(content))
	}

	verboseFlag := cmd.Flag("verbose")
	if verboseFlag.Changed {
		fmt.Printf("status code is %d\n", response.StatusCode)
	}

	return nil
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Provides a simple HTTP/HTTPS client function",
	Long:  `Provides a simple HTTP/HTTPS client function to access clay's model store without any other HTTP clients.`,
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.PersistentFlags().StringP("output", "o", "", "Filename to output retrieved resource data")
	clientCmd.PersistentFlags().BoolP("verbose", "v", false, "Display status code")
}
