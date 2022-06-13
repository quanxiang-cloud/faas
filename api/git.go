package restful

import (
	"context"

	"github.com/quanxiang-cloud/faas/pkg/basic/k8s"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	error2 "github.com/quanxiang-cloud/cabin/error"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/faas/internal/logic"
	"github.com/quanxiang-cloud/faas/pkg/basic/define/code"
	"github.com/quanxiang-cloud/faas/pkg/config"
)

// Git git api
type Git struct {
	git  logic.Git
	k8sc k8s.Client
}

// NewGitAPI new
func NewGitAPI(c context.Context, conf *config.Config, db *gorm.DB, kc k8s.Client) *Git {
	return &Git{
		git:  logic.NewGit(c, db),
		k8sc: kc,
	}
}

// Create create
func (f *Git) Create(c *gin.Context) {
	r := &logic.CreateGitRequest{}
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	res, err := f.git.Create(ginheader.MutateContext(c), r)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	// create git to k8s
	err = f.k8sc.CreateGitSSH(ginheader.MutateContext(c), r.KnownHosts, r.KeyScanKnownHosts, r.SSH)
	if err != nil {
		resp.Format(res, err).Context(c)
		return
	}
	resp.Format(res, nil).Context(c)
}

// Update update
func (f *Git) Update(c *gin.Context) {
	r := &logic.UpdateGitRequest{}
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	resp.Format(f.git.Update(ginheader.MutateContext(c), r)).Context(c)
}

// Delete delete
func (f *Git) Delete(c *gin.Context) {
	r := &logic.DeleteGitRequest{}
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	resp.Format(f.git.Delete(ginheader.MutateContext(c), r)).Context(c)
}

// Get get
func (f *Git) Get(c *gin.Context) {
	r := &logic.GetGitRequest{}
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	resp.Format(f.git.Get(ginheader.MutateContext(c), r)).Context(c)
}
