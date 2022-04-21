package restful

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	"github.com/quanxiang-cloud/faas/pkg/k8s"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/quanxiang-cloud/cabin/logger"
	ginlogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
	"github.com/quanxiang-cloud/faas/pkg/config"
	"github.com/quanxiang-cloud/faas/pkg/probe"
	"github.com/quanxiang-cloud/faas/pkg/util"
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
func NewRouter(ctx context.Context, c *config.Config, log logger.AdaptedLogger, db *gorm.DB, kc k8s.Client, rc redis.UniversalClient, es *elastic.Client) (*Router, error) {
	engine, err := newRouter(c)
	if err != nil {
		return nil, err
	}

	v1 := engine.Group("/api/v1/faas")

	{
		gitAPI := NewGitAPI(ctx, c, db, kc)
		g := v1.Group("/git")
		{
			g.POST("", gitAPI.Create)
			g.PUT("/:id", gitAPI.Update)
			g.DELETE("/:id", gitAPI.Delete)
			g.GET("", gitAPI.Get)
		}
		dockerAPI := NewDockerAPI(ctx, c, db, kc)
		d := v1.Group("/docker")
		{
			d.POST("", dockerAPI.Create)
			d.PUT("/:id", dockerAPI.Update)
			d.DELETE("/:id", dockerAPI.Delete)
			d.GET("", dockerAPI.Get)
		}
		cm := NewCompoundAPI(ctx, rc)
		cmGroup := v1.Group("/cm")
		{
			cmGroup.POST("/subscribe", cm.Subscribe)
		}
	}
	{
		fnAPI := NewFunctionAPI(ctx, c, db, kc, rc, es)
		f := v1.Group("/fn")
		{
			f.POST("/create", fnAPI.Create)
			f.POST("/update/status", fnAPI.UpdateStatus)
			f.DELETE("/del", fnAPI.Delete)
			f.GET("/get", fnAPI.Get)
			f.GET("/:resourceRef/logger", fnAPI.ListLog)
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
		group := v1.Group("/group")
		{
			group.POST("", groupAPI.Create)
			group.POST("/:groupID/member", groupAPI.BindingGroup)
			group.GET("/:groupID/projects", projectAPI.GetList)
		}
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
	}
	{
		probe := probe.New(util.LoggerFromContext(ctx))
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
