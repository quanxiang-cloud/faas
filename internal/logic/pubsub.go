package logic

import (
	"context"
	"log"
	"time"

	"github.com/quanxiang-cloud/faas/internal/models"
	re "github.com/quanxiang-cloud/faas/internal/models/redis"
	"github.com/quanxiang-cloud/faas/pkg/event"
	"github.com/quanxiang-cloud/faas/pkg/publish"

	"github.com/go-redis/redis/v8"
)

type PubSub interface {
	Subscribe(ctx context.Context, req *SubscribeReq) (*SubscribeResp, error)
	Publish(msgBus *event.MsgBus) error
}

var ttl = time.Duration(300) * time.Second

func NewPubSub(ctx context.Context, client redis.UniversalClient) PubSub {
	return &pubSub{
		ps: re.NewPubSub(client),
		ph: publish.NewMessage(),
	}
}

type pubSub struct {
	ps models.PubSubRepo
	ph publish.Publish
}

type SubscribeReq struct {
	UserID string `json:"userID"`
	UUID   string `json:"uuid"`
	Topic  string `json:"topic"`
	Key    string `json:"key"`
}

type SubscribeResp struct{}

func (p *pubSub) Subscribe(ctx context.Context, req *SubscribeReq) (*SubscribeResp, error) {
	err := p.ps.Subscribe(ctx, &models.PubSub{
		UserID: req.UserID,
		Topic:  req.Topic,
		UUID:   req.UUID,
		Key:    req.Key,
	}, ttl)

	return &SubscribeResp{}, err
}

type PublishReq struct {
	Topic string
	Key   string
}

func (p *pubSub) Publish(msg *event.MsgBus) error {
	b := event.Convert(msg)

	consumers, err := p.ps.Get(msg.CTX, b.Topic, msg.Data)
	if err != nil {
		return err
	}

	for _, consumer := range consumers {
		_, err := p.ph.Publish(msg.CTX, &publish.PublishReq{
			UserID: consumer.UserID,
			UUID:   consumer.UUID,
			Content: map[string]string{
				"topic": consumer.Topic,
				"key":   consumer.Key,
			},
		})
		if err != nil {
			return nil
		}
		log.Printf("push letter: %#v\n", consumer)
	}
	return nil
}
