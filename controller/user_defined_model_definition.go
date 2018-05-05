package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
	"github.com/qb0C80aE/clay/version"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"time"
)

type userDefinedModelDefinitionController struct {
	BaseController
}

func newUserDefinedModelRegistrationController() *userDefinedModelDefinitionController {
	return CreateController(&userDefinedModelDefinitionController{}, model.NewUserDefinedModelDefinition()).(*userDefinedModelDefinitionController)
}

func (receiver *userDefinedModelDefinitionController) GetResourceSingleURL() (string, error) {
	modelKey, err := receiver.model.GetModelKey(receiver.model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	resourceName, err := receiver.GetResourceName()
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s/:%s", resourceName, modelKey.KeyParameter), nil
}

func (receiver *userDefinedModelDefinitionController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
			extension.URLMulti:  receiver.GetMulti,
		},
		extension.MethodPost: {
			extension.URLMulti: receiver.Create,
		},
	}
	return routeMap
}

func (receiver *userDefinedModelDefinitionController) Create(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	resourceName, err := receiver.model.GetResourceName(receiver.model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	container, err := receiver.binder.Bind(c, resourceName)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)
	result, err := receiver.model.Create(receiver.model, db, c.Params, c.Request.URL.Query(), container)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	typeName := reflect.ValueOf(container).Elem().FieldByName("TypeName").Interface().(string)
	newModel, err := extension.GetAssociatedModelWithTypeName(typeName)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	apiPath := ""
	if newUserDefinedManyToManyAssociationModel, ok := newModel.(*model.UserDefinedManyToManyAssociationModel); ok {
		if newUserDefinedManyToManyAssociationModel.IsControllerEnabled() {
			newParentController := newUserDefinedManyToManyAssociationModelController(newModel)
			newChildLeftController := newUserDefinedManyToManyAssociationChildModelController(newModel, true)
			newChildRightController := newUserDefinedManyToManyAssociationChildModelController(newModel, false)
			extension.RegisterController(newParentController)
			extension.RegisterController(newChildLeftController)
			extension.RegisterController(newChildRightController)
			newControllerList := []extension.Controller{
				newParentController,
				newChildLeftController,
				newChildRightController,
			}
			if err := extension.SetupController(apiPath, newControllerList); err != nil {
				logging.Logger().Debug(err.Error())
				receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
				return
			}
		}
	} else if newUserDefinedModel, ok := newModel.(*model.UserDefinedModel); ok {
		if newUserDefinedModel.IsControllerEnabled() {
			newController := newUserDefinedModelController(newModel)
			extension.RegisterController(newController)
			newControllerList := []extension.Controller{newController}
			if err := extension.SetupController(apiPath, newControllerList); err != nil {
				logging.Logger().Debug(err.Error())
				receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
				return
			}
		}
	} else {
		logging.Logger().Debug("model is not an UserDefinedModel or an UserDefinedManyToManyAssociationModel")
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, errors.New("model is not an UserDefinedManyToManyAssociationModel"))
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	receiver.outputHandler.OutputCreate(c, http.StatusCreated, result)
}

func (receiver *userDefinedModelDefinitionController) registerUserDefinedModels() error {
	modelConfigFilePath := os.Getenv("CLAY_MODEL_CONFIG_FILE_PATH")

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

	if len(modelConfigFilePath) == 0 {
		return nil
	}

	modelConfigJSONData, err := ioutil.ReadFile(modelConfigFilePath)
	if err != nil {
		logging.Logger().Critical(err.Error())
		os.Exit(1)
	}

	modelConfig := []string{}
	if err := json.Unmarshal(modelConfigJSONData, &modelConfig); err != nil {
		logging.Logger().Critical(err.Error())
		os.Exit(1)
	}

	modelPath := path.Dir(modelConfigFilePath)

	for _, fileName := range modelConfig {
		filePath := filepath.Join(modelPath, fileName)
		jsonData, err := ioutil.ReadFile(filePath)
		if err != nil {
			logging.Logger().Critical(err.Error())
			os.Exit(1)
		}

		url := fmt.Sprintf("http://%s:%s/user_defined_model_definitions", host, port)
		request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			logging.Logger().Critical(err.Error())
			os.Exit(1)
		}
		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logging.Logger().Critical(err.Error())
			os.Exit(1)
		}

		if response.StatusCode != http.StatusCreated {
			logging.Logger().Critical(string(responseBody))
			logging.Logger().Criticalf("failed to load %s, it might not be a correct model file", filePath)
			os.Exit(1)
		}
	}

	return nil
}

// DoAfterRouterSetup execute initialization process after Router initialization
// It registers user defined models automatically at boot time
func (receiver *userDefinedModelDefinitionController) DoAfterRouterSetup(r *gin.Engine) error {
	go receiver.registerUserDefinedModels()
	return nil
}

func init() {
	extension.RegisterController(newUserDefinedModelRegistrationController())
	extension.RegisterInitializer(newUserDefinedModelRegistrationController())
}
