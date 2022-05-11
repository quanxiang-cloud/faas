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
	"github.com/quanxiang-cloud/faas/pkg/basic/define/code"
	"github.com/quanxiang-cloud/faas/pkg/basic/define/consts"
	"github.com/quanxiang-cloud/faas/pkg/basic/event"
	"github.com/quanxiang-cloud/faas/pkg/basic/k8s"
	"github.com/quanxiang-cloud/faas/pkg/basic/strutil"
	"github.com/quanxiang-cloud/faas/pkg/config"
	"gorm.io/gorm"
)

type Function interface {
	Create(c context.Context, r *CreateFunctionRequest) (*CreateFunctionResponse, error)
	Delete(c context.Context, r *DeleteFunctionRequest) (*DeleteFunctionResponse, error)
	Get(c context.Context, r *GetFunctionRequest) (*GetFunctionResponse, error)
	UpdateDescribe(c context.Context, r *UpdateFuncDescribeReq) (*UpdateFuncDescribeResp, error)

	Build(c context.Context, r *BuildFunctionRequest) (*BuildFunctionResponse, error)
	DelFunction(c context.Context, r *DelBuildFunctionRequest) (*DelBuildFunctionResponse, error)
	ListLog(c context.Context, r *ListLogRequest) (*ListLogResponse, error)
	List(c context.Context, r *ListRequest) (*ListResponse, error)

	RegSwagger(c context.Context, r *RegSwaggerReq) (*RegSwaggerResp, error)
	DeleteRegPipeline(msg *event.MsgBus) error

	UpdateStatus(*event.MsgBus) error
	UpdateDocStatus(*event.MsgBus) error
	UpdateServingStatus(bus *event.MsgBus) error
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
	CreatedBy string            `json:"-"`
}
type CreateFunctionResponse struct {
	ID string `json:"id"`
}

// TODO: default env
func (g *function) Create(c context.Context, r *CreateFunctionRequest) (*CreateFunctionResponse, error) {
	data := &models.Function{}
	data.ID = id.ShortID(0)
	data.GroupID = r.GroupID
	data.Describe = r.Describe
	data.ProjectID = r.ProjectID
	data.Version = r.Version
	data.CreatedBy = r.CreatedBy
	data.DocStatus = consts.DocNotExists
	if r.Env != nil {
		marshal, _ := json.Marshal(appendDefaultEnv(r.Env))
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
	data.BuiltAt = unix
	one := g.functionRepo.GetByName(c, g.db, data.Name)
	if one != nil {
		return nil, error2.New(code.ErrFunctionExist)
	}
	return &CreateFunctionResponse{
		ID: data.ID,
	}, g.functionRepo.Insert(c, g.db, data)
}

func appendDefaultEnv(env map[string]string) map[string]string {
	env[consts.DefaultEnvFuncClearSource] = consts.DefaultEnvFuncClearSourceVal
	env[consts.DefaultEnvFuncName] = consts.DefaultEnvFuncNameVal
	if _, ok := env[consts.DefaultEnvGoProxy]; !ok {
		env[consts.DefaultEnvGoProxy] = consts.DefaultEnvGoProxyVal
	}
	return env
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
	StatusServing
	StatusOnlineFailed
)

var result = map[string]int{
	"Building":  int(StatusBuilding),
	"Succeeded": int(StatusOK),
	"Failed":    int(StatusFailed),
	"Cancelled": int(StatusCancelled),
	"True":      int(StatusOnline),
	"Unknown":   int(StatusOnline),
	"False":     int(StatusOnlineFailed),
}

func (g *function) UpdateStatus(bus *event.MsgBus) error {
	data := g.functionRepo.GetByName(bus.CTX, g.db, bus.Msg.Fn.Name)
	if data == nil {
		return error2.New(code.ErrDataNotExist)
	}
	if v, ok := result[bus.Msg.Fn.State]; ok && v != 0 {
		data.Status = v
	}
	unix := time2.NowUnix()
	data.UpdatedAt = unix
	data.BuiltAt = unix
	data.ResourceRef = bus.Msg.Fn.ResourceRef
	if err := g.functionRepo.Update(bus.CTX, g.db, data); err != nil {
		return err
	}

	_, err := g.DelFunction(bus.CTX, &DelBuildFunctionRequest{
		ID:     data.ID,
		Status: data.Status,
	})

	bus.Data = data.ID
	return err
}

func (g *function) UpdateServingStatus(bus *event.MsgBus) error {
	fnName := strutil.Reverse(bus.Msg.Svc.Name)
	data := g.functionRepo.GetByName(bus.CTX, g.db, fnName)
	if data == nil {
		return error2.New(code.ErrDataNotExist)
	}
	if v, ok := result[bus.Msg.Svc.State]; ok && v != 0 {
		data.Status = v
	}
	unix := time2.NowUnix()
	data.UpdatedAt = unix
	if err := g.functionRepo.Update(bus.CTX, g.db, data); err != nil {
		return err
	}
	bus.Data = data.ID
	return nil
}

type DeleteFunctionRequest struct {
	FunctionID string `json:"functionID" form:"functionID" uri:"functionID"`
}
type DeleteFunctionResponse struct {
}

func (g *function) Delete(c context.Context, r *DeleteFunctionRequest) (*DeleteFunctionResponse, error) {
	data := g.functionRepo.Get(c, g.db, r.FunctionID)
	if data == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	return &DeleteFunctionResponse{}, g.functionRepo.Delete(c, g.db, r.FunctionID)
}

type GetFunctionRequest struct {
	ID string `json:"functionID" form:"functionID" uri:"functionID"`
}
type GetFunctionResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	GroupName   string `json:"groupName"`
	Project     string `json:"project"`
	Version     string `json:"version"`
	Describe    string `json:"describe"`
	UpdatedAt   int64  `json:"updatedAt"`
	CreatedBy   string `json:"createdBy"`
	CreatedAt   int64  `json:"createdAt"`
	Status      int    `json:"status"`
	DocStatus   int    `json:"docStatus"`
	BuiltAt     int64  `json:"builtAt"`
	ResourceRef string `json:"resourceRef"`
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
		ID:          data.ID,
		Name:        data.Name,
		GroupName:   group.GroupName,
		Project:     project.ProjectName,
		Version:     data.Version,
		Describe:    data.Describe,
		UpdatedAt:   data.UpdatedAt,
		CreatedBy:   data.CreatedBy,
		CreatedAt:   data.CreatedAt,
		Status:      data.Status,
		DocStatus:   data.DocStatus,
		ResourceRef: data.ResourceRef,
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

	builder, err := g.k8sc.GetBuilder(project.Language, project.Version)
	if err != nil {
		return nil, err
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
		Builder: builder,
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

type ListLogRequest struct {
	ResourceRef string `json:"resourceRef" form:"resourceRef" uri:"resourceRef"`
	Step        string `json:"step" form:"step" uri:"step"`
	Index       int    `json:"index" form:"index"`
	Timestamp   int64  `json:"timestamp" form:"timestamp"`
}
type ListLogResponse struct {
	Logs  []*models.LogVO `json:"logs"`
	Count int64           `json:"count"`
}

func (g *function) ListLog(c context.Context, r *ListLogRequest) (*ListLogResponse, error) {
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
	DocStatus   int    `json:"docStatus"`
	BuiltAt     int64  `json:"builtAt"`
	Env         string `json:"env"`
	ResourceRef string `json:"resourceRef"`
	Name        string ` json:"name"`

	CreatedAt int64  `json:"createdAt,omitempty" `
	CreatedBy string `json:"createdBy,omitempty"` //创建者
}

func (g *function) List(c context.Context, r *ListRequest) (*ListResponse, error) {
	fns, count := g.functionRepo.Search(c, g.db, r.ProjectID, r.GroupID, r.Page, r.Limit)
	if len(fns) == 0 {
		return &ListResponse{
			Data:  nil,
			Count: 0,
		}, nil
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
			DocStatus:   fns[k].DocStatus,
			Env:         fns[k].Env,
			ResourceRef: fns[k].ResourceRef,
			Name:        fns[k].Name,
			CreatedBy:   fns[k].CreatedBy,
			CreatedAt:   fns[k].CreatedAt,
		})

	}

	return &ListResponse{
		Data:  res,
		Count: count,
	}, nil
}

