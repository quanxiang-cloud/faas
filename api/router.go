package restful

import (
	"context"

	elastic2 "github.com/quanxiang-cloud/cabin/tailormade/db/elastic"
	mysql2 "github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"

	"github.com/quanxiang-cloud/faas/pkg/event"
	"github.com/quanxiang-cloud/faas/pkg/k8s"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	ginlogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
	"github.com/quanxiang-cloud/faas/pkg/config"
	"github.com/quanxiang-cloud/faas/pkg/probe"
)

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
)

// Router router
type Router struct {
	c      *config.Config
	engine *gin.Engine
}

// NewRouter 开启路由
func NewRouter(ctx context.Context, c *config.Config, log logger.AdaptedLogger) (*Router, error) {
	redisClient, err := redis2.NewClient(c.Redis)
	if err != nil {
		return nil, err
	}
	esClient, err := elastic2.NewClient(&c.Elastic, log)
	if err != nil {
		return nil, err
	}
	db, err := mysql2.New(c.Mysql, log)
	if err != nil {
		return nil, err
	}
	k8sClient := k8s.NewClient(c.K8s.NameSpace)

	engine, err := newRouter(c)
	if err != nil {
		return nil, err
	}

	v1 := engine.Group("/api/v1/faas")
	graphAPI := NewGraphAPI(ctx, c)
	graph := v1.Group("/graph")
	{
		graph.GET("", graphAPI.Get)
	}
	{
		gitAPI := NewGitAPI(ctx, c, db, k8sClient)
		g := v1.Group("/git")
		{
			g.POST("", gitAPI.Create)
			g.PUT("/:id", gitAPI.Update)
			g.DELETE("/:id", gitAPI.Delete)
			g.GET("", gitAPI.Get)
		}
		dockerAPI := NewDockerAPI(ctx, c, db, k8sClient)
		d := v1.Group("/docker")
		{
			d.POST("", dockerAPI.Create)
			d.PUT("/:id", dockerAPI.Update)
			d.DELETE("/:id", dockerAPI.Delete)
			d.GET("", dockerAPI.Get)
		}
		cm := NewCompoundAPI(ctx, redisClient)
		cmGroup := v1.Group("/cm")
		{
			cmGroup.POST("/subscribe", cm.Subscribe)
		}
	}
	{
		userAPI := NewUserAPI(ctx, c, db)
		groupAPI := NewGroupAPI(ctx, c, db)
		projectAPI := NewProjectAPI(ctx, c, db)
		user := v1.Group("/user")
		{
			user.POST("", userAPI.CreateUser)
		}
		g := v1.Group("")
		{
			g.POST("/group", groupAPI.Create)
			g.GET("/groups", groupAPI.ListGroup)
			g.POST("/group.bind", groupAPI.BindGroup)
			g.POST("/:groupID/member", groupAPI.AddMember)
			g.GET("/:groupID/projects", projectAPI.GetList)
		}
		group := v1.Group("/group")
		check := v1.Group("/check")
		{
			check.GET("/group", groupAPI.CheckGroup)
			check.GET("/member", groupAPI.CheckMember)
			check.GET("/developer", userAPI.CheckUser)
		}

		project := group.Group("/:groupID/project")
		{
			project.POST("", projectAPI.CreateProject)
			project.GET("/:projectID", projectAPI.GetProjectByID)
			project.PATCH("/:projectID/desc", projectAPI.UpdDescribe)
			project.DELETE("/:projectID", projectAPI.DelProject)
		}
		fnAPI := NewFunctionAPI(ctx, c, db, k8sClient, redisClient, esClient)
		f := project.Group("")
		{
			f.POST("/create", fnAPI.Create)
			f.DELETE("/:projectID/:functionID", fnAPI.Delete)
			f.GET("/get", fnAPI.Get)
			f.GET("/logger/:resourceRef", fnAPI.ListLog)
			f.GET("/list/:projectID", fnAPI.List)
			f.GET("/:projectID/:functionID", fnAPI.Get)
			f.POST("/regSwagger", fnAPI.RegSwagger)
			f.PATCH("/:projectID/:functionID/desc", fnAPI.UpdateDescribe)
		}
		svcApi := NewServing(db, c, k8sClient)
		{
			f.POST("/serve", svcApi.serve)
			f.DELETE("/offline", svcApi.offline)
		}

		er := newEvent(db)
		{
			v1.GET("/event/:id", er.GetEvent)
		}

		event.New(log, event.WithRouter(engine.Group("")),
			event.WithHandle(event.Function, fnAPI.fn.UpdateStatus, fnAPI.ps.Publish),
			// TODO:
			event.WithHandle(event.APIDoc, er.e.Save, fnAPI.fn.DeleteRegPipeline, fnAPI.ps.Publish),
			event.WithHandle(event.Serving, er.e.Save, fnAPI.fn.UpdateStatus, fnAPI.ps.Publish))

	}
	{
		probe := probe.New(log)
		engine.GET("liveness", func(c *gin.Context) {
			probe.LivenessProbe(c.Writer, c.Request)
		})

		engine.Any("readiness", func(c *gin.Context) {
			probe.ReadinessProbe(c.Writer, c.Request)
		})

	}

	return &Router{
		c:      c,
		engine: engine,
	}, nil
}

func newRouter(c *config.Config) (*gin.Engine, error) {
	if c.Model == "" || (c.Model != ReleaseMode && c.Model != DebugMode) {
		c.Model = ReleaseMode
	}
	gin.SetMode(c.Model)
	engine := gin.New()

	engine.Use(ginlogger.LoggerFunc(), ginlogger.LoggerFunc())

	return engine, nil
}

// Run run
func (r *Router) Run() {
	r.engine.Run(r.c.Port)
}

// Close close
func (r *Router) Close() {
}
