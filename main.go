package main

import (
	"os"
	"strconv"

	"flag"
	"fmt"
	"github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/server"
	"github.com/qb0C80aE/clay/submodules"
)

var showVersion = flag.Bool("version", false, "show version")

func main() {
	flag.Parse()
	if *showVersion {
		programInformation := submodules.BuildInformation()
		fmt.Printf("Clay build %s\n", programInformation.BuildTime)
		for _, subModuleInformation := range programInformation.SubModuleInformationList {
			fmt.Printf("  module %s\n    version: %s\n", subModuleInformation.Name, subModuleInformation.Revision)
		}
		return
	}

	host := "localhost"
	port := "8080"

	if h := os.Getenv("HOST"); h != "" {
		host = h
	}

	if p := os.Getenv("PORT"); p != "" {
		if _, err := strconv.Atoi(p); err == nil {
			port = p
		}
	}

	os.Setenv("HOST", host)
	os.Setenv("PORT", port)

	database := db.Connect()
	s := server.Setup(database)

	s.Run(fmt.Sprintf("%s:%s", host, port))

}
