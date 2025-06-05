package common

type NotificationEvent struct {
	Type      string `json:"type"`         // "email", "sms", etc
	Recipient string `json:"recipient"`    // email or phone
	Subject   string `json:"subject"`      // optional
	Body      string `json:"body"`         // the actual message
	TriggerBy string `json:"triggered_by"` // maybe userID or action
}
