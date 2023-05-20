package util

import (
	"html/template"
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
	testTemplateName := "test.tmpl"
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
