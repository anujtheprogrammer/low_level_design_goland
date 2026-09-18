package main

import (
	"context"
	"fmt"
	"time"
)

// NotificationService is the Facade: the single entry point the rest
// of the codebase (or an HTTP handler) calls to send a notification.
// It wires together preference checks, templating, rate limiting,
// the correct Notifier (via factory) and retries.
type NotificationService struct {
	observable
	factory     *ChannelFactory
	prefs       PreferenceRepository
	templates   TemplateService
	RateLimiter *RateLimiter
	retry       *RetryPolicy
}

func NewNotificationService(factory *ChannelFactory, prefs PreferenceRepository, templates TemplateService, rl *RateLimiter, retry *RetryPolicy) *NotificationService {
	return &NotificationService{factory: factory, prefs: prefs, templates: templates, RateLimiter: rl, retry: retry}
}

// Send is the main use case. Steps:
// 1. validate the request
// 2. check user preference (skip silently if opted out)
// 3. render the template into title/body
// 4. rate-limit check
// 5. resolve the right Notifier via the factory
// 6. send, retrying with backoff on failure
// 7. update status + notify observers

func (s *NotificationService) Send(ctx context.Context, n *Notification) error {
	if err := n.Validate(); err != nil {
		return err
	}
	if !s.prefs.IsOptedIN(n.UserID, n.ChannelId) {
		n.Status = StatusFailed
		return ErrUserOptedOut
	}

	title, body, err := s.templates.Render(n.TemplateID, n.Data)
	if err != nil {
		return fmt.Errorf("template render failed : %w", err)
	}
	n.Title, n.Body = title, body

	if !s.RateLimiter.Allow() {
		return ErrRateLimited
	}

	notifier, err := s.factory.GetNotifier(n.ChannelId)
	if err != nil {
		return err
	}
	sendErr := s.sendWithRetry(ctx, notifier, n)
	s.notify(n) // Observer pattern: tell metrics/audit/webhook listeners
	return sendErr
}

// sendWithRetry keeps trying notifier.Send until it succeeds, the
// retry budget is exhausted, or the caller's context is cancelled.
func (s *NotificationService) sendWithRetry(ctx context.Context, notifier Notifier, n *Notification) error {
	for attempt := 0; ; attempt++ {
		err := notifier.Send(ctx, n)
		if err == nil {
			n.Status = StatusSent
			return nil
		}
		n.RetryCount = attempt + 1
		if !s.retry.ShouldRetry(attempt) {
			n.Status = StatusFailed
			return err
		}
		select {
		case <-ctx.Done():
			n.Status = StatusFailed
			return ctx.Err()
		case <-time.After(s.retry.NextDelay(attempt)):
			// backoff elapsed, loop again and retry
		}
	}
}
