package db

import (
	"math"
	"strconv"

	"errors"
	"net/url"
	"reflect"
)

const (
	defaultLimit = "10"
	maxLimit     = int(^uint(0) >> 1)
	defaultPage  = "1"
	defaultOrder = "desc"
)

// Parameter is the struct what contains parameters related to db
type Parameter struct {
	Filters  map[string]string
	Preloads string
	Sort     string
	Limit    int
	Page     int
	LastID   int
	Order    string
	IsLastID bool
}

// NewParameter creates a new Parameter instance
func NewParameter(query url.Values, model interface{}) (*Parameter, error) {
	parameter := &Parameter{}

	if err := parameter.initialize(query, model); err != nil {
		return nil, err
	}

	return parameter, nil
}

func (parameter *Parameter) getQueryArray(query url.Values, key string) ([]string, bool) {
	if values, ok := query[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

func (parameter *Parameter) getQuery(query url.Values, key string) (string, bool) {
	if values, ok := parameter.getQueryArray(query, key); ok {
		return values[0], ok
	}
	return "", false
}

// DefaultQuery get a query parameter. If key does not exist, it returns defaultValue
func (parameter *Parameter) DefaultQuery(query url.Values, key string, defaultValue string) string {
	if value, ok := parameter.getQuery(query, key); ok {
		return value
	}
	return defaultValue
}

func (parameter *Parameter) initialize(query url.Values, model interface{}) error {
	vs := reflect.ValueOf(model)
	for vs.Kind() == reflect.Ptr {
		vs = vs.Elem()
	}
	if !vs.IsValid() {
		return errors.New("the model is invalid")
	}
	if !vs.CanInterface() {
		return errors.New("the model cannot interface")
	}
	value := vs.Interface()

	parameter.Filters = filterToMap(query, value)
	parameter.Preloads = query.Get("preloads")
	parameter.Sort = query.Get("sort")

	limit, err := validate(parameter.DefaultQuery(query, "limit", defaultLimit))
	if err != nil {
		return err
	}

	parameter.Limit = int(math.Max(1, math.Min(float64(maxLimit), float64(limit))))
	page, err := validate(parameter.DefaultQuery(query, "page", defaultPage))
	if err != nil {
		return err
	}

	parameter.Page = int(math.Max(1, float64(page)))
	lastID, err := validate(query.Get("last_id"))
	if err != nil {
		return err
	}

	if lastID != -1 {
		parameter.IsLastID = true
		parameter.LastID = int(math.Max(0, float64(lastID)))
	}

	parameter.Order = parameter.DefaultQuery(query, "order", defaultOrder)
	return nil
}

func validate(s string) (int, error) {
	if s == "" {
		return -1, nil
	}

	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return num, nil
}
