package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/asset"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
)

type clayRuntimeInitializer struct {
}

type clayConfig struct {
	General                *clayConfigGeneral                 `json:"general"`
	UserDefinedModels      []*clayConfigUserDefinedModel      `json:"user_defined_models"`
	EphemeralTemplates     []*clayConfigEphemeralTemplate     `json:"ephemeral_templates"`
	EphemeralBinaryObjects []*clayConfigEphemeralBinaryObject `json:"ephemeral_binary_objects"`
	EphemeralScripts       []*clayConfigEphemeralScript       `json:"ephemeral_scripts"`
	URLAliases             []*model.URLAliasDefinition        `json:"url_aliases"`
}

type clayConfigGeneral struct {
	UserDefinedModelsDirectory      string `json:"user_defined_models_directory"`
	EphemeralTemplatesDirectory     string `json:"ephemeral_templates_directory"`
	EphemeralBinaryObjectsDirectory string `json:"ephemeral_binary_objects_directory"`
	EphemeralScriptsDirectory       string `json:"ephemeral_scripts_directory"`
}

type clayConfigUserDefinedModel struct {
	FileName string `json:"file_name"`
}

type clayConfigEphemeralTemplate struct {
	Name     string `json:"name"`
	FileName string `json:"file_name"`
}

type clayConfigEphemeralBinaryObject struct {
	Name     string `json:"name"`
	FileName string `json:"file_name"`
}

type clayConfigEphemeralScript struct {
	Name     string `json:"name"`
	FileName string `json:"file_name"`
}

func (receiver *clayRuntimeInitializer) readFile(filePath string) ([]byte, error) {
	environmentalVariableSet := extension.GetCurrentEnvironmentalVariableSet()
	switch environmentalVariableSet.GetClayAssetMode() {
	case "internal":
		// since go-assets uses path(always slashed path)
		filePath = filepath.ToSlash(filePath)
		file, err := asset.Assets.Open(filePath)
		if err != nil {
			logging.Logger().Debug(filePath)
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		defer file.Close()

		return ioutil.ReadAll(file)
	case "external":
		return ioutil.ReadFile(filePath)
	}

	logging.Logger().Debugf("invalid asset mode %s", environmentalVariableSet.GetClayAssetMode())
	return nil, fmt.Errorf("invalid asset mode %s", environmentalVariableSet.GetClayAssetMode())
}

func (receiver *clayRuntimeInitializer) copyFromFile(writer io.Writer, filePath string) error {
	environmentalVariableSet := extension.GetCurrentEnvironmentalVariableSet()
	switch environmentalVariableSet.GetClayAssetMode() {
	case "internal":
		// since go-assets uses path(always slashed path)
		filePath = filepath.ToSlash(filePath)
		file, err := asset.Assets.Open(filePath)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return err
		}
		defer file.Close()

		if _, err = io.Copy(writer, file); err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		return nil
	case "external":
		file, err := os.Open(filePath)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}
		defer file.Close()

		if _, err = io.Copy(writer, file); err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		return nil
	}

	logging.Logger().Debugf("invalid asset mode %s", environmentalVariableSet.GetClayAssetMode())
	return fmt.Errorf("invalid asset mode %s", environmentalVariableSet.GetClayAssetMode())
}

func (receiver *clayRuntimeInitializer) initialize() {
	environmentalVariableSet := extension.GetCurrentEnvironmentalVariableSet()
	configFilePath := environmentalVariableSet.GetClayConfigFilePath()

	if len(configFilePath) == 0 {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logging.Logger().Critical(err.Error())
			os.Exit(1)
		}

		configFilePath = filepath.Join(dir, "clay_config.json")
	}

	configJSONData, err := receiver.readFile(configFilePath)
	if err != nil {
		logging.Logger().Debugf("cloud not load %s, just boot up without initial configuration", configFilePath)
		return
	}

	config := &clayConfig{}
	if err := json.Unmarshal(configJSONData, config); err != nil {
		logging.Logger().Critical(err.Error())
		os.Exit(1)
	}

	if err := receiver.loadUserDefinedModels(config); err != nil {
		logging.Logger().Critical(err.Error())
		os.Exit(1)
	}

	if err := receiver.loadEphemeralTemplates(config); err != nil {
		logging.Logger().Critical(err.Error())
		os.Exit(1)
	}

	if err := receiver.loadEphemeralBinaryObjects(config); err != nil {
		logging.Logger().Critical(err.Error())
		os.Exit(1)
	}

	if err := receiver.loadEphemeralScripts(config); err != nil {
		logging.Logger().Critical(err.Error())
		os.Exit(1)
	}

	if err := receiver.loadURLAliases(config); err != nil {
		logging.Logger().Critical(err.Error())
		os.Exit(1)
	}
}

