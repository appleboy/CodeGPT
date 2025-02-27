package util

import (
	"bytes"
	"embed"
	"html/template"
	"os"
	"testing"
)

func TestNewTemplateByString(t *testing.T) {
	data := map[string]interface{}{
		"Name": "John Doe",
	}

	expected := "Hello, John Doe!"

	actual, err := NewTemplateByString("Hello, {{.Name}}!", data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Errorf("Expected %q but got %q", expected, actual)
	}
}

func TestNewTemplateByStringWithCustomVars(t *testing.T) {
	data := map[string]interface{}{}
	vars := ConvertToMap([]string{"Name=John Doe", "Message=Hello"})
	for k, v := range vars {
		data[k] = v
	}

	expected := "Hello, John Doe! Hello"

	actual, err := NewTemplateByString("Hello, {{.Name}}! {{.Message}}", data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Errorf("Expected %q but got %q", expected, actual)
	}
}

func TestProcessTemplate(t *testing.T) {
	// Set up test data
	testTemplateName := "foo.tmpl"
	testTemplateText := "Hello {{.Name}}!"
	testData := Data{"Name": "World"}

	// Parse test template
	tmpl, err := template.New(testTemplateName).Parse(testTemplateText)
	if err != nil {
		t.Errorf("Failed to parse test template: %v", err)
	}

	// Add test template to templates map
	templates = make(map[string]*template.Template)
	templates[testTemplateName] = tmpl

	// Process test template
	buf, err := processTemplate(testTemplateName, testData)
	if err != nil {
		t.Errorf("Failed to process template: %v", err)
	}

	// Check the output
	expected := "Hello World!"
	if buf.String() != expected {
		t.Errorf("Unexpected output. Got: %v, Want: %v", buf.String(), expected)
	}
}

func TestLoadTemplatesFromDir(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a sample template file in the temporary directory
	templateContent := "Hello, {{.Name}}!"
	templateFile := "test.tmpl"
	err := os.WriteFile(tempDir+"/"+templateFile, []byte(templateContent), 0o600)
	if err != nil {
		t.Fatalf("Failed to create test template file: %v", err)
	}

	// Load templates from the temporary directory
	err = LoadTemplatesFromDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to load templates from directory: %v", err)
	}

	// Check if the template was loaded correctly
	tmpl, ok := templates[templateFile]
	if !ok {
		t.Fatalf("Template %s not found in loaded templates", templateFile)
	}

	// Process the loaded template
	data := Data{"Name": "World"}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		t.Fatalf("Failed to execute loaded template: %v", err)
	}

	// Check the output
	expected := "Hello, World!"
	if buf.String() != expected {
		t.Errorf("Unexpected output. Got: %v, Want: %v", buf.String(), expected)
	}
}

// Create an embedded filesystem with a sample template
//
//go:embed templates/*
var testFiles embed.FS

func TestLoadTemplates(t *testing.T) {
	// Load templates from the embedded filesystem
	err := LoadTemplates(testFiles)
	if err != nil {
		t.Fatalf("Failed to load templates from embedded filesystem: %v", err)
	}

	// Check if the template was loaded correctly
	templateFile := "test.tmpl"
	tmpl, ok := templates[templateFile]
	if !ok {
		t.Fatalf("Template %s not found in loaded templates", templateFile)
	}

	// Process the loaded template
	data := Data{"Name": "World"}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		t.Fatalf("Failed to execute loaded template: %v", err)
	}

	// Check the output
	expected := "Hello, World!"
	if buf.String() != expected {
		t.Errorf("Unexpected output. Got: %v, Want: %v", buf.String(), expected)
	}
}
