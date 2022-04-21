package restful

import (
	"context"
	"github.com/gin-gonic/gin"
	error2 "github.com/quanxiang-cloud/cabin/error"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/faas/internal/logic"
	"github.com/quanxiang-cloud/faas/pkg/code"
	"github.com/quanxiang-cloud/faas/pkg/config"
	"gorm.io/gorm"
)

// ProjectAPI ProjectAPI
type ProjectAPI struct {
	projectService logic.ProjectService
}

// NewProjectAPI NewProjectAPI
func NewProjectAPI(ctx context.Context, conf *config.Config, db *gorm.DB) *ProjectAPI {
	return &ProjectAPI{
		projectService: logic.NewProjectService(ctx, db, conf),
	}
}

// CreateProject CreateProject
func (p *ProjectAPI) CreateProject(c *gin.Context) {
	req := &logic.CreateProjectReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}

	ctx := ginheader.MutateContext(c)
	req.UserID = c.GetHeader(_userID)
	req.GroupID = c.Param("groupID")
	resp.Format(p.projectService.CreateProject(ctx, req)).Context(c)
}

// GetProjectByID GetProjectByID
func (p *ProjectAPI) GetProjectByID(c *gin.Context) {
	req := &logic.GetProjectByIDReq{}
	ctx := ginheader.MutateContext(c)
	req.ProjectID = c.Param(":projectID")
	resp.Format(p.projectService.GetProjectByID(ctx, req)).Context(c)
}

// GetList GetList
func (p *ProjectAPI) GetList(c *gin.Context) {
	req := &logic.GetListReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	req.GroupID = c.Param("groupID")
	ctx := ginheader.MutateContext(c)
	resp.Format(p.projectService.GetList(ctx, req)).Context(c)
}

// DelProject DelProject
func (p *ProjectAPI) DelProject(c *gin.Context) {
	req := &logic.DelProjectReq{}
	req.ProjectID = c.Param("projectID")
	ctx := ginheader.MutateContext(c)
	resp.Format(p.projectService.DelProject(ctx, req)).Context(c)
}

// UpdDescribe UpdDescribe
func (p *ProjectAPI) UpdDescribe(c *gin.Context) {
	req := &logic.UpdateDescribeReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	req.ProjectID = c.Param("projectID")
	ctx := ginheader.MutateContext(c)
	resp.Format(p.projectService.UpdateDescribe(ctx, req)).Context(c)
}