type RegSwaggerReq struct {
	ID      string `json:"buildID"`
	GroupID string `json:"-"`
}

type RegSwaggerResp struct {
}

func (g *function) RegSwagger(c context.Context, r *RegSwaggerReq) (*RegSwaggerResp, error) {
	fn := g.functionRepo.Get(c, g.db, r.ID)
	group, err := g.groupRepo.Get(g.db, r.GroupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, fmt.Errorf("group not exist")
	}

	project, err := g.projectRepo.Get(g.db, fn.ProjectID)
	if err != nil {
		return nil, err
	}

	git := g.gitRepo.Get(c, g.db)

	err = g.k8sc.RegistAPI(c, &k8s.Function{
		Version:   fn.Version,
		Project:   project.ProjectName,
		GroupName: group.GroupName,
		Git: &k8s.Git{
			Name: git.Name,
			Host: git.KnownHosts,
		},
	}, group.AppID)

	fn.DocStatus = consts.DocTaskExists
	g.functionRepo.Update(c, g.db, fn)

	return &RegSwaggerResp{}, err
}

func (g *function) DeleteRegPipeline(bus *event.MsgBus) error {
	var err error
	if bus.Msg.Pr.State == "True" || bus.Msg.Pr.State == "False" {
		err = g.k8sc.DeleteReigstRun(bus.CTX, bus.Msg.Pr.Name)
	}
	return err
}

type UpdateFuncDescribeReq struct {
	ID       string `json:"functionID"`
	Describe string `json:"describe"`
}

type UpdateFuncDescribeResp struct {
}

func (g *function) UpdateDescribe(c context.Context, r *UpdateFuncDescribeReq) (*UpdateFuncDescribeResp, error) {
	fnData := g.functionRepo.Get(c, g.db, r.ID)
	if fnData == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	fnData.Describe = r.Describe
	fnData.UpdatedAt = time2.NowUnix()
	err := g.functionRepo.UpdateDescribe(c, g.db, fnData)
	if err != nil {
		return nil, err
	}
	return &UpdateFuncDescribeResp{}, nil
}

func (g *function) UpdateDocStatus(bus *event.MsgBus) error {
	fnName := strutil.Reverse(bus.Msg.Pr.Name)
	data := g.functionRepo.GetByName(bus.CTX, g.db, fnName)
	if data == nil {
		return error2.New(code.ErrDataNotExist)
	}

	if v, ok := consts.DocStatusMapping[bus.Msg.Pr.State]; ok {
		data.DocStatus = v
	}
	unix := time2.NowUnix()
	data.UpdatedAt = unix
	if err := g.functionRepo.Update(bus.CTX, g.db, data); err != nil {
		return err
	}

	bus.Data = data.ID
	return nil
}
