package restful

import (
	"github.com/gin-gonic/gin"
	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/faas/internal/logic"
	"github.com/quanxiang-cloud/faas/pkg/code"
	"github.com/quanxiang-cloud/faas/pkg/config"
	"github.com/quanxiang-cloud/faas/pkg/k8s"
	"gorm.io/gorm"
)

type Serving struct {
	svc logic.Serving
}

func NewServing(db *gorm.DB, config *config.Config, kc k8s.Client) *Serving {
	svc := logic.NewServing(db, config)
	return &Serving{
		svc: svc,
	}
}

func (s *Serving) serve(c *gin.Context) {
	req := &logic.ServeReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	resp.Format(s.svc.Serve(header.MutateContext(c), req)).Context(c)
}

func (s *Serving) offline(c *gin.Context) {
	req := &logic.OffLineReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	resp.Format(s.svc.OffLine(header.MutateContext(c), req)).Context(c)
}
