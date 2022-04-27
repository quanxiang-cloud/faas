package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/id"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/faas/internal/models"
	"github.com/quanxiang-cloud/faas/internal/models/es"
	"github.com/quanxiang-cloud/faas/internal/models/mysql"
	"github.com/quanxiang-cloud/faas/pkg/code"
	"github.com/quanxiang-cloud/faas/pkg/config"
	"github.com/quanxiang-cloud/faas/pkg/k8s"
	"gorm.io/gorm"
)

type Function interface {
	Create(c context.Context, r *CreateFunctionRequest) (*CreateFunctionResponse, error)
	UpdateStatus(c context.Context, r *UpdateFunctionRequest) (*UpdateFunctionResponse, error)
	Delete(c context.Context, r *DeleteFunctionRequest) (*DeleteFunctionResponse, error)
	Get(c context.Context, r *GetFunctionRequest) (*GetFunctionResponse, error)

	Build(c context.Context, r *BuildFunctionRequest) (*BuildFunctionResponse, error)
	DelFunction(c context.Context, r *DelBuildFunctionRequest) (*DelBuildFunctionResponse, error)
	ListLog(c context.Context, r *ListlogRequest) (*ListLogResponse, error)
	List(c context.Context, r *ListRequest) (*ListResponse, error)

	RegSwagger(c context.Context, r *RegSwaggerReq) (*RegSwaggerResp, error)
	UpdateDoc(c context.Context, req *UpdateDocReq) (*UpdateDocResp, error)
}

type function struct {
	db           *gorm.DB
	functionRepo models.FunctionRepo
	gitRepo      models.GitRepo
	dockerRepo   models.DockerRepo
	k8sc         k8s.Client
	conf         config.Config
	buildLogRepo models.BuilderLogRepo
	groupRepo    models.GroupRepo
	projectRepo  models.ProjectRepo
}

func NewFunction(c context.Context, db *gorm.DB, conf config.Config, kc k8s.Client, esClient *elastic.Client) Function {
	return &function{
		db:           db,
		functionRepo: mysql.NewFunctionRepo(),
		gitRepo:      mysql.NewGitRepo(),
		dockerRepo:   mysql.NewDockerRepo(),
		conf:         conf,
		k8sc:         kc,
		buildLogRepo: es.NewBuildLogRepo(esClient),
		projectRepo:  mysql.NewProjectRepo(),
		groupRepo:    mysql.NewGroupRepo(),
	}
}

type CreateFunctionRequest struct {
	GroupID   string            `json:"groupID"  form:"groupID" uri:"groupID"`
	ProjectID string            `json:"projectID" form:"projectID" uri:"projectID"`
	Version   string            `json:"version" form:"version" uri:"version"`
	Describe  string            `json:"describe" form:"describe" uri:"describe"`
	Env       map[string]string `json:"env"`
}
type CreateFunctionResponse struct {
	ID string `json:"id"`
}

func (g *function) Create(c context.Context, r *CreateFunctionRequest) (*CreateFunctionResponse, error) {
	data := &models.Function{}
	data.ID = id.ShortID(0)
	data.GroupID = r.GroupID
	data.Describe = r.Describe
	data.ProjectID = r.ProjectID
	data.Version = r.Version
	if r.Env != nil {
		marshal, _ := json.Marshal(r.Env)
		data.Env = string(marshal)
	}
	group, err := g.groupRepo.Get(g.db, r.GroupID)
	if err != nil || group == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	project, err := g.projectRepo.Get(g.db, r.ProjectID)
	if err != nil || project == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	data.Name = strings.ToLower(group.GroupName) + "-" + project.ProjectName + "-" + data.Version
	unix := time2.NowUnix()
	data.CreatedAt = unix
	data.UpdatedAt = unix
	one := g.functionRepo.GetByName(c, g.db, data.Name)
	if one != nil {
		return nil, error2.New(code.ErrFunctionExist)
	}
	return &CreateFunctionResponse{
		ID: data.ID,
	}, g.functionRepo.Insert(c, g.db, data)
}

type UpdateFunctionRequest struct {
	//Labels map[string]string `json:"labels"`
	State       string `json:"state"`
	ResourceRef string `json:"resourceRef"`
	Topic       string `json:"topic"`
	Name        string `json:"name"`
}
type UpdateFunctionResponse struct {
	ID     string `json:"-"`
	Status int    `json:"-"`
	Topic  string `json:"-"`
}

type functionStatus int

const (
	StatusNull functionStatus = iota
	StatusBuilding
	StatusFailed
	StatusOK
	StatusCancelled

	StatusOnline
	StatusOffline
)

