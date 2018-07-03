package extension

import (
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"strings"
)

// Runtime is the interface what handles the runtime environment
// * Run boots server
type Runtime interface {
	Run()
}

var registeredRuntime Runtime
var registeredEngine *gin.Engine

// RegisterRuntime registers the runtime
func RegisterRuntime(runtime Runtime) {
	registeredRuntime = runtime
}

// GetRegisteredRuntime returns the registered runtime
func GetRegisteredRuntime() Runtime {
	return registeredRuntime
}

// RegisterEngine registers the engine
func RegisterEngine(engine *gin.Engine) {
	registeredEngine = engine
}

// GetRegisteredEngine returns the registered engine
func GetRegisteredEngine() *gin.Engine {
	return registeredEngine
}

// EnvironmentalVariableSet is the interface what handles variable environments used in Clay
// * GetClayConfigFilePath returns CLAY_CONFIG_FILE_PATH
// * GetClayHost returns CLAY_HOST
// * GetClayPort returns CLAY_PORT
// * GetClayPortInt returns CLAY_PORT as interger
// * GetClayDBMode returns CLAY_DB_MODE
// * GetClayDBFilePath returns CLAY_DB_FILE_PATH
// * GetClayAssetMode returns CLAY_ASSET_MODE
type EnvironmentalVariableSet interface {
	GetClayConfigFilePath() string
	GetClayHost() string
	GetClayPort() string
	GetClayPortInt() int
	GetClayDBMode() string
	GetClayDBFilePath() string
	GetClayAssetMode() string
}

type environmentalVariableSet struct {
	clayConfigFilePath string
	clayHost           string
	clayPort           string
	clayDBMode         string
	clayDBFilePath     string
	clayAssetMode      string
}

func getEnvironmentalVariable(name string, defaultValue string) string {
	value := strings.Trim(os.Getenv(name), " ")
	if len(value) == 0 {
		return defaultValue
	}

	return value
}

var defaultEnvironmentalVariableSet EnvironmentalVariableSet
var currentEnvironmentalVariableSet EnvironmentalVariableSet

// RegisterDefaultEnvironmentalVariableSet registers default environmental variable set, and create current environmental variable set
func RegisterDefaultEnvironmentalVariableSet(newEnvironmentalVariableSet EnvironmentalVariableSet) {
	defaultEnvironmentalVariableSet = newEnvironmentalVariableSet
	currentEnvironmentalVariableSet = &environmentalVariableSet{
		clayConfigFilePath: getEnvironmentalVariable("CLAY_CONFIG_FILE_PATH", defaultEnvironmentalVariableSet.GetClayConfigFilePath()),
		clayHost:           getEnvironmentalVariable("CLAY_HOST", defaultEnvironmentalVariableSet.GetClayHost()),
		clayPort:           getEnvironmentalVariable("CLAY_PORT", defaultEnvironmentalVariableSet.GetClayPort()),
		clayDBMode:         getEnvironmentalVariable("CLAY_DB_MODE", defaultEnvironmentalVariableSet.GetClayDBMode()),
		clayDBFilePath:     getEnvironmentalVariable("CLAY_DB_FILE_PATH", defaultEnvironmentalVariableSet.GetClayDBFilePath()),
		clayAssetMode:      getEnvironmentalVariable("CLAY_ASSET_MODE", defaultEnvironmentalVariableSet.GetClayAssetMode()),
	}
}

// GetDefaultEnvironmentalVariableSet returns default environmental variable set
func GetDefaultEnvironmentalVariableSet() EnvironmentalVariableSet {
	return defaultEnvironmentalVariableSet
}

// GetCurrentEnvironmentalVariableSet returns CLAY_CONFIG_FILE_PATH
func GetCurrentEnvironmentalVariableSet() EnvironmentalVariableSet {
	return currentEnvironmentalVariableSet
}

func (receiver *environmentalVariableSet) GetClayConfigFilePath() string {
	return receiver.clayConfigFilePath
}

func (receiver *environmentalVariableSet) GetClayHost() string {
	return receiver.clayHost
}

func (receiver *environmentalVariableSet) GetClayPort() string {
	return receiver.clayPort
}

func (receiver *environmentalVariableSet) GetClayPortInt() int {
	port, err := strconv.Atoi(receiver.clayPort)
	if err != nil {
		port, _ = strconv.Atoi(receiver.clayPort)
	}

	return port
}

func (receiver *environmentalVariableSet) GetClayDBMode() string {
	return receiver.clayDBMode
}

func (receiver *environmentalVariableSet) GetClayDBFilePath() string {
	return receiver.clayDBFilePath
}

func (receiver *environmentalVariableSet) GetClayAssetMode() string {
	return receiver.clayAssetMode
}
