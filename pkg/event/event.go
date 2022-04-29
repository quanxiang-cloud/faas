package event

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EventAdaptor struct {
	handleMapping map[EventType][]Handle

	rg *gin.RouterGroup
}

type Handle func(*MsgBus) error

func New(opt ...Option) *EventAdaptor {
	e := &EventAdaptor{
		handleMapping: map[EventType][]Handle{},
	}
	for _, f := range opt {
		f(e)
	}
	e.init()
	return e
}

type Option func(*EventAdaptor)

func WithHandle(t EventType, h ...Handle) Option {
	return func(ea *EventAdaptor) {
		ea.handleMapping[t] = h
	}
}

func WithRouter(group *gin.RouterGroup) Option {
	return func(ea *EventAdaptor) {
		ea.rg = group
	}
}

func (ea *EventAdaptor) init() {
	ea.rg.POST("/event", func(c *gin.Context) {
		msg := &MsgBus{
			Message: &Message{},
			CTX:     context.Background(),
		}

		if err := c.ShouldBind(msg.Message); err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK, nil)
		}

		var h []Handle
		switch {
		case msg.FnMessage != nil:
			msg.Type = Function
			h = ea.handleMapping[Function]
		case msg.APIDocMessage != nil:
			msg.Type = APIDoc
			h = ea.handleMapping[APIDoc]
		}

		if err := do(msg, h); err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK, nil)
		}
	})
}

func do(msg *MsgBus, handle []Handle) error {
	for _, fn := range handle {
		if err := fn(msg); err != nil {
			return err
		}
	}
	return nil
}

func Convert(msg *MsgBus) *BaseMessage {
	switch msg.Type {
	case Function:
		return &BaseMessage{
			Name:  msg.FnMessage.Name,
			State: msg.FnMessage.State,
			Topic: msg.FnMessage.Topic,
		}
	case APIDoc:
		return &BaseMessage{
			Name:  msg.APIDocMessage.Name,
			State: msg.APIDocMessage.State,
			Topic: msg.APIDocMessage.Topic,
		}
	default:
		return nil
	}
}
