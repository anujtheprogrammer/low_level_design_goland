package main

import "time"

// channeltype rerpresents the delivery channel for a notification
type ChannelType string

const (
	ChannelEmail ChannelType = "EMAIL"
	ChannelSMS   ChannelType = "SMS"
	ChannelPush  ChannelType = "PUSH"
)

// status represents the lifecycle state of notification
type Status string

const (
	StatusPending   Status = "PENDING"
	StatusSent      Status = "SENT"
	StatusFailed    Status = "Failed"
	StatusDelivered Status = "DELIVERED"
)

// priority sees how urgent the notification is
type Priority int

const (
	PriorityLow Priority = iota
	PriorityNormal
	PriorityHigh
)

// notification is the core struct flow through the application
type Notification struct {
	ID         string
	UserID     string
	ChannelId  ChannelType
	TemplateID string
	Data       map[string]string
	Title      string
	Body       string
	Status     Status
	Priority   Priority
	CreatedAt  time.Time
	RetryCount int
}

// validate basic check
func (n *Notification) Validate() error {
	if n.UserID == "" {
		return ErrMissingUser
	}
	if n.ChannelId == "" {
		return ErrMissingChannel
	}
	return nil
}
