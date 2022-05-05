package event

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
)

type EventAdaptor struct {
	handleMapping map[EventType][]Handle

	rg  *gin.RouterGroup
	log logger.AdaptedLogger
}

type Handle func(*MsgBus) error

func New(log logger.AdaptedLogger, opt ...Option) *EventAdaptor {
	e := &EventAdaptor{
		handleMapping: map[EventType][]Handle{},
		log:           log,
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
		daprEvent := &DaprEvent{}
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK, nil)
		}
		ea.log.Debugf("receiving event \n %s", string(b))
		if err := json.Unmarshal(b, daprEvent); err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK, nil)
		}

		bus := &MsgBus{
			CTX: context.Background(),
			Msg: daprEvent.Data,
		}
		var h []Handle
		switch {
		case bus.Msg.Fn != nil:
			ea.log.Debugf("handle function event")
			bus.Type = Function
			h = ea.handleMapping[Function]
		case bus.Msg.Pr != nil:
			ea.log.Debugf("handle pipeline event")
			bus.Type = APIDoc
			h = ea.handleMapping[APIDoc]
		case bus.Msg.Svc != nil:
			ea.log.Debugf("handle pipeline event")
			bus.Type = Serving
			h = ea.handleMapping[Serving]
		}

		if err := do(bus, h); err != nil {
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

func Convert(bus *MsgBus) *BaseMessage {
	switch bus.Type {
	case Function:
		return &BaseMessage{
			Name:  bus.Msg.Fn.Name,
			State: bus.Msg.Fn.State,
			Topic: bus.Msg.Fn.Topic,
		}
	case APIDoc:
		return &BaseMessage{
			Name:  bus.Msg.Pr.Name,
			State: bus.Msg.Pr.State,
			Topic: bus.Msg.Pr.Topic,
		}
	default:
		return nil
	}
}
