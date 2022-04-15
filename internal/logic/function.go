package logic

import (
	"context"
	"encoding/json"
	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/faas/internal/models"
	"github.com/quanxiang-cloud/faas/internal/models/mysql"
	"github.com/quanxiang-cloud/faas/pkg/code"
	"github.com/quanxiang-cloud/faas/pkg/config"
	"github.com/quanxiang-cloud/faas/pkg/k8s"
	"gorm.io/gorm"
	"strings"
)

type Function interface {
	Create(c context.Context, r *CreateFunctionRequest) (*CreateFunctionResponse, error)
	UpdateStatus(c context.Context, r *UpdateFunctionRequest) (*UpdateFunctionResponse, error)
	Delete(c context.Context, r *DeleteFunctionRequest) (*DeleteFunctionResponse, error)
	Get(c context.Context, r *GetFunctionRequest) (*GetFunctionResponse, error)

	Build(c context.Context, r *BuildFunctionRequest) (*BuildFunctionResponse, error)
	DelFunction(c context.Context, r *DelBuildFunctionRequest) (*DelBuildFunctionResponse, error)
}

type function struct {
	db           *gorm.DB
	functionRepo models.FunctionRepo
	gitRepo      models.GitRepo
	dockerRepo   models.DockerRepo
	k8sc         k8s.Client
	conf         config.Config
}

func NewFunction(c context.Context, db *gorm.DB, conf config.Config, kc k8s.Client) Function {
	return &function{
		db:           db,
		functionRepo: mysql.NewFunctionRepo(),
		gitRepo:      mysql.NewGitRepo(),
		dockerRepo:   mysql.NewDockerRepo(),
		conf:         conf,
	}
}

type CreateFunctionRequest struct {
	GroupName string            `json:"groupName"`
	Project   string            `json:"project"`
	Version   string            `json:"version"`
	Language  string            `json:"language"`
	Env       map[string]string `json:"env"`
}
type CreateFunctionResponse struct {
	ID string `json:"id"`
}

func (g *function) Create(c context.Context, r *CreateFunctionRequest) (*CreateFunctionResponse, error) {
	data := &models.Function{}
	data.ID = id.ShortID(0)
	data.GroupName = r.GroupName
	data.Language = r.Language
	data.Project = r.Project
	data.Version = r.Version
	if r.Env != nil {
		marshal, _ := json.Marshal(r.Env)
		data.Env = string(marshal)
	}
	unix := time.NowUnix()
	data.CreatedAt = unix
	data.UpdatedAt = unix
	return &CreateFunctionResponse{}, g.functionRepo.Insert(c, g.db, data)
}

type UpdateFunctionRequest struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
}
type UpdateFunctionResponse struct {
	ID string `json:"-"`
}

type functionStatus int

const (
	StatusNull functionStatus = iota
	StatusBuilding
	StatusFailed
	StatusOK
)

func (g *function) UpdateStatus(c context.Context, r *UpdateFunctionRequest) (*UpdateFunctionResponse, error) {
	data := g.functionRepo.Get(c, g.db, r.ID)
	if data == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	data.Status = r.Status
	unix := time.NowUnix()
	data.UpdatedAt = unix
	return &UpdateFunctionResponse{}, g.functionRepo.Update(c, g.db, data)
}

type DeleteFunctionRequest struct {
	ID string `json:"id"`
}
type DeleteFunctionResponse struct {
}

func (g *function) Delete(c context.Context, r *DeleteFunctionRequest) (*DeleteFunctionResponse, error) {
	data := g.functionRepo.Get(c, g.db, r.ID)
	if data == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	return &DeleteFunctionResponse{}, g.functionRepo.Delete(c, g.db, r.ID)
}

type GetFunctionRequest struct {
	ID string `json:"id"`
}
type GetFunctionResponse struct {
	ID        string `json:"id"`
	GroupName string `json:"groupName"`
	Project   string `json:"project"`
	Version   string `json:"version"`
	Language  string `json:"language"`
	UpdatedAt int64  `json:"updatedAt"`
}

func (g *function) Get(c context.Context, r *GetFunctionRequest) (*GetFunctionResponse, error) {
	data := g.functionRepo.Get(c, g.db, r.ID)
	if data == nil {
		return nil, nil
	}
	res := &GetFunctionResponse{
		ID:        data.ID,
		GroupName: data.GroupName,
		Project:   data.Project,
		Version:   data.Version,
		Language:  data.Language,
		UpdatedAt: data.UpdatedAt,
	}
	return res, nil

}

type BuildFunctionRequest struct {
	ID string `json:"id"`
}
type BuildFunctionResponse struct {
}

func (g *function) Build(c context.Context, r *BuildFunctionRequest) (*BuildFunctionResponse, error) {
	fnData := g.functionRepo.Get(c, g.db, r.ID)
	if fnData == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	gitData := g.gitRepo.Get(c, g.db)
	if gitData == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	dockerData := g.dockerRepo.Get(c, g.db)
	if dockerData == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	env := make(map[string]string)
	if fnData.Env != "" {
		json.Unmarshal([]byte(fnData.Env), &env)
	}
	return &BuildFunctionResponse{}, g.k8sc.Build(c, &k8s.Function{
		Version:   fnData.Version,
		Project:   fnData.Project,
		GroupName: fnData.GroupName,
		Git: &k8s.Git{
			Name: gitData.Name,
			Host: gitData.Host,
		},
		Docker: &k8s.Docker{
			NameSpace: g.conf.Docker.NameSpace,
			Name:      dockerData.Name,
			Host:      dockerData.Host,
		},
		Builder: k8s.GetBuilder(fnData.Language),
		ENV:     env,
	})
}

type DelBuildFunctionRequest struct {
	ID string `json:"id"`
}
type DelBuildFunctionResponse struct {
}

func (g *function) DelFunction(c context.Context, r *DelBuildFunctionRequest) (*DelBuildFunctionResponse, error) {
	fnData := g.functionRepo.Get(c, g.db, r.ID)
	if fnData.Status == int(StatusOK) || fnData.Status == int(StatusFailed) {
		return nil, g.k8sc.DelFunction(c, &k8s.DelFunction{
			Name: strings.ToLower(fnData.GroupName) + "-" + fnData.Project + "-" + fnData.Version,
		})
	}
	return nil, nil

}