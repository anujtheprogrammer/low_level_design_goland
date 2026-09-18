package main

// StatusObserver lets other parts of the system (metrics, audit log,
// webhook callbacks) react whenever a notification's status changes,
// without NotificationService needing to know about them directly.
type StatusObserver interface {
	OnStatusChange(n *Notification)
}

// Subject-side helper embedded in NotificationService.
type observable struct {
	observers []StatusObserver
}

func (o *observable) Subscribe(obs StatusObserver) {
	o.observers = append(o.observers, obs)
}

func (o *observable) notify(n *Notification) {
	for _, ob := range o.observers {
		ob.OnStatusChange(n) // kept synchronous & simple for the interview;
	} // could be made async via a channel if asked.
}
