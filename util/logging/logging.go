package logging

import (
	"errors"
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
func (receiver *Utility) Debug(log interface{}) error {
	message, ok := log.(string)
	if !ok {
		return errors.New("log is not a string")
	}

	logging.Logger().Debug(message)
	return nil
}

// Debugf outputs formatted debug message
func (receiver *Utility) Debugf(log interface{}, parameters []interface{}) error {
	message, ok := log.(string)
	if !ok {
		return errors.New("log is not a string")
	}

	logging.Logger().Debugf(message, parameters...)
	return nil
}

// Info outputs information message
func (receiver *Utility) Info(log interface{}) error {
	message, ok := log.(string)
	if !ok {
		return errors.New("log is not a string")
	}

	logging.Logger().Info(message)
	return nil
}

// Infof outputs formatted information message
func (receiver *Utility) Infof(log interface{}, parameters []interface{}) error {
	message, ok := log.(string)
	if !ok {
		return errors.New("log is not a string")
	}

	logging.Logger().Infof(message, parameters...)
	return nil
}

// Warn outputs warning message
func (receiver *Utility) Warn(log interface{}) error {
	message, ok := log.(string)
	if !ok {
		return errors.New("log is not a string")
	}

	logging.Logger().Warn(message)
	return nil
}

// Warnf outputs formatted warning message
func (receiver *Utility) Warnf(log interface{}, parameters []interface{}) error {
	message, ok := log.(string)
	if !ok {
		return errors.New("log is not a string")
	}

	logging.Logger().Warnf(message, parameters...)
	return nil
}

// Critical outputs critical message
func (receiver *Utility) Critical(log interface{}) error {
	message, ok := log.(string)
	if !ok {
		return errors.New("log is not a string")
	}

	logging.Logger().Critical(message)
	return nil
}

// Criticalf outputs formatted critial message
func (receiver *Utility) Criticalf(log interface{}, parameters []interface{}) error {
	message, ok := log.(string)
	if !ok {
		return errors.New("log is not a string")
	}

	logging.Logger().Criticalf(message, parameters...)
	return nil
}

// Panic throws panic with message
func (receiver *Utility) Panic(log interface{}) error {
	message, ok := log.(string)
	if !ok {
		return errors.New("log is not a string")
	}

	panic(message)
	return nil
}

// Panicf throws panic with formatted message
func (receiver *Utility) Panicf(log interface{}, parameters []interface{}) error {
	message, ok := log.(string)
	if !ok {
		return errors.New("log is not a string")
	}

	message = fmt.Sprintf(message, parameters...)
	panic(message)
	return nil
}
