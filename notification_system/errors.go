package main

import "errors"

var (
	ErrMissingUser    = errors.New("Notification : UserID is required")
	ErrMissingChannel = errors.New("notification: channel is required")
	ErrUnknownChannel = errors.New("notification: no notifier registered for channel")
	ErrUserOptedOut   = errors.New("notification: user has opted out of this channel")
	ErrRateLimited    = errors.New("notification: rate limit exceeded, try later")
)
