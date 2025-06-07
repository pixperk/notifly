package common

import "github.com/google/uuid"

type NotificationEvent struct {
	NotificationId uuid.UUID `json:"notification_id"` // unique ID for the notification
	Type           string    `json:"type"`            // "email", "sms", etc
	Recipient      string    `json:"recipient"`       // email or phone
	Subject        string    `json:"subject"`         // optional
	Body           string    `json:"body"`            // the actual message
	TriggerBy      string    `json:"triggered_by"`    // maybe userID or action
}
