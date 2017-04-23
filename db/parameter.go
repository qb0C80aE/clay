package db

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"reflect"
)

const (
	defaultLimit = "25"
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
func NewParameter(c *gin.Context, model interface{}) (*Parameter, error) {
	parameter := &Parameter{}

	if err := parameter.initialize(c, model); err != nil {
		return nil, err
	}

	return parameter, nil
}

func (parameter *Parameter) initialize(c *gin.Context, model interface{}) error {
	vs := reflect.ValueOf(model)
	for vs.Kind() == reflect.Ptr {
		vs = vs.Elem()
	}
	if !vs.IsValid() {
		return nil
	}
	if !vs.CanInterface() {
		return nil
	}
	value := vs.Interface()

	parameter.Filters = filterToMap(c, value)
	parameter.Preloads = c.Query("preloads")
	parameter.Sort = c.Query("sort")

	limit, err := validate(c.DefaultQuery("limit", defaultLimit))
	if err != nil {
		return err
	}

	parameter.Limit = int(math.Max(1, math.Min(10000, float64(limit))))
	page, err := validate(c.DefaultQuery("page", defaultPage))
	if err != nil {
		return err
	}

	parameter.Page = int(math.Max(1, float64(page)))
	lastID, err := validate(c.Query("last_id"))
	if err != nil {
		return err
	}

	if lastID != -1 {
		parameter.IsLastID = true
		parameter.LastID = int(math.Max(0, float64(lastID)))
	}

	parameter.Order = c.DefaultQuery("order", defaultOrder)
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
