package io

import (
	"io/ioutil"
	"os"
)

var utility = &Utility{}

// Utility handles io operation
type Utility struct {
}

// GetUtility returns the instance of utility
func GetUtility() *Utility {
	return utility
}

// ReadFile reads file
func (receiver *Utility) ReadFile(fileName interface{}) ([]byte, error) {
	return ioutil.ReadFile(fileName.(string))
}

// WriteFile writes contents into the file
func (receiver *Utility) WriteFile(fileName interface{}, data interface{}, perm os.FileMode) error {
	return ioutil.WriteFile(fileName.(string), data.([]byte), perm)
}
