package main

import (
	"github.com/qb0C80aE/clay/cmd"
	_ "github.com/qb0C80aE/loam" // Install Loam module by importing
	_ "github.com/qb0C80aE/pottery" // Install Pottery module by importing
)

func main() {
	cmd.Execute()
}
