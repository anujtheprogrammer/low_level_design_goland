package main

import (
	"context"
	"fmt"
)

// SMTPClient is a tiny interface over the real email provider SDK
// (e.g. AWS SES). Keeping it as an interface lets us mock it in tests.
type SMTPClient interface {
	SendMail(ctx context.Context, to, subject, body string) error
}

// EmailNotifier is one concrete Strategy implementing Notifier.
type EmailNotifier struct {
	Client SMTPClient
}

func NewEmailNotifier(c SMTPClient) *EmailNotifier {
	return &EmailNotifier{Client: c}
}

func (e *EmailNotifier) Name() ChannelType {
	return ChannelEmail
}

func (e *EmailNotifier) Send(ctx context.Context, n Notification) error {
	err := e.Client.SendMail(ctx, n.UserID, n.Title, n.Body)
	if err != nil {
		return fmt.Errorf("email send failed: %w", err)
	}
	return nil
}

// similarly we can do for all and this is the strategy design pattern
