// Package twilio provides utilities for validating Twilio webhook requests.
// See: https://www.twilio.com/docs/usage/webhooks/webhooks-security
package twilio

import (
	"crypto/hmac"
	"crypto/sha1" //nolint:gosec // Twilio mandates HMAC-SHA1 per their spec
	"encoding/base64"
	"net/url"
	"sort"
)

// ComputeSignature returns the base64-encoded HMAC-SHA1 signature for the given
// URL and POST parameters. This is exported primarily for use in tests that need
// to build valid signed requests.
func ComputeSignature(authToken, rawURL string, params url.Values) string {
	return computeSignature(authToken, rawURL, params)
}

// ValidateSignature verifies that an incoming request originated from Twilio.
//
// Parameters:
//   - authToken: the Twilio account auth token used as the HMAC-SHA1 secret key
//   - signature: the value of the X-Twilio-Signature header from the request
//   - rawURL: the full request URL exactly as it was configured in Twilio (scheme + host + path + query)
//   - params: the POST form parameters from the request body (empty for GET requests)
//
// Returns true when the computed signature matches the provided header value.
func ValidateSignature(authToken, signature, rawURL string, params url.Values) bool {
	expected := computeSignature(authToken, rawURL, params)
	// Use hmac.Equal for constant-time comparison to prevent timing attacks.
	return hmac.Equal([]byte(expected), []byte(signature))
}

// computeSignature builds the HMAC-SHA1 signature Twilio expects.
//
// Algorithm (from Twilio docs):
//  1. Take the full URL of the request URL
//  2. If it is a POST request, sort all POST parameters alphabetically
//     and append each name+value (no delimiters) to the URL string
//  3. Sign with HMAC-SHA1 using the auth token as the key
//  4. Base64-encode the result
func computeSignature(authToken, rawURL string, params url.Values) string {
	s := rawURL

	if len(params) > 0 {
		// Sort parameter names alphabetically.
		keys := make([]string, 0, len(params))
		for k := range params {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			s += k + params.Get(k)
		}
	}

	mac := hmac.New(sha1.New, []byte(authToken)) //nolint:gosec // required by Twilio spec
	mac.Write([]byte(s))                          //nolint:errcheck // hash.Hash.Write never returns an error
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
