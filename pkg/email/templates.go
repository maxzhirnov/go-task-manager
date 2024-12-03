package email

import (
	"bytes"
	"html/template"
)

type EmailTemplate struct {
	templates *template.Template
}

type WelcomeEmailData struct {
	Username string
	LoginURL string
	Year     int
}

func NewEmailTemplate() (*EmailTemplate, error) {
	// Загружаем все шаблоны из директории templates/email
	templates, err := template.ParseGlob("templates/email/*.html")
	if err != nil {
		return nil, err
	}
	return &EmailTemplate{templates: templates}, nil
}

func (et *EmailTemplate) ExecuteTemplate(name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := et.templates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
