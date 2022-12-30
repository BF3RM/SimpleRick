package sentry

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

var (
	ErrInvalidHTTPMethod      = errors.New("invalid HTTP Method")
	ErrMissingSignatureHeader = errors.New("missing signature header")
	ErrMissingResourceHeader  = errors.New("missing resource header")
	ErrParsingPayload         = errors.New("error parsing payload")
	ErrHMACVerificationFailed = errors.New("HMAC verification failed")
)

type Event string

const (
	InstallationEvent   Event = "installation"
	UninstallationEvent Event = "uninstallation"
	IssueAlertEvent     Event = "event_alert"
	MetricAlertEvent    Event = "metric_alert"
	IssueEvent          Event = "issue"
	ErrorEvent          Event = "error"
)

type EventAction string

const (
	InstallationCreatedAction EventAction = "created"
	InstallationDeletedAction EventAction = "deleted"
	IssueAlertTriggeredAction EventAction = "triggered"
	MetricAlertCriticalAction EventAction = "critical"
	MetricAlertWarningAction  EventAction = "warning"
	MetricAlertResolvedAction EventAction = "resolved"
	IssueCreatedAction        EventAction = "created"
	IssueResolvedAction       EventAction = "resolved"
	IssueAssignedAction       EventAction = "assigned"
	IssueIgnoredAction        EventAction = "ignored"
	ErrorCreatedAction        EventAction = "created"
)

func ValidatePayload(req *http.Request, secret []byte) (payload []byte, err error) {
	if req.Method != http.MethodPost {
		return nil, ErrInvalidHTTPMethod
	}

	payload, err = ioutil.ReadAll(req.Body)
	if err != nil || len(payload) == 0 {
		return nil, ErrParsingPayload
	}

	if len(secret) > 0 {
		signature := req.Header.Get("Sentry-Hook-Signature")
		if len(signature) == 0 {
			return nil, ErrMissingSignatureHeader
		}
		mac := hmac.New(sha256.New, secret)
		_, _ = mac.Write(payload)

		digest := hex.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(signature), []byte(digest)) {
			return nil, ErrHMACVerificationFailed
		}
	}

	return payload, nil
}

func WebhookResource(req *http.Request) string {
	return req.Header.Get("Sentry-Hook-Resource")
}

func ParseWebhook(resource string, payload []byte) (EventAction, interface{}, error) {
	if len(resource) == 0 {
		return EventAction(""), nil, ErrMissingResourceHeader
	}

	ctx := new(WebhookPayload)
	ctx.Event = Event(resource)

	if err := json.Unmarshal(payload, &ctx); err != nil {
		return EventAction(""), nil, err
	}

	return ctx.Action, ctx.Data, nil
}
