package mail

import (
	"bytes"
	"html/template"
	"net/smtp"
)

type SMTPMailer struct {
	client *smtp.Client
}

func (m *SMTPMailer) SendTemplatedMail(name string, data interface{}) error {
	t := template.Must(template.New("name").Parse(""))
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return err
	}
	return nil
}

func CreateMailer() (*SMTPMailer, error) {
	c, err := smtp.Dial("mail.example.com:25")
	if err != nil {
		return nil, err
	}
	return &SMTPMailer{client: c}, nil
}
