package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type clayRuntimeInitializer struct {
}

type clayConfig struct {
	General            *clayConfigGeneral             `json:"general"`
	UserDefinedModels  []*clayConfigUserDefinedModel  `json:"user_defined_models"`
	EphemeralTemplates []*clayConfigEphemeralTemplate `json:"ephemeral_templates"`
	URLAliases         []*clayConfigURLAlias          `json:"url_aliases"`
}

type clayConfigGeneral struct {
	UserDefinedModelsDirectory  string `json:"user_defined_models_directory"`
	EphemeralTemplatesDirectory string `json:"ephemeral_templates_directory"`
}

type clayConfigUserDefinedModel struct {
	FileName string `json:"file_name"`
}

type clayConfigEphemeralTemplate struct {
	Name     string `json:"name"`
	FileName string `json:"file_name"`
}

type clayConfigURLAlias struct {
	Name  string `json:"name"`
	From  string `json:"from"`
	To    string `json:"to"`
	Query string `json:"query"`
}

func (receiver *clayRuntimeInitializer) initialize() {
	configFilePath := os.Getenv("CLAY_CONFIG_FILE_PATH")

	host := "localhost"
	port := "8080"

	if h := os.Getenv("CLAY_HOST"); h != "" {
		host = h
	}

	if p := os.Getenv("CLAY_PORT"); p != "" {
		if _, err := strconv.Atoi(p); err == nil {
			port = p
		}
	}

	maxRetry := 10
	for i := 1; i <= maxRetry; i++ {
		url := fmt.Sprintf("http://%s:%s", host, port)
		request, _ := http.NewRequest("GET", url, nil)
		client := &http.Client{}
		response, err := client.Do(request)

		if err == nil {
			defer response.Body.Close()
			if response.StatusCode == http.StatusOK {
				break
			} else {
				os.Exit(1)
			}
		}

		if i == maxRetry {
			logging.Logger().Critical("failed to connect server at model loading")
			os.Exit(1)
		}

		time.Sleep(time.Second)
	}

	if len(configFilePath) == 0 {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logging.Logger().Critical(err.Error())
			os.Exit(1)
		}

		configFilePath = path.Join(dir, "clay_config.json")
	}

	configJSONData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		logging.Logger().Debugf("cloud not load %s, just boot up without initial configuration", configFilePath)
		return
	}

	config := &clayConfig{}
	if err := json.Unmarshal(configJSONData, config); err != nil {
		logging.Logger().Critical(err.Error())
		os.Exit(1)
	}

	if err := receiver.loadUserDefinedModels(config, host, port); err != nil {
		logging.Logger().Critical(err.Error())
		os.Exit(1)
	}

	if err := receiver.loadEphemeralTemplates(config, host, port); err != nil {
		logging.Logger().Critical(err.Error())
		os.Exit(1)
	}

	if err := receiver.loadURLAliases(config, host, port); err != nil {
		logging.Logger().Critical(err.Error())
		os.Exit(1)
	}
}

func (receiver *clayRuntimeInitializer) loadUserDefinedModels(config *clayConfig, host string, port string) error {
	for _, userDefinedModel := range config.UserDefinedModels {
		filePath := filepath.Join(config.General.UserDefinedModelsDirectory, userDefinedModel.FileName)
		jsonData, err := ioutil.ReadFile(filePath)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		url := fmt.Sprintf("http://%s:%s/user_defined_model_definitions", host, port)
		request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}
		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		if response.StatusCode != http.StatusCreated {
			logging.Logger().Critical(fmt.Errorf("status code was %d", response.StatusCode))
			logging.Logger().Critical(string(responseBody))
			logging.Logger().Criticalf("failed to load %s, it might not be a correct model file", filePath)
			return fmt.Errorf("status code was %d", response.StatusCode)
		}
	}
	return nil
}

func (receiver *clayRuntimeInitializer) loadEphemeralTemplates(config *clayConfig, host string, port string) error {
	for _, ephemeralTemplate := range config.EphemeralTemplates {
		filePath := filepath.Join(config.General.EphemeralTemplatesDirectory, ephemeralTemplate.FileName)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		ephemeralTemplateModel := model.NewEphemeralTemplate()
		ephemeralTemplateModel.Name = ephemeralTemplate.Name
		ephemeralTemplateModel.TemplateContent = string(data)

		jsonData, err := json.Marshal(ephemeralTemplateModel)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		url := fmt.Sprintf("http://%s:%s/ephemeral_templates", host, port)
		request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}
		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		if response.StatusCode != http.StatusCreated {
			logging.Logger().Critical(fmt.Errorf("status code was %d", response.StatusCode))
			logging.Logger().Critical(string(responseBody))
			return fmt.Errorf("status code was %d", response.StatusCode)
		}
	}
	return nil
}

func (receiver *clayRuntimeInitializer) loadURLAliases(config *clayConfig, host string, port string) error {
	for _, urlAliasDefinition := range config.URLAliases {
		urlAliasDefinitionModel := model.NewURLAliasDefinition()
		urlAliasDefinitionModel.Name = urlAliasDefinition.Name
		urlAliasDefinitionModel.From = urlAliasDefinition.From
		urlAliasDefinitionModel.To = urlAliasDefinition.To
		urlAliasDefinitionModel.Query = urlAliasDefinition.Query

		jsonData, err := json.Marshal(urlAliasDefinitionModel)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		url := fmt.Sprintf("http://%s:%s/url_alias_definitions", host, port)
		request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}
		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		if response.StatusCode != http.StatusCreated {
			logging.Logger().Critical(fmt.Errorf("status code was %d", response.StatusCode))
			logging.Logger().Critical(string(responseBody))
			return fmt.Errorf("status code was %d", response.StatusCode)
		}
	}
	return nil
}

func (receiver *clayRuntimeInitializer) DoBeforeDBMigration(_ *gorm.DB) error {
	return nil
}

func (receiver *clayRuntimeInitializer) DoAfterDBMigration(_ *gorm.DB) error {
	return nil
}

func (receiver *clayRuntimeInitializer) DoBeforeRouterSetup(_ *gin.Engine) error {
	return nil
}

func (receiver *clayRuntimeInitializer) DoAfterRouterSetup(_ *gin.Engine) error {
	go receiver.initialize()
	return nil
}

func init() {
	runtimeInitializer := &clayRuntimeInitializer{}
	extension.RegisterInitializer(runtimeInitializer)
}
