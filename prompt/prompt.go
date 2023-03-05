package prompt

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"text/template"
)

//go:embed templates/*
var files embed.FS

type Data map[string]interface{}

func init() {
	if err := loadTemplates(); err != nil {
		log.Fatal(err)
	}
}

const (
	SummarizeCommitTemplate = "summarize_commit.tpl"
)

var (
	templates    map[string]*template.Template
	templatesDir = "templates"
)

func GetTemplate(name string, data map[string]interface{}) (string, error) {
	t, ok := templates[name]
	if !ok {
		return "", fmt.Errorf("template %s not found", name)
	}

	var tpl bytes.Buffer

	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func loadTemplates() error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	tmplFiles, err := fs.ReadDir(files, templatesDir)
	if err != nil {
		return err
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(files, templatesDir+"/"+tmpl.Name())
		if err != nil {
			return err
		}

		templates[tmpl.Name()] = pt
	}
	return nil
}
