package publish

import "context"

type Publish interface {
	Publish(context.Context, *PublishReq) (*PublishResp, error)
}

type PublishReq struct {
	UserID  string
	UUID    string
	Content interface{}
}

type PublishResp struct{}

type LetterSpec struct {
	ID      string   `json:"id,omitempty"`
	UUID    []string `json:"uuid,omitempty"`
	Content []byte   `json:"content,omitempty"`
}
