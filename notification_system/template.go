package main

import (
	"fmt"
	"strings"
)

// TemplateService renders a template ID + dynamic data into a
// concrete (subject, body) pair. Kept as an interface so it can be
// swapped for a DB-backed or i18n-aware implementation later.
type TemplateService interface {
	Render(templateID string, data map[string]string) (subject, body string, err error)
}

type Template struct {
	Subject string
	Body    string
}

type inMemoryTemplateService struct {
	templates map[string]Template
}

func NewInMemoryTemplateService(templates map[string]Template) *inMemoryTemplateService {
	return &inMemoryTemplateService{templates: templates}
}

// Render does simple {{placeholder}} substitution. A production
// version would use text/template for safety and richer logic.
func (t *inMemoryTemplateService) Render(templateID string, data map[string]string) (string, string, error) {
	tmpl, ok := t.templates[templateID]
	if !ok {
		return "", "", fmt.Errorf("template %q not found", templateID)
	}
	subject, body := tmpl.Subject, tmpl.Body
	for key, val := range data {
		placeholder := "{{" + key + "}}"
		subject = strings.ReplaceAll(subject, placeholder, val)
		body = strings.ReplaceAll(body, placeholder, val)
	}
	return subject, body, nil
}
