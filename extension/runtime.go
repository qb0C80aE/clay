package extension

import (
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

// RegisterRuntime registers the runtime
func RegisterRuntime(runtime Runtime) {
	registeredRuntime = runtime
}

// GetRegisteredRuntime returns the registered runtime
func GetRegisteredRuntime() Runtime {
	return registeredRuntime
}

// EnvironmentalVariableSet contains variable environments used in Clay
type EnvironmentalVariableSet struct {
	clayConfigFilePath string
	clayHost           string
	clayPort           string
	clayDBMode         string
	clayDBFilePath     string
}

func getEnvironmentalVariable(name string, defaultValue string) string {
	value := strings.Trim(os.Getenv(name), " ")
	if len(value) == 0 {
		return defaultValue
	}

	return value
}

var defaultEnvironmentVariableSet = EnvironmentalVariableSet{
	clayConfigFilePath: "./clay_config.json",
	clayHost:           "localhost",
	clayPort:           "8080",
	clayDBMode:         "file",
	clayDBFilePath:     "./clay.db",
}

var currentEnvironmentalVariableSet = EnvironmentalVariableSet{
	clayConfigFilePath: getEnvironmentalVariable("CLAY_CONFIG_FILE_PATH", defaultEnvironmentVariableSet.clayConfigFilePath),
	clayHost:           getEnvironmentalVariable("CLAY_HOST", defaultEnvironmentVariableSet.clayHost),
	clayPort:           getEnvironmentalVariable("CLAY_PORT", defaultEnvironmentVariableSet.clayPort),
	clayDBMode:         getEnvironmentalVariable("CLAY_DB_MODE", defaultEnvironmentVariableSet.clayDBMode),
	clayDBFilePath:     getEnvironmentalVariable("CLAY_DB_FILE_PATH", defaultEnvironmentVariableSet.clayDBFilePath),
}

// GetDefaultEnvironmentalVariableSet returns default environmental variable set
func GetDefaultEnvironmentalVariableSet() EnvironmentalVariableSet {
	return defaultEnvironmentVariableSet
}

// GetCurrentEnvironmentalVariableSet returns CLAY_CONFIG_FILE_PATH
func GetCurrentEnvironmentalVariableSet() EnvironmentalVariableSet {
	return currentEnvironmentalVariableSet
}

// GetClayConfigFilePath returns CLAY_CONFIG_FILE_PATH
func (receiver *EnvironmentalVariableSet) GetClayConfigFilePath() string {
	return receiver.clayConfigFilePath
}

// GetClayHost returns CLAY_HOST
func (receiver *EnvironmentalVariableSet) GetClayHost() string {
	return receiver.clayHost
}

// GetClayPort returns CLAY_POTY
func (receiver *EnvironmentalVariableSet) GetClayPort() string {
	return receiver.clayPort
}

// GetClayPortInt returns CLAY_POTY as interger
func (receiver *EnvironmentalVariableSet) GetClayPortInt() int {
	port, err := strconv.Atoi(receiver.clayPort)
	if err != nil {
		port, _ = strconv.Atoi(defaultEnvironmentVariableSet.clayPort)
	}

	return port
}

// GetClayDBMode returns CLAY_DB_MODE
func (receiver *EnvironmentalVariableSet) GetClayDBMode() string {
	return receiver.clayDBMode
}

// GetClayDBFilePath returns CLAY_DB_FILE_PATH
func (receiver *EnvironmentalVariableSet) GetClayDBFilePath() string {
	return receiver.clayDBFilePath
}
