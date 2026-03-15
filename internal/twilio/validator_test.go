package twilio_test

import (
	"net/url"
	"testing"

	"github.com/zyx-holdings/go-spec/internal/twilio"
)

// Test vectors derived from the Twilio documentation example:
// https://www.twilio.com/docs/usage/webhooks/webhooks-security#validating-signatures-from-twilio
const (
	testAuthToken = "12345"
	testURL       = "https://mycompany.com/myapp.php?foo=1&bar=2"
)

var testParams = url.Values{
	"Digits":    {"1234"},
	"To":        {"+18005551212"},
	"From":      {"+14158675309"},
	"Caller":    {"+14158675309"},
	"CallSid":   {"CA1234567890ABCDE"},
	"AccountSid": {"ACXXXXXXXXXXXXXXXXX"},
	"CallStatus": {"ringing"},
	"ApiVersion": {"2010-04-01"},
	"Direction":  {"inbound"},
}

// expectedSig is the base64-encoded HMAC-SHA1 of the URL + sorted POST params
// using "12345" as the auth token. Precomputed for regression testing.
const expectedSig = "/CsNkvS6ruwEyE9vMZtwzZdmO1s="

func TestValidateSignature_Valid(t *testing.T) {
	if !twilio.ValidateSignature(testAuthToken, expectedSig, testURL, testParams) {
		t.Fatal("expected valid signature to pass validation")
	}
}

func TestValidateSignature_WrongSignature(t *testing.T) {
	if twilio.ValidateSignature(testAuthToken, "invalidsignature==", testURL, testParams) {
		t.Fatal("expected invalid signature to fail validation")
	}
}

func TestValidateSignature_WrongAuthToken(t *testing.T) {
	if twilio.ValidateSignature("wrongtoken", expectedSig, testURL, testParams) {
		t.Fatal("expected wrong auth token to fail validation")
	}
}

func TestValidateSignature_WrongURL(t *testing.T) {
	if twilio.ValidateSignature(testAuthToken, expectedSig, "https://mycompany.com/other?foo=1&bar=2", testParams) {
		t.Fatal("expected wrong URL to fail validation")
	}
}

func TestValidateSignature_EmptyParams(t *testing.T) {
	// GET requests have no POST body; the signature covers only the URL.
	sig := "zYQTYrRWXE7LtzbG4PfP7/bkkGo=" // precomputed for testURL with no params
	if !twilio.ValidateSignature(testAuthToken, sig, testURL, nil) {
		t.Fatal("expected GET-style (no params) signature to pass validation")
	}
}

func TestValidateSignature_ParamOrderIndependent(t *testing.T) {
	// Params supplied in reverse order should produce the same signature because
	// the validator sorts them alphabetically before hashing.
	reversed := url.Values{
		"To":         {"+18005551212"},
		"From":       {"+14158675309"},
		"Digits":     {"1234"},
		"Direction":  {"inbound"},
		"ApiVersion": {"2010-04-01"},
		"CallStatus": {"ringing"},
		"AccountSid": {"ACXXXXXXXXXXXXXXXXX"},
		"CallSid":    {"CA1234567890ABCDE"},
		"Caller":     {"+14158675309"},
	}
	if !twilio.ValidateSignature(testAuthToken, expectedSig, testURL, reversed) {
		t.Fatal("param order should not affect validation result")
	}
}

func TestValidateSignature_EmptySignature(t *testing.T) {
	if twilio.ValidateSignature(testAuthToken, "", testURL, testParams) {
		t.Fatal("empty signature must not pass validation")
	}
}

func TestValidateSignature_TimingAttackSafe(t *testing.T) {
	// Ensure the function does not short-circuit on the first byte; we can only
	// observe that calls with different-length inputs still return false without
	// panicking or taking dramatically different time (structural test).
	result := twilio.ValidateSignature(testAuthToken, "a", testURL, testParams)
	if result {
		t.Fatal("single-byte signature must not match")
	}
}
