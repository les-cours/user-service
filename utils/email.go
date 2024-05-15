package utils

import (
	"bytes"
	"fmt"
	"github.com/les-cours/user-service/env"
	"html/template"
	"log"
	"net/smtp"
)

var auth smtp.Auth

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject string) *Request {
	auth = smtp.PlainAuth("", "chouaib.chouache@univ-constantine2.dz", "s3yN5Ffz6MkCod", "smtp.gmail.com")

	return &Request{
		from:    env.Settings.NoreplyEmail,
		to:      to,
		subject: subject,
	}
}

func (r *Request) SendEmail() error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, r.from, r.to, msg); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GenerateEmail(receiver, subject, templateName string, data interface{}) error {

	r := NewRequest([]string{receiver}, subject)

	err := r.ParseTemplate("utils/templates/"+templateName+".html", data)

	if err != nil {
		fmt.Println(err)
		return err
	}
	err = r.SendEmail()
	if err != nil {
		return err
	}

	return nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
