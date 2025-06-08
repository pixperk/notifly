package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/pixperk/notifly/common"
	"github.com/pixperk/notifly/notification"
)

var ErrInvalidEmail = fmt.Errorf("invalid email format")

type BrevoEmailRequest struct {
	Sender struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"sender"`
	To []struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"to"`
	Subject     string `json:"subject"`
	HTMLContent string `json:"htmlContent"`
}

func SendEmail(event common.NotificationEvent, cfg notification.Config) error {
	url := "https://api.brevo.com/v3/smtp/email"

	reqBody := BrevoEmailRequest{}
	reqBody.Sender.Name = "No-Reply"
	reqBody.Sender.Email = cfg.FromAddress

	// Validate recipient email
	if !checkIfEmailValid(event.Recipient) {
		return ErrInvalidEmail
	}

	reqBody.To = []struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}{
		{Email: event.Recipient, Name: "Recipient"},
	}
	reqBody.Subject = event.Subject

	// Check if the body is HTML content
	if checkIfHTML(event.Body) {
		reqBody.HTMLContent = event.Body
	} else {
		// If not HTML, wrap in paragraph tags to make it valid HTML
		reqBody.HTMLContent = fmt.Sprintf("<p>%s</p>", event.Body)
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal email request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("api-key", cfg.BrevoAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Brevo request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status from Brevo API: %s", resp.Status)
	}

	return nil
}

// checkIfHTML determines if the provided string is likely HTML content
// It checks for common HTML indicators like tags, doctype declarations, etc.
func checkIfHTML(body string) bool {
	// Empty content can't be HTML
	if len(body) == 0 {
		return false
	}

	// Look for HTML doctype declaration
	if bytes.Contains([]byte(body), []byte("<!DOCTYPE html>")) ||
		bytes.Contains([]byte(body), []byte("<!doctype html>")) {
		return true
	}

	// Check for common HTML tags
	commonTags := []string{"<html", "<body", "<div", "<p>", "<span", "<h1", "<table"}
	for _, tag := range commonTags {
		if bytes.Contains([]byte(body), []byte(tag)) {
			return true
		}
	}
	// Basic check for opening and closing angle brackets
	// This is a simple heuristic - the content should start with < and end with >
	// and have a reasonable number of angle brackets throughout
	if body[0] == '<' && body[len(body)-1] == '>' {
		// Count angle brackets - proper HTML should have a balanced number
		openCount := bytes.Count([]byte(body), []byte("<"))
		closeCount := bytes.Count([]byte(body), []byte(">"))

		// If we have multiple balanced brackets, it's likely HTML
		return openCount > 1 && closeCount > 1 && openCount == closeCount
	}

	return false
}

func checkIfEmailValid(email string) bool {
	// Basic email validation regex
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
