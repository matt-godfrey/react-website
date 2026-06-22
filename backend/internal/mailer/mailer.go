package mailer

import "context"

type Message struct {
	From    string
	To      []string
	Subject string
	HTML    string
	Text    string
}

type Mailer interface {
	SendMail(ctx context.Context, msg Message) (string, error)
}
