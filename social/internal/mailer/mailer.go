package mailer

type Client interface {
	Send(templateFile, userame , email string, data any , isSandBox bool) error
}