// +build prebuild
// execute 'go generate -tags=prebuild prebuild/generate.go' or 'go generate -tags=prebuild ./...' or 'go generate -tags=prebuild prebuild/...' to generate build_information.go manually

package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

//go:generate go run generate.go

var clayVersionTemplate = template.Must(template.New("template").Parse(`package buildtime

import "github.com/qb0C80aE/clay/extension"

func init() {
	var programInformation = &clayProgramInformation{
		buildTime:  "{{ .BuildTime }}",
		branch:     "{{ .Branch }}",
		version:    "{{ .Version }}",
		commitHash: "{{ .CommitHash }}",
	}
	extension.RegisterProgramInformation(programInformation)
}
`))

func main() {
	now := time.Now().UTC().Format(time.RFC3339)

	cwd, err := filepath.Abs("..")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	f, err := os.Create(filepath.Join(cwd, "buildtime", "build_information.go"))
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

	clayVersionTemplate.Execute(f, struct {
		BuildTime  string
		Branch     string
		Version    string
		CommitHash string
	}{
		BuildTime:  now,
		Branch:     head.Name().Short(),
		Version:    version,
		CommitHash: head.Hash().String(),
	})
}
