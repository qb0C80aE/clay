package logging

import (
	"fmt"
	"github.com/qb0C80aE/clay/logging"
)

var utility = &Utility{}

// Utility handles logging operation
type Utility struct {
}

// GetUtility returns the instance of utility
func GetUtility() *Utility {
	return utility
}

// Debug outputs debug message
func (receiver *Utility) Debug(log interface{}) interface{} {
	logging.Logger().Debug(log)
	return nil
}

// Debugf outputs formatted debug message
func (receiver *Utility) Debugf(format string, parameters []interface{}) interface{} {
	logging.Logger().Debugf(format, parameters...)
	return nil
}

// Info outputs information message
func (receiver *Utility) Info(log interface{}) error {
	logging.Logger().Info(log)
	return nil
}

// Infof outputs formatted information message
func (receiver *Utility) Infof(format string, parameters []interface{}) interface{} {
	logging.Logger().Infof(format, parameters...)
	return nil
}

// Warn outputs warning message
func (receiver *Utility) Warn(log interface{}) error {
	logging.Logger().Warn(log)
	return nil
}

// Warnf outputs formatted warning message
func (receiver *Utility) Warnf(format string, parameters []interface{}) interface{} {
	logging.Logger().Warnf(format, parameters...)
	return nil
}

// Critical outputs critical message
func (receiver *Utility) Critical(log interface{}) error {
	logging.Logger().Critical(log)
	return nil
}

// Criticalf outputs formatted critial message
func (receiver *Utility) Criticalf(format string, parameters []interface{}) interface{} {
	logging.Logger().Criticalf(format, parameters...)
	return nil
}

// Panic throws panic with message
func (receiver *Utility) Panic(log interface{}) error {
	panic(log)
}

// Panicf throws panic with formatted message
func (receiver *Utility) Panicf(format string, parameters []interface{}) interface{} {
	message := fmt.Sprintf(format, parameters...)
	panic(message)
}
