// Package email provides email template management and rendering functionality.
package email

import (
	"bytes"
	"html/template"
)

// EmailTemplate manages HTML email templates for the application.
// It wraps the standard template.Template to provide email-specific
// template execution and management.
type EmailTemplate struct {
	templates *template.Template // Compiled HTML templates
}

// NewEmailTemplate initializes a new email template manager by loading
// all HTML templates from the templates/email directory.
//
// The function expects templates to be located in "templates/email/*.html"
// and will parse all files matching this pattern.
//
// Returns:
//   - *EmailTemplate: Template manager instance
//   - error: Any error encountered during template parsing
//
// Example Usage:
//
//	templates, err := NewEmailTemplate()
//	if err != nil {
//	    return fmt.Errorf("failed to initialize email templates: %w", err)
//	}
//
// Note: Templates must be valid HTML and follow Go template syntax.
func NewEmailTemplate() (*EmailTemplate, error) {
	// Load and parse all email templates
	templates, err := template.ParseGlob("templates/email/*.html")
	if err != nil {
		return nil, err
	}
	return &EmailTemplate{templates: templates}, nil
}

// ExecuteTemplate renders a specific template with the provided data.
//
// Parameters:
//   - name: Name of the template to execute (e.g., "welcome.html")
//   - data: Data to be passed to the template
//
// Returns:
//   - string: The rendered template as a string
//   - error: Any error encountered during template execution
//
// Example Usage:
//
//	data := WelcomeEmailData{
//	    Username: "John",
//	    LoginURL: "http://example.com/login",
//	    Year: 2024,
//	}
//	body, err := templates.ExecuteTemplate("welcome.html", data)
//
// Note: The data parameter must match the structure expected by the template.
func (et *EmailTemplate) ExecuteTemplate(name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := et.templates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
