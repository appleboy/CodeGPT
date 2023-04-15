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
func GetTemplateByString(name string, data map[string]interface{}) (string, error) {
	tpl, err := processTemplate(name, data)
	return tpl.String(), err
}

// GetTemplateByBytes returns the parsed template as a byte.
func GetTemplateByBytes(name string, data map[string]interface{}) ([]byte, error) {
	tpl, err := processTemplate(name, data)
	return tpl.Bytes(), err
}

// LoadTemplates loads all the templates found in the templates directory.
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
