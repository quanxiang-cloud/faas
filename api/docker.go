package restful

import (
	"context"
	"github.com/quanxiang-cloud/faas/pkg/k8s"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	error2 "github.com/quanxiang-cloud/cabin/error"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/faas/internal/logic"
	"github.com/quanxiang-cloud/faas/pkg/code"
	"github.com/quanxiang-cloud/faas/pkg/config"
)

// Docker docker api
type Docker struct {
	dokcer logic.Docker
	k8sc   k8s.Client
}

// NewDockerAPI new
func NewDockerAPI(c context.Context, conf *config.Config, db *gorm.DB, kc k8s.Client) *Docker {
	return &Docker{
		dokcer: logic.NewDocker(c, db),
		k8sc:   kc,
	}
}

// Create create
func (f *Docker) Create(c *gin.Context) {
	r := &logic.CreateDockerRequest{}
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	res, err := f.dokcer.Create(ginheader.MutateContext(c), r)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	err = f.k8sc.CreateDocker(ginheader.MutateContext(c), r.Host, r.UserName, r.Secret)
	if err != nil {
		resp.Format(res, err).Context(c)
	}
	resp.Format(res, nil).Context(c)
}

// Update update
func (f *Docker) Update(c *gin.Context) {
	r := &logic.UpdateDockerRequest{}
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	resp.Format(f.dokcer.Update(ginheader.MutateContext(c), r)).Context(c)
}

// Delete delete
func (f *Docker) Delete(c *gin.Context) {
	r := &logic.DeleteDockerRequest{}
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	resp.Format(f.dokcer.Delete(ginheader.MutateContext(c), r)).Context(c)
}

// Get get
func (f *Docker) Get(c *gin.Context) {
	r := &logic.GetDockerRequest{}
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	resp.Format(f.dokcer.Get(ginheader.MutateContext(c), r)).Context(c)
}
