package restful

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	error2 "github.com/quanxiang-cloud/cabin/error"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/faas/internal/logic"
	"github.com/quanxiang-cloud/faas/pkg/basic/define/code"
	"github.com/quanxiang-cloud/faas/pkg/basic/k8s"
	"github.com/quanxiang-cloud/faas/pkg/config"
	"gorm.io/gorm"
)

// Function Function api
type Function struct {
	fn logic.Function
	ps logic.PubSub
}

// NewFunctionAPI new
func NewFunctionAPI(c context.Context, conf *config.Config, db *gorm.DB, kc k8s.Client, rc redis.UniversalClient, esClient *elastic.Client) *Function {
	return &Function{
		fn: logic.NewFunction(c, db, *conf, kc, esClient),
		ps: logic.NewPubSub(c, rc),
	}
}

// Create create
func (f *Function) Create(c *gin.Context) {
	r := &logic.CreateFunctionRequest{
		GroupID:   c.Param("groupID"),
		CreatedBy: c.GetHeader("User-Id"),
	}
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	res, err := f.fn.Create(ginheader.MutateContext(c), r)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	buildFunctionRequest := &logic.BuildFunctionRequest{ID: res.ID}
	_, err = f.fn.Build(ginheader.MutateContext(c), buildFunctionRequest)
	if err != nil {
		resp.Format(res, err).Context(c)
	}
	resp.Format(res, nil).Context(c)
}

// Delete delete
func (f *Function) Delete(c *gin.Context) {
	r := &logic.DeleteFunctionRequest{
		FunctionID: c.Param("functionID"),
	}
	if r.FunctionID == "" {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	resp.Format(f.fn.Delete(ginheader.MutateContext(c), r)).Context(c)
}

// Get get
func (f *Function) Get(c *gin.Context) {
	r := &logic.GetFunctionRequest{
		ID: c.Param("functionID"),
	}
	err := c.ShouldBind(r)
	if err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	resp.Format(f.fn.Get(ginheader.MutateContext(c), r)).Context(c)
}

// ListLog ListLog
func (f *Function) ListLog(c *gin.Context) {
	ctx := c.Request.Context()
	req := &logic.ListLogRequest{
		ResourceRef: c.Param("resourceRef"),
	}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	resp.Format(f.fn.ListLog(ctx, req)).Context(c)
}

// List List
func (f *Function) List(c *gin.Context) {
	ctx := c.Request.Context()
	req := &logic.ListRequest{
		GroupID:   c.Param("groupID"),
		ProjectID: c.Param("projectID"),
	}
	resp.Format(f.fn.List(ctx, req)).Context(c)
}

func (f *Function) RegSwagger(c *gin.Context) {
	req := &logic.RegSwaggerReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	req.GroupID = c.Param("groupID")

	resp.Format(f.fn.RegSwagger(ginheader.MutateContext(c), req)).Context(c)
}

func (f *Function) UpdateDescribe(c *gin.Context) {
	req := &logic.UpdateFuncDescribeReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	req.ID = c.Param("functionID")
	data, err := f.fn.UpdateDescribe(ginheader.MutateContext(c), req)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	resp.Format(data, nil).Context(c)
}
