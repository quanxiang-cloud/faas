package restful

import (
	"context"
	"github.com/gin-gonic/gin"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/faas/internal/logic"
	"github.com/quanxiang-cloud/faas/pkg/config"
)

// GraphAPI graphAPI
type GraphAPI struct {
	graph logic.Graph
}

// NewGraphAPI NewProjectAPI
func NewGraphAPI(ctx context.Context, conf *config.Config) *GraphAPI {
	return &GraphAPI{
		graph: logic.NewGraph(ctx, *conf),
	}
}

func (g *GraphAPI) Get(c *gin.Context) {
	r := &logic.GetGraphRequest{}

	resp.Format(g.graph.Get(ginheader.MutateContext(c), r)).Context(c)
}
