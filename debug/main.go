// +build debug
// execute 'go run debug/main.go' and access localhost:8080/<resource> for API debugging or manual testing

package main

import (
	"github.com/qb0C80aE/clay/extensions"
	_ "github.com/qb0C80aE/clay/runtime" // Import runtime package to register Clay runtime
)

func main() {
	extensions.RegisteredRuntime().Run()
}
