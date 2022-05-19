package publish

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

const (
	messageURL = "http://message/api/v1/message/send"

	faas    = "faas"
	channel = "letter"
)

type message struct {
	client http.Client
}

func NewMessage() *message {
	return &message{client: client.New(client.Config{
		Timeout:      time.Minute,
		MaxIdleConns: 10,
	})}
}

func (m *message) Publish(ctx context.Context, req *PublishReq) (*PublishResp, error) {
	content := map[string]interface{}{
		"type":    faas,
		"content": req.Content,
	}
	contentByte, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	resp := new(PublishResp)
	err = client.POST(ctx, &m.client, messageURL, map[string]interface{}{
		channel: LetterSpec{
			ID:      req.UserID,
			UUID:    []string{req.UUID},
			Content: contentByte,
		},
	}, resp)
	return resp, err
}
