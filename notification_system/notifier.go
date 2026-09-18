package main

import "context"

// notifier is the strategy interface. this is strategy design pattern

type Notifier interface {
	//Send delivers the notification. Implementations should
	// respect ctx cancellation/deadlines (e.g. HTTP calls to
	// a 3rd party provider like SES/Twilio/FCM).
	Send(ctx context.Context, n *Notification) error

	// Name returns which channel this notifier serves - used
	// by the factory/registry.
	Name() ChannelType
}
