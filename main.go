package main

import (
	"flag"
	"fmt"
	"github.com/qb0C80aE/clay/extensions"
	_ "github.com/qb0C80aE/clay/runtime" // Import runtime package to register Clay runtime
)

var showVersion = flag.Bool("version", false, "show version")

func main() {
	flag.Parse()
	if *showVersion {
		programInformation := extensions.RegisteredProgramInformation()
		fmt.Printf("Clay build-%s\n", programInformation.BuildTime())
		subModuleInformationList := programInformation.SubModuleInformationList()
		for _, subModuleInformation := range subModuleInformationList {
			fmt.Printf("  module %s\n    revision: %s\n", subModuleInformation.Name(), subModuleInformation.Revision())
		}
		return
	}

	extensions.RegisteredRuntime().Run()
}
