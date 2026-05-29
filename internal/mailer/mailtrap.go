package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"time"

	gomail "gopkg.in/mail.v2"
)

type mailtrapClient struct {
	fromEmail string
	apiKey    string
}

func NewMailTrapClient(apiKey, fromEmail string) (Client, error) {
	if apiKey == "" {
		return nil, errors.New("api key is required")
	}

	return &mailtrapClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m *mailtrapClient) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return 0, err
	}

	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return 0, err
	}

	body := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(body, "body", data); err != nil {
		return 0, err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", message.FormatAddress(m.fromEmail, FromName))
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())
	message.AddAlternative("text/html", body.String())

	host := "live.smtp.mailtrap.io"
	if isSandbox {
		host = "sandbox.smtp.mailtrap.io"
	}

	dialer := gomail.NewDialer(host, 587, "d68a754927924d", m.apiKey)

	var retryErr error
	for i := range maxRetries {
		retryErr = dialer.DialAndSend(message)
		if retryErr != nil {
			// exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		return 200, nil
	}

	return -1, fmt.Errorf("failed to send email after %d attempt, error: %v", maxRetries, retryErr)
}
