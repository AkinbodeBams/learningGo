package mailer

import "github.com/sendgrid/sendgrid-go"

type SendGRidMailer struct {
	fromEmail string
	apiKey string
	client *sendgrid.Client
}