var result = map[string]int{
	"Building":  int(StatusBuilding),
	"Succeeded": int(StatusOK),
	"Failed":    int(StatusFailed),
	"Cancelled": int(StatusCancelled),
}

func (g *function) UpdateStatus(c context.Context, r *UpdateFunctionRequest) (*UpdateFunctionResponse, error) {

	data := g.functionRepo.GetByName(c, g.db, r.Name)
	if data == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	if v, ok := result[r.State]; ok && v != 0 {
		data.Status = v
	}
	unix := time2.NowUnix()
	data.UpdatedAt = unix
	data.ResourceRef = r.ResourceRef
	return &UpdateFunctionResponse{
		ID:     data.ID,
		Status: result[r.State],
		Topic:  r.Topic,
	}, g.functionRepo.Update(c, g.db, data)
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
	Describe  string `json:"describe"`
	UpdatedAt int64  `json:"updatedAt"`
}

func (g *function) Get(c context.Context, r *GetFunctionRequest) (*GetFunctionResponse, error) {
	data := g.functionRepo.Get(c, g.db, r.ID)
	if data == nil {
		return nil, nil
	}
	group, err := g.groupRepo.Get(g.db, data.GroupID)
	if err != nil || group == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	project, err := g.projectRepo.Get(g.db, data.ProjectID)
	if err != nil || project == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	res := &GetFunctionResponse{
		ID:        data.ID,
		GroupName: group.GroupName,
		Project:   project.ProjectName,
		Version:   data.Version,
		Describe:  data.Describe,
		UpdatedAt: data.UpdatedAt,
	}
	return res, nil

}

type BuildFunctionRequest struct {
	ID string `json:"id"`
}
type BuildFunctionResponse struct {
	PodName string `json:"-"`
}

func (g *function) Build(c context.Context, r *BuildFunctionRequest) (*BuildFunctionResponse, error) {
	fnData := g.functionRepo.Get(c, g.db, r.ID)
	if fnData == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	group, err := g.groupRepo.Get(g.db, fnData.GroupID)
	if err != nil || group == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	project, err := g.projectRepo.Get(g.db, fnData.ProjectID)
	if err != nil || project == nil {
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
		Project:   project.ProjectName,
		GroupName: group.GroupName,
		Git: &k8s.Git{
			Name: gitData.Name,
			Host: gitData.KnownHosts,
		},
		Docker: &k8s.Docker{
			NameSpace: dockerData.NameSpace,
			Name:      dockerData.Name,
			Host:      dockerData.Host,
		},
		Builder: k8s.GetBuilder(project.Language),
		ENV:     env,
	})
}

type DelBuildFunctionRequest struct {
	ID     string `json:"id"`
	Status int    `json:"-"`
}
type DelBuildFunctionResponse struct {
}

func (g *function) DelFunction(c context.Context, r *DelBuildFunctionRequest) (*DelBuildFunctionResponse, error) {

	if r.Status == int(StatusOK) || r.Status == int(StatusFailed) {
		fnData := g.functionRepo.Get(c, g.db, r.ID)
		group, err := g.groupRepo.Get(g.db, fnData.GroupID)
		if err != nil || group == nil {
			return nil, error2.New(code.ErrDataNotExist)
		}
		project, err := g.projectRepo.Get(g.db, fnData.ProjectID)
		if err != nil || project == nil {
			return nil, error2.New(code.ErrDataNotExist)
		}
		if fnData != nil {
			return nil, g.k8sc.DelFunction(c, &k8s.DelFunction{
				Name: strings.ToLower(group.GroupName) + "-" + project.ProjectName + "-" + fnData.Version,
			})
		}
		return nil, nil
	}
	return nil, nil

}

type ListlogRequest struct {
	ResourceRef string `json:"resourceRef" form:"resourceRef" uri:"resourceRef"`
	Step        string `json:"step" form:"step" uri:"step"`
	Index       int    `json:"index" form:"index"`
	Timestamp   int64  `json:"timestamp" form:"timestamp"`
}
type ListLogResponse struct {
	Logs  []*models.LogVO `json:"logs"`
	Count int64           `json:"count"`
}

