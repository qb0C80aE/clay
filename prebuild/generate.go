// +build prebuild
// execute 'go generate -tags=prebuild prebuild/generate.go' or 'go generate -tags=prebuild ./...' or 'go generate -tags=prebuild prebuild/...' to generate build_information.go manually

package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

//go:generate go run generate.go

var clayVersionTemplate = template.Must(template.New("template").Parse(`package buildtime

import "github.com/qb0C80aE/clay/extension"

func init() {
	var programInformation = &clayProgramInformation{
		buildTime: "{{.BuildTime}}",
		claySubModuleInformationList: []*claySubModuleInformation{
			{{- range $i, $subModule := .SubModules}}
			{
				name:     "{{$subModule.Name}}",
				revision: "{{$subModule.Revision}}",
				version:  "{{$subModule.Version}}",
			},
			{{- end }}
		},
	}
	extension.RegisterProgramInformation(programInformation)
}
`))

type depLock struct {
	Projects []*depLockProject `toml:"projects"`
}

type depLockProject struct {
	Name     string `toml:"name"`
	Revision string `toml:"revision"`
	Version  string `toml:"version"`
}

func claySubModules(cwd string, depLockProjects []*depLockProject) ([]*depLockProject, error) {
	result := []*depLockProject{}

	out, err := exec.Command("git", "rev-parse", "@").Output()
	if err != nil {
		return nil, err
	}
	clayInfo := &depLockProject{
		Name:     "clay",
		Revision: strings.TrimSpace(string(out)),
		Version:  strings.TrimSpace(string(out)),
	}
	result = append(result, clayInfo)

	for _, depLockProject := range depLockProjects {
		_, err := os.Stat(filepath.Join(cwd, "vendor", depLockProject.Name, "clay_module.go"))
		if err != nil {
			continue
		}
		result = append(result, depLockProject)
	}
	return result, nil
}

func main() {
	cwd, err := filepath.Abs("..")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	buf, err := ioutil.ReadFile(filepath.Join(cwd, "Gopkg.toml"))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var depLock depLock
	err = toml.Unmarshal(buf, &depLock)
	clayModules, err := claySubModules(cwd, depLock.Projects)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	t := time.Now()
	const layout = "20060102150405"
	now := t.Format(layout)

	f, err := os.Create(filepath.Join(cwd, "buildtime", "build_information.go"))
	defer f.Close()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	clayVersionTemplate.Execute(f, struct {
		BuildTime  string
		SubModules []*depLockProject
	}{
		BuildTime:  now,
		SubModules: clayModules,
	})
}
