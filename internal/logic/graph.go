package logic

import (
	"context"
	"github.com/quanxiang-cloud/faas/pkg/config"
)

type Graph interface {
	Get(c context.Context, r *GetGraphRequest) (*GetGraphResponse, error)
}

type graph struct {
	conf config.Config
}

func NewGraph(c context.Context, conf config.Config) Graph {
	return &graph{
		conf: conf,
	}
}

type GetGraphRequest struct {
}
type GetGraphResponse struct {
	Data interface{} `json:"data"`
}

func (g *graph) Get(c context.Context, r *GetGraphRequest) (*GetGraphResponse, error) {

	return &GetGraphResponse{
		Data: g.conf.Graph,
	}, nil

}
