package sms_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/zyx-holdings/go-spec/internal/twilio"
	sms "github.com/zyx-holdings/go-spec/server/api/sms"
)

const (
	testAuthToken  = "test-auth-token-abc123"
	testWebhookURL = "https://example.com/api/sms/incoming"
)

var baseParams = url.Values{
	"From": {"+14155551234"},
	"To":   {"+18005550000"},
	"Body": {"Hello"},
}

// buildRequest constructs a signed or unsigned POST request with form params.
func buildRequest(t *testing.T, params url.Values, sign bool) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, testWebhookURL, strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if sign {
		req.Header.Set("X-Twilio-Signature", twilio.ComputeSignature(testAuthToken, testWebhookURL, params))
	}
	return req
}

func TestIncomingHandler_ValidRequest(t *testing.T) {
	h := sms.NewHandler(testAuthToken, testWebhookURL)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, buildRequest(t, baseParams, true))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/xml" {
		t.Errorf("expected Content-Type text/xml, got %q", ct)
	}
	if !strings.Contains(w.Body.String(), "<Response>") {
		t.Errorf("expected TwiML <Response>, got: %s", w.Body.String())
	}
}

func TestIncomingHandler_MissingSignature(t *testing.T) {
	h := sms.NewHandler(testAuthToken, testWebhookURL)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, buildRequest(t, baseParams, false))

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for missing signature, got %d", w.Code)
	}
}

func TestIncomingHandler_InvalidSignature(t *testing.T) {
	h := sms.NewHandler(testAuthToken, testWebhookURL)
	req := buildRequest(t, baseParams, false)
	req.Header.Set("X-Twilio-Signature", "invalidsignature==")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for invalid signature, got %d", w.Code)
	}
}

func TestIncomingHandler_WrongAuthToken(t *testing.T) {
	// Handler uses a different auth token than the one used to sign the request.
	h := sms.NewHandler("wrong-token", testWebhookURL)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, buildRequest(t, baseParams, true)) // signed with testAuthToken

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for wrong auth token, got %d", w.Code)
	}
}

func TestIncomingHandler_WrongURL(t *testing.T) {
	// Handler configured with a URL different from the one used to sign.
	h := sms.NewHandler(testAuthToken, "https://other.example.com/api/sms/incoming")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, buildRequest(t, baseParams, true)) // signed with testWebhookURL

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for mismatched URL, got %d", w.Code)
	}
}

func TestIncomingHandler_MethodNotAllowed(t *testing.T) {
	h := sms.NewHandler(testAuthToken, testWebhookURL)
	req := httptest.NewRequest(http.MethodGet, testWebhookURL, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 for GET request, got %d", w.Code)
	}
}
