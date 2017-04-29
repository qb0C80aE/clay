// +build generate

package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

//go:generate go run generate.go

var clayVersionTemplate = template.Must(template.New("template").Parse(`package submodules

var programInformation = &ProgramInformation{
	BuildTime: "{{.BuildTime}}",
	SubModuleInformationList: []*SubModuleInformation{
		{{- range $i, $subModule := .SubModules}}
		{
			Name:     "{{$subModule.Name}}",
			Revision: "{{$subModule.Version}}",
		},
		{{- end }}
	},
}
`))

type glideLock struct {
	Imports []*glideLockImport `yaml:"imports"`
}

type glideLockImport struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

func claySubModules(cwd string, glideLockImports []*glideLockImport) ([]*glideLockImport, error) {
	result := []*glideLockImport{}

	out, err := exec.Command("git", "rev-parse", "@").Output()
	if err != nil {
		return nil, err
	}
	clayInfo := &glideLockImport{
		Name:    "clay",
		Version: strings.TrimSpace(string(out)),
	}
	result = append(result, clayInfo)

	for _, glideLockImport := range glideLockImports {
		_, err := os.Stat(filepath.Join(cwd, "vendor", glideLockImport.Name, "clay_submodule.go"))
		if err != nil {
			continue
		}
		result = append(result, glideLockImport)
	}
	return result, nil
}

func main() {
	cwd, err := filepath.Abs(".")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	buf, err := ioutil.ReadFile(filepath.Join(cwd, "glide.lock"))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var glideLock glideLock
	err = yaml.Unmarshal(buf, &glideLock)
	clayModules, err := claySubModules(cwd, glideLock.Imports)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	t := time.Now()
	const layout = "20060102150405"
	now := t.Format(layout)

	f, err := os.Create(filepath.Join(cwd, "submodules", "submodule_version.go"))
	defer f.Close()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	clayVersionTemplate.Execute(f, struct {
		BuildTime  string
		SubModules []*glideLockImport
	}{
		BuildTime:  now,
		SubModules: clayModules,
	})
}
