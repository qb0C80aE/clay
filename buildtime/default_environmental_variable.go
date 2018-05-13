package buildtime

import (
	"os"
	"strconv"
	"strings"
)

type defaultEnvironmentalVariableSet struct {
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

func (receiver *defaultEnvironmentalVariableSet) GetClayConfigFilePath() string {
	return receiver.clayConfigFilePath
}

func (receiver *defaultEnvironmentalVariableSet) GetClayHost() string {
	return receiver.clayHost
}

func (receiver *defaultEnvironmentalVariableSet) GetClayPort() string {
	return receiver.clayPort
}

func (receiver *defaultEnvironmentalVariableSet) GetClayPortInt() int {
	port, err := strconv.Atoi(receiver.clayPort)
	if err != nil {
		port, _ = strconv.Atoi(receiver.clayPort)
	}

	return port
}

func (receiver *defaultEnvironmentalVariableSet) GetClayDBMode() string {
	return receiver.clayDBMode
}

func (receiver *defaultEnvironmentalVariableSet) GetClayDBFilePath() string {
	return receiver.clayDBFilePath
}

func (receiver *defaultEnvironmentalVariableSet) GetClayAssetMode() string {
	return receiver.clayAssetMode
}
