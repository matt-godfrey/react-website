package mailer

import (
	"context"
	"log"

	"github.com/resend/resend-go/v3"
)

type ResendMailer struct {
	apiKey string
	from   string
}

func NewResendMailer(apiKey, from string) *ResendMailer {
	return &ResendMailer{
		apiKey: apiKey,
		from:   from,
	}
}

func (m *ResendMailer) SendMail(ctx context.Context, msg Message) (string, error) {
	client := resend.NewClient(m.apiKey)

	params := &resend.SendEmailRequest{
		From:    m.from,
		To:      msg.To,
		Subject: msg.Subject,
		Html:    msg.HTML,
		Text:    msg.Text,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return sent.Id, nil
}
