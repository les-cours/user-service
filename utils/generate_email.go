package utils

import (
	"bytes"
	"embed"
	"github.com/les-cours/user-service/env"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
)

//go:embed templates/*
var htmlFS embed.FS

func GenerateEmail(receiver, subject, templateName string, data interface{}) (*mail.SGMailV3, error) {
	var t = template.New(templateName + ".html")
	t, err := t.ParseFS(htmlFS, "templates/"+templateName+".html")
	if err != nil {
		return nil, err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return nil, err
	}

	result := tpl.String()

	var from = mail.NewEmail("Emplorium", env.Settings.Noreply.Email)
	var to = mail.NewEmail("", receiver)
	var message = mail.NewSingleEmail(from, subject, to, "", result)

	return message, nil
}
