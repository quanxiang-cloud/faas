package event

import "context"

type EventType int

const (
	Function EventType = iota
	APIDoc
	Serving
)

type DaprEvent struct {
	Topic           string   `json:"topic"`
	Pubsubname      string   `json:"pubsubname"`
	Traceid         string   `json:"traceid"`
	ID              string   `json:"id"`
	Datacontenttype string   `json:"datacontenttype"`
	Data            *Message `json:"data"`
	Type            string   `json:"type"`
	Specversion     string   `json:"specversion"`
	Source          string   `json:"source"`
}

type MsgBus struct {
	Msg  *Message
	CTX  context.Context
	Type EventType
	Data string
}

type Message struct {
	Fn *Fn `json:"fn"`
	Pr *Pr `json:"pr"`
}

type BaseMessage struct {
	Name  string `json:"name,omitempty"`
	Topic string `json:"topic,omitempty"`
	State string `json:"state,omitempty"`
}

type Fn struct {
	Name        string `json:"name,omitempty"`
	Topic       string `json:"topic,omitempty"`
	State       string `json:"state,omitempty"`
	ResourceRef string `json:"resource_ref,omitempty"`
}

type Pr struct {
	Name  string `json:"name,omitempty"`
	Topic string `json:"topic,omitempty"`
	State string `json:"state,omitempty"`
}
