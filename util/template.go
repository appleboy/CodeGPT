package util

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
)

// Data defines a custom type for the template data.
type Data map[string]interface{}

var (
	templates    map[string]*template.Template
	templatesDir = "templates"
)

// NewTemplateByString parses a template from a string and executes it with the provided data.
// It returns the resulting string or an error if the template parsing or execution fails.
func NewTemplateByString(format string, data map[string]interface{}) (string, error) {
	t, err := template.New("message").Parse(format)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

// processTemplate processes the template with the given name and data.
// It returns the resulting bytes.Buffer or an error if the template execution fails.
func processTemplate(name string, data map[string]interface{}) (*bytes.Buffer, error) {
	t, ok := templates[name]
	if !ok {
		return nil, fmt.Errorf("template %s not found", name)
	}

	var tpl bytes.Buffer

	if err := t.Execute(&tpl, data); err != nil {
		return nil, err
	}

	return &tpl, nil
}

// GetTemplateByString returns the parsed template as a string.
// It returns an error if the template processing fails.
func GetTemplateByString(name string, data map[string]interface{}) (string, error) {
	tpl, err := processTemplate(name, data)
	return tpl.String(), err
}

// GetTemplateByBytes returns the parsed template as a byte slice.
// It returns an error if the template processing fails.
func GetTemplateByBytes(name string, data map[string]interface{}) ([]byte, error) {
	tpl, err := processTemplate(name, data)
	return tpl.Bytes(), err
}

// LoadTemplates loads all the templates found in the templates directory from the embedded filesystem.
// It returns an error if reading the directory or parsing any template fails.
func LoadTemplates(files embed.FS) error {
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
