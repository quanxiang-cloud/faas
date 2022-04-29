package event

import "context"

type EventType int

const (
	Function EventType = iota
	APIDoc
	Serving
)

type MsgBus struct {
	*Message
	CTX  context.Context
	Type EventType
	Data string
}

type Message struct {
	*FnMessage
	*APIDocMessage
}

type BaseMessage struct {
	Name  string `json:"name,omitempty"`
	Topic string `json:"topic,omitempty"`
	State string `json:"state,omitempty"`
}

type FnMessage struct {
	Name        string `json:"name,omitempty"`
	Topic       string `json:"topic,omitempty"`
	State       string `json:"state,omitempty"`
	ResourceRef string `json:"resource_ref,omitempty"`
}

type APIDocMessage struct {
	Name  string `json:"name,omitempty"`
	Topic string `json:"topic,omitempty"`
	State string `json:"state,omitempty"`
}
