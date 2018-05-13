// +build prebuild
// execute '[BUILD_ASSET=true|false] [BUILD_ASSET_SOURCE=...] [CLAY_...=...] go generate -tags=prebuild prebuild/generate.go' or '[BUILD_ASSET=true|false] [BUILD_ASSET_SOURCE=...] [CLAY_...=...] go generate -tags=prebuild ./...' or '[BUILD_ASSET=true|false] [BUILD_ASSET_SOURCE=...] [CLAY_...=...] go generate -tags=prebuild prebuild/...' to generate build_information.go manually
// ex.internal asset mode:
// CLAY_ASSET_MODE=internal BUILD_ASSET=true BUILD_ASSET_SOURCE=../examples/api_and_gui go generate -tags=prebuild prebuild/generate.go
// ex.external asset mode (default):
// go generate -tags=prebuild prebuild/generate.go

package main

import (
	"fmt"
	_ "github.com/qb0C80aE/clay/buildtime"
	"github.com/qb0C80aE/clay/extension"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"
)

//go:generate go run generate.go

var clayBuildTimeTemplate = template.Must(template.New("template").Parse(`package buildtime

import "github.com/qb0C80aE/clay/extension"

func init() {
	var programInformation = &clayProgramInformation{
		buildTime:  "{{ .BuildTime }}",
		branch:     "{{ .Branch }}",
		version:    "{{ .Version }}",
		commitHash: "{{ .CommitHash }}",
	}
	extension.RegisterProgramInformation(programInformation)

	var environmentalVariableSet = &defaultEnvironmentalVariableSet{
		clayConfigFilePath: "{{ .ClayConfigFilePath }}",
		clayHost:           "{{ .ClayHost }}",
		clayPort:           "{{ .ClayPort }}",
		clayDBMode:         "{{ .ClayDBMode }}",
		clayDBFilePath:     "{{ .ClayDBFilePath }}",
		clayAssetMode:      "{{ .ClayAssetMode }}",
	}
	extension.RegisterDefaultEnvironmentalVariableSet(environmentalVariableSet)
}
`))

func main() {
	now := time.Now().UTC().Format(time.RFC3339)

	cwd, err := filepath.Abs("..")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	f, err := os.Create(filepath.Join(cwd, "buildtime", "buildtime.go"))
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	repository, err := git.PlainOpen(cwd)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	head, err := repository.Head()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	tagsIterator, err := repository.Tags()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	version := "v0.0.0"
	tagsIterator.ForEach(func(reference *plumbing.Reference) error {
		if head.Hash().String() == reference.Hash().String() {
			version = reference.Name().Short()
		}
		return nil
	})

	environmentalVariable := extension.GetCurrentEnvironmentalVariableSet()

	clayBuildTimeTemplate.Execute(f, struct {
		BuildTime          string
		Branch             string
		Version            string
		CommitHash         string
		ClayConfigFilePath string
		ClayHost           string
		ClayPort           string
		ClayDBMode         string
		ClayDBFilePath     string
		ClayAssetMode      string
	}{
		BuildTime:          now,
		Branch:             head.Name().Short(),
		Version:            version,
		CommitHash:         head.Hash().String(),
		ClayConfigFilePath: environmentalVariable.GetClayConfigFilePath(),
		ClayHost:           environmentalVariable.GetClayHost(),
		ClayPort:           environmentalVariable.GetClayPort(),
		ClayDBMode:         environmentalVariable.GetClayDBMode(),
		ClayDBFilePath:     environmentalVariable.GetClayDBFilePath(),
		ClayAssetMode:      environmentalVariable.GetClayAssetMode(),
	})

	if os.Getenv("BUILD_ASSET") == "true" {
		output, err := exec.Command("go-assets-builder", "-p", "asset", "-o", "../asset/asset.go", os.Getenv("BUILD_ASSET_SOURCE"), "-s", "/").CombinedOutput()
		fmt.Println(string(output))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
