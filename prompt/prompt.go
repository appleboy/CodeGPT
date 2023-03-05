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

// Data define a custom type for the template data.
type Data map[string]interface{}

func init() {
	if err := loadTemplates(); err != nil {
		log.Fatal(err)
	}
}

const (
	SummarizeFileDiffTemplate = "summarize_file_diff.tmpl"
	SummarizeTitleTemplate    = "summarize_title.tmpl"
	TranslationTemplate       = "translation.tmpl"
)

var (
	templates    map[string]*template.Template
	templatesDir = "templates"
)

// GetTemplate returns the parsed template as a string.
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

// loadTemplates loads all the templates found in the templates directory.
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
