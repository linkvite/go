package linkvite

import (
	"errors"
	"fmt"
	"net/http"
)

// Error is a failed request. Status is 0 for network errors. Quote RequestID
// when contacting support.
type Error struct {
	Status  int
	Message string
	// Errors maps field -> problem for validation errors (400).
	Errors    map[string]string
	RequestID string
	Method    string
	URL       string
	Body      []byte
	cause     error
}

func (e *Error) Error() string {
	s := fmt.Sprintf("linkvite: %s %s: %d %s", e.Method, e.URL, e.Status, e.Message)
	if e.RequestID != "" {
		s += " (request " + e.RequestID + ")"
	}
	return s
}

func (e *Error) Unwrap() error { return e.cause }

// IsNotFound, IsUnauthorized, IsForbidden, IsRateLimited and IsProRequired
// report whether err is an *Error with that status.
func IsNotFound(err error) bool     { return hasStatus(err, http.StatusNotFound) }
func IsUnauthorized(err error) bool { return hasStatus(err, http.StatusUnauthorized) }
func IsForbidden(err error) bool    { return hasStatus(err, http.StatusForbidden) }
func IsRateLimited(err error) bool  { return hasStatus(err, http.StatusTooManyRequests) }
func IsProRequired(err error) bool  { return hasStatus(err, http.StatusPaymentRequired) }

func hasStatus(err error, status int) bool {
	var e *Error
	return errors.As(err, &e) && e.Status == status
}
