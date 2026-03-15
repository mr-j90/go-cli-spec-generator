// Package sms provides HTTP handlers for incoming Twilio SMS webhooks.
package sms

import (
	"fmt"
	"net/http"

	"github.com/zyx-holdings/go-spec/internal/twilio"
)

// Handler handles incoming Twilio SMS webhook POST requests.
// It validates the X-Twilio-Signature header before processing the message.
//
// The authToken parameter must be the Twilio account auth token; the fullURL
// parameter must be the exact URL (scheme + host + path + query) that was
// registered with Twilio for this webhook — any mismatch will cause valid
// requests to be rejected.
type Handler struct {
	authToken string
	fullURL   string
}

// NewHandler creates a Handler with the supplied Twilio credentials.
func NewHandler(authToken, fullURL string) *Handler {
	return &Handler{authToken: authToken, fullURL: fullURL}
}

// ServeHTTP implements http.Handler.
//
// Request lifecycle:
//  1. Parse the POST form body so we have access to Twilio's parameters.
//  2. Validate X-Twilio-Signature against the HMAC-SHA1 of the URL + params.
//  3. Reject with 403 if the signature is missing or invalid.
//  4. Process the verified webhook payload.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ParseForm populates r.PostForm, which we need for signature validation.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	sig := r.Header.Get("X-Twilio-Signature")
	if sig == "" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if !twilio.ValidateSignature(h.authToken, sig, h.fullURL, r.PostForm) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	from := r.PostFormValue("From")
	body := r.PostFormValue("Body")

	// Respond with TwiML. An empty <Response> acknowledges receipt without
	// sending a reply message; add <Message> elements here as needed.
	w.Header().Set("Content-Type", "text/xml")
	fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<Response></Response>\n")

	// Log the verified inbound message. Replace with your application logic.
	_ = from
	_ = body
}