func (receiver *clayRuntimeInitializer) loadUserDefinedModels(config *clayConfig) error {
	for _, userDefinedModel := range config.UserDefinedModels {
		filePath := filepath.Join(config.General.UserDefinedModelsDirectory, userDefinedModel.FileName)
		jsonData, err := receiver.readFile(filePath)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		request, err := http.NewRequest(extension.LookUpMethodName(extension.MethodPost), "/user_defined_model_definitions", bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", "application/json")
		responseWriter := httptest.NewRecorder()

		extension.GetRegisteredEngine().ServeHTTP(responseWriter, request)
		response := responseWriter.Result()
		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		if response.StatusCode != http.StatusCreated {
			logging.Logger().Critical(fmt.Errorf("status code was %d", response.StatusCode))
			logging.Logger().Critical(string(responseBody))
			logging.Logger().Criticalf("failed to load %s, it might not be a correct model file, or loading dependencies sequence is wrong", filePath)
			return fmt.Errorf("status code was %d", response.StatusCode)
		}
	}
	return nil
}

func (receiver *clayRuntimeInitializer) loadEphemeralTemplates(config *clayConfig) error {
	for _, ephemeralTemplate := range config.EphemeralTemplates {
		filePath := filepath.Join(config.General.EphemeralTemplatesDirectory, ephemeralTemplate.FileName)
		data, err := receiver.readFile(filePath)
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

		request, err := http.NewRequest(extension.LookUpMethodName(extension.MethodPost), "/ephemeral_templates", bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", "application/json")
		responseWriter := httptest.NewRecorder()

		extension.GetRegisteredEngine().ServeHTTP(responseWriter, request)
		response := responseWriter.Result()
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

func (receiver *clayRuntimeInitializer) loadEphemeralBinaryObjects(config *clayConfig) error {
	for _, ephemeralBinaryObject := range config.EphemeralBinaryObjects {
		filePath := filepath.Join(config.General.EphemeralBinaryObjectsDirectory, ephemeralBinaryObject.FileName)

		var bytesBuffer bytes.Buffer
		multipartWriter := multipart.NewWriter(&bytesBuffer)

		nameWriter, err := multipartWriter.CreateFormField("name")
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		if _, err := nameWriter.Write([]byte(ephemeralBinaryObject.Name)); err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		contentWriter, err := multipartWriter.CreateFormFile("content", ephemeralBinaryObject.FileName)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		if err := receiver.copyFromFile(contentWriter, filePath); err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		if err := multipartWriter.Close(); err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		request, err := http.NewRequest(extension.LookUpMethodName(extension.MethodPost), "/ephemeral_binary_objects", &bytesBuffer)
		request.Header.Set("Content-Type", multipartWriter.FormDataContentType())
		responseWriter := httptest.NewRecorder()

		extension.GetRegisteredEngine().ServeHTTP(responseWriter, request)
		response := responseWriter.Result()
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

func (receiver *clayRuntimeInitializer) loadEphemeralScripts(config *clayConfig) error {
	for _, ephemeralScript := range config.EphemeralScripts {
		filePath := filepath.Join(config.General.EphemeralScriptsDirectory, ephemeralScript.FileName)
		data, err := receiver.readFile(filePath)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		ephemeralScriptModel := model.NewEphemeralScript()
		ephemeralScriptModel.Name = ephemeralScript.Name
		ephemeralScriptModel.ScriptContent = string(data)

		jsonData, err := json.Marshal(ephemeralScriptModel)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		request, err := http.NewRequest(extension.LookUpMethodName(extension.MethodPost), "/ephemeral_scripts", bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", "application/json")
		responseWriter := httptest.NewRecorder()

		extension.GetRegisteredEngine().ServeHTTP(responseWriter, request)
		response := responseWriter.Result()
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

func (receiver *clayRuntimeInitializer) loadURLAliases(config *clayConfig) error {
	for _, urlAliasDefinition := range config.URLAliases {
		jsonData, err := json.Marshal(urlAliasDefinition)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		request, err := http.NewRequest(extension.LookUpMethodName(extension.MethodPost), "/url_alias_definitions", bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", "application/json")
		responseWriter := httptest.NewRecorder()

		extension.GetRegisteredEngine().ServeHTTP(responseWriter, request)
		response := responseWriter.Result()
		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logging.Logger().Critical(err.Error())
			return err
		}

		if response.StatusCode != http.StatusCreated {
			logging.Logger().Critical(fmt.Errorf("invalid url alias definition of %s", urlAliasDefinition.Name))
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