func (g *function) ListLog(c context.Context, r *ListlogRequest) (*ListLogResponse, error) {
	fn := g.functionRepo.GetByResourceRef(c, g.db, r.ResourceRef)
	if fn == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}

	t := time.Unix(r.Timestamp, 0)
	fullLogs, count, err := g.buildLogRepo.Search(c, r.ResourceRef, r.Step, t, r.Index, 100)
	if err != nil {
		return nil, err
	}

	logs := make([]*models.LogVO, 0, count)
	for _, e := range fullLogs {
		logs = append(logs, &models.LogVO{
			Run:       e.Labels.PipelineTask,
			Step:      e.ContainerName,
			Log:       e.Log,
			Timestamp: e.Time.Unix(),
			PodName:   e.PodName,
		})
	}

	return &ListLogResponse{
		Logs:  logs,
		Count: count,
	}, nil
}

type ListRequest struct {
	Page      int    `json:"page" form:"page" uri:"page"`
	Limit     int    `json:"limit" form:"limit" uri:"limit"`
	GroupID   string `json:"groupID" form:"groupID"  uri:"groupID"`
	ProjectID string `json:"projectID" form:"projectID" uri:"projectID"`
}
type ListResponse struct {
	Data  []RespFunction `json:"data"`
	Count int64          `json:"count"`
}

type RespFunction struct {
	ID          string `json:"id"`
	GroupID     string `json:"groupID"`
	ProjectID   string `json:"projectID"`
	Version     string `json:"version"`
	Describe    string `json:"describe"`
	Status      int    `json:"status"`
	Env         string `json:"env"`
	ResourceRef string `json:"resourceRef"`
	Name        string ` json:"name"`

	CreatedAt int64  `json:"createdAt,omitempty" `
	CreatedBy string `json:"createdBy,omitempty"` //创建者
}

func (g *function) List(c context.Context, r *ListRequest) (*ListResponse, error) {
	fns, count := g.functionRepo.Search(c, g.db, r.ProjectID, r.GroupID, r.Page, r.Limit)
	if len(fns) == 0 {
		return nil, error2.New(code.ErrDataNotExist)
	}

	res := make([]RespFunction, 0, len(fns))
	for k := range fns {
		res = append(res, RespFunction{
			ID:          fns[k].ID,
			GroupID:     fns[k].GroupID,
			ProjectID:   fns[k].ProjectID,
			Version:     fns[k].Version,
			Describe:    fns[k].Describe,
			Status:      fns[k].Status,
			Env:         fns[k].Env,
			ResourceRef: fns[k].ResourceRef,
			Name:        fns[k].Name,
		})

	}

	return &ListResponse{
		Data:  res,
		Count: count,
	}, nil
}

type RegSwaggerReq struct {
	ID    string `json:"id"`
	AppID string `json:"appID"`
}

type RegSwaggerResp struct {
}

func (g *function) RegSwagger(c context.Context, r *RegSwaggerReq) (*RegSwaggerResp, error) {
	fn := g.functionRepo.Get(c, g.db, r.ID)
	if fn.Status != int(StatusOK) {
		return nil, error2.New(code.ErrDataIllegal)
	}

	project, err := g.projectRepo.Get(g.db, fn.ProjectID)
	if err != nil {
		return nil, err
	}

	group, err := g.groupRepo.Get(g.db, fn.GroupID)
	if err != nil {
		return nil, err
	}

	git := g.gitRepo.Get(c, g.db)

	// TODO: credentials
	err = g.k8sc.RegistAPI(c, &k8s.Function{
		Version:   fn.Version,
		Project:   project.ProjectName,
		GroupName: group.GroupName,
		Git: &k8s.Git{
			Name: git.Name,
			Host: git.KnownHosts,
		},
	}, r.AppID)
	return &RegSwaggerResp{}, err
}

type UpdateDocReq struct {
	Name   string `json:"name"`
	Status string `json:"state"`
	Topic  string `json:"topic"`
}

type UpdateDocResp struct {
	ID    string `json:"id"`
	Topic string `json:"topic"`
}

func (g *function) UpdateDoc(c context.Context, req *UpdateDocReq) (*UpdateDocResp, error) {
	name, err := k8s.ReverseName(req.Name)
	if err != nil {
		return nil, err
	}
	fn := g.functionRepo.GetByName(c, g.db, name)
	if fn == nil {
		return nil, fmt.Errorf("can not find function(%s)", name)
	}

	// TODO: ???
	fn.Doc = 0
	if req.Status == "Succeeded" {
		fn.Doc = 1
	}

	if err := g.k8sc.DeleteReigstRun(c, req.Name); err != nil {
		return nil, err
	}

	if err := g.functionRepo.Update(c, g.db, fn); err != nil {
		return nil, err
	}

	return &UpdateDocResp{
		ID:    fn.ID,
		Topic: req.Topic,
	}, nil
}
