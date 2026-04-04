package notification

import (
	"log"
)

// Notifier defines the interface for sending notifications
type Notifier interface {
	Send(recipient, subject, message string) error
}

// logNotifier is a simple implementation that logs to console
type logNotifier struct{}

// NewLogNotifier creates a new instance of logNotifier
func NewLogNotifier() Notifier {
	return &logNotifier{}
}

// Send logs the notification details to the console
func (n *logNotifier) Send(recipient, subject, message string) error {
	log.Printf("[NOTIFICATION] To: %s | Subject: %s | Message: %s", recipient, subject, message)
	return nil
}
