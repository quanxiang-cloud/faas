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

// GroupAPI GroupAPI
type GroupAPI struct {
	groupService logic.GroupService
}

// NewGroupAPI NewGroupAPI
func NewGroupAPI(c context.Context, conf *config.Config, db *gorm.DB) *GroupAPI {
	return &GroupAPI{
		groupService: logic.NewGroupService(c, db, conf),
	}
}

// Create Create
func (g *GroupAPI) Create(c *gin.Context) {
	req := &logic.CreateGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	ctx := ginheader.MutateContext(c)
	req.UserID = c.GetHeader(_userID)
	resp.Format(g.groupService.CreateGroup(ctx, req)).Context(c)
}

// BindingGroup BindingGroup
func (g *GroupAPI) BindingGroup(c *gin.Context) {
	req := &logic.BindingGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	req.GroupID = c.Param("groupID")
	req.UserID = c.GetHeader(_userID)
	ctx := ginheader.MutateContext(c)
	resp.Format(g.groupService.BindingGroup(ctx, req)).Context(c)
}

// CheckGroup CheckGroup
func (g *GroupAPI) CheckGroup(c *gin.Context) {
	req := &logic.CheckGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	ctx := ginheader.MutateContext(c)
	resp.Format(g.groupService.CheckGroup(ctx, req)).Context(c)
}

// CheckMember CheckMember
func (g *GroupAPI) CheckMember(c *gin.Context) {
	req := &logic.CheckMemberReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	req.UserID = c.GetHeader(_userID)
	ctx := ginheader.MutateContext(c)
	resp.Format(g.groupService.CheckMember(ctx, req)).Context(c)
}
