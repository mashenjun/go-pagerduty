package pagerduty

import (
	"bytes"
	"encoding/json"
	"errors"
)

const ChangeEventPath = "/v2/change/enqueue"

// ChangeEvent represents a ChangeEvent's request parameters
// https://developer.pagerduty.com/docs/events-api-v2/send-change-events/#parameters
type ChangeEvent struct {
	RoutingKey string  `json:"routing_key"`
	Payload    Payload `json:"payload"`
	Links      []Link  `json:"links"`
}

// Payload ChangeEvent Payload
// https://developer.pagerduty.com/docs/events-api-v2/send-change-events/#example-request-payload
type Payload struct {
	Source        string            `json:"source"`
	Summary       string            `json:"summary"`
	Timestamp     string            `json:"timestamp"`
	CustomDetails map[string]string `json:"custom_details"`
}

// Link represents a single link in a ChangeEvent
// https://developer.pagerduty.com/docs/events-api-v2/send-change-events/#the-links-property
type Link struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

// ChangeEventResponse is the json response body for an event
type ChangeEventResponse struct {
	Status  string   `json:"status,omitempty"`
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// CreateChangeEvent Sends PagerDuty a single ChangeEvent to record
// The v2EventsAPIEndpoint parameter must be set on the client
// Documentation can be found at https://developer.pagerduty.com/docs/events-api-v2/send-change-events
func (c *Client) CreateChangeEvent(e ChangeEvent) (*ChangeEventResponse, error) {
	if c.v2EventsAPIEndpoint == "" {
		return nil, errors.New("v2EventsAPIEndpoint field must be set on Client")
	}

	headers := make(map[string]string)

	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	resp, err := c.doWithEndpoint(
		c.v2EventsAPIEndpoint,
		"POST",
		ChangeEventPath,
		false,
		bytes.NewBuffer(data),
		&headers,
	)
	if err != nil {
		return nil, err
	}

	var eventResponse ChangeEventResponse
	if err := json.NewDecoder(resp.Body).Decode(&eventResponse); err != nil {
		return nil, err
	}

	return &eventResponse, nil
}
