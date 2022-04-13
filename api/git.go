package restful

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	error2 "github.com/quanxiang-cloud/cabin/error"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/faas/internal/logic"
	"github.com/quanxiang-cloud/faas/pkg/code"
	"github.com/quanxiang-cloud/faas/pkg/config"
)

// Git git api
type Git struct {
	git logic.Git
}

// NewGitAPI new
func NewGitAPI(c context.Context, conf *config.Config, db *gorm.DB) *Git {
	return &Git{
		git: logic.NewGit(c, db),
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
	resp.Format(f.git.Create(ginheader.MutateContext(c), r)).Context(c)
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
