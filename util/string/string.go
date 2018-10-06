package string

import (
	"fmt"
	"github.com/qb0C80aE/clay/logging"
	mapstructutilpkg "github.com/qb0C80aE/clay/util/mapstruct"
	"regexp"
	"strings"
)

var utility = &Utility{}

// Utility handles string operation
type Utility struct {
}

// GetUtility returns the instance of utility
func GetUtility() *Utility {
	return utility
}

// Sprintf apply sprintf to given strings
func (receiver *Utility) Sprintf(format string, parameters ...interface{}) interface{} {
	return fmt.Sprintf(format, parameters...)
}

// Split splits string by given separator into an slice
func (receiver *Utility) Split(value interface{}, separator string) interface{} {
	data := fmt.Sprintf("%v", value)
	return strings.Split(data, separator)
}

// Join joins given slice with give separator into a string
func (receiver *Utility) Join(slice interface{}, separator string) (interface{}, error) {
	interfaceSlice, err := mapstructutilpkg.GetUtility().SliceToInterfaceSlice(slice)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	stringSlice := make([]string, len(interfaceSlice))

	for index, item := range interfaceSlice {
		stringSlice[index] = fmt.Sprintf("%v", item)
	}

	return strings.Join(stringSlice, separator), nil
}

// Trim eliminate cutset string from both side of given string
func (receiver *Utility) Trim(value interface{}, cutset string) interface{} {
	data := fmt.Sprintf("%v", value)
	return strings.Trim(data, cutset)
}

// Replace replaces parts of string into another ones in a string
func (receiver *Utility) Replace(value interface{}, search string, replace string) interface{} {
	originalString := fmt.Sprintf("%v", value)
	rep := regexp.MustCompile(search)
	return rep.ReplaceAllString(originalString, replace)
}
