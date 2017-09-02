package db

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID      uint   `json:"id,omitempty" form:"id"`
	Name    string `json:"name,omitempty" form:"name"`
	Engaged bool   `json:"engaged,omitempty" form:"engaged"`
}

func contains(ss map[string]string, s string) bool {
	_, ok := ss[s]

	return ok
}

func TestFilterToMap(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?q[id]=1,5,100&q[name]=hoge,fuga&q[unexisted_field]=null&q[nested.id]=1,2,3", nil)
	c := &gin.Context{
		Request: req,
	}
	value, preloadFilter := filterToMap(c.Request.URL.Query())

	if !contains(value, "id") {
		t.Fatalf("Filter should have `id` key.")
	}

	if !contains(value, "name") {
		t.Fatalf("Filter should have `name` key.")
	}

	if contains(value, "nested.id") {
		t.Fatalf("Filter should not have `nested.id` key.")
	}

	if !contains(preloadFilter["nested"], "id") {
		t.Fatalf("Preload Filter should have `id` key.")
	}

	// unexisted keys should be allowed because those are ignored in GORM.

	if value["id"] != "1,5,100" {
		t.Fatalf("filters[\"id\"] expected: `1,5,100`, actual: %s", value["id"])
	}

	if value["name"] != "hoge,fuga" {
		t.Fatalf("filters[\"name\"] expected: `hoge,fuga`, actual: %s", value["id"])
	}

	if preloadFilter["nested"]["id"] != "1,2,3" {
		t.Fatalf("nested_filters[\"id\"] expected: `1,2,3`, actual: %s", value["id"])
	}
}
