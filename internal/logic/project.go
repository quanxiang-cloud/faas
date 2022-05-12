package logic

import (
	"context"
	"fmt"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/faas/internal/models"
	"github.com/quanxiang-cloud/faas/internal/models/mysql"
	"github.com/quanxiang-cloud/faas/pkg/basic/define/code"
	git2 "github.com/quanxiang-cloud/faas/pkg/basic/git"
	"github.com/quanxiang-cloud/faas/pkg/config"
	"gorm.io/gorm"
)

type ProjectService interface {
	CreateProject(ctx context.Context, req *CreateProjectReq) (*CreateProjectResp, error)
	GetProjectByID(ctx context.Context, req *GetProjectByIDReq) (*GetProjectByIDResp, error)
	GetList(ctx context.Context, req *GetListReq) (*GetListResp, error)
	DelProject(ctx context.Context, req *DelProjectReq) (*DelProjectResp, error)
	UpdateDescribe(ctx context.Context, req *UpdateDescribeReq) (*UpdateDescribeResp, error)
	ListGITProjects(ctx context.Context, req *ListGITProjectsReq) (*ListGITProjectsResp, error)
}

type project struct {
	db            *gorm.DB
	projectRepo   models.ProjectRepo
	gitRepo       models.GitRepo
	groupRepo     models.GroupRepo
	userGroupRepo models.UserGroupRepo
}

func NewProjectService(ctx context.Context, db *gorm.DB, conf *config.Config) ProjectService {
	return &project{
		db:            db,
		projectRepo:   mysql.NewProjectRepo(),
		gitRepo:       mysql.NewGitRepo(),
		groupRepo:     mysql.NewGroupRepo(),
		userGroupRepo: mysql.NewUserGroupRepo(),
	}
}

type CreateProjectReq struct {
	GroupID  string `json:"-"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Alias    string `json:"alias"`
	Language string `json:"language"`
	Version  string `json:"version"`
	Describe string `json:"description"`
	UserID   string `json:"-"`
}

type CreateProjectResp struct {
	ID        string `json:"id"`
	RepoUrl   string `json:"repoUrl"`
	CreatedAt int64  `json:"createdAt"`
	CreatedBy string `json:"creator"`
}

func (p *project) CreateProject(ctx context.Context, req *CreateProjectReq) (*CreateProjectResp, error) {
	tx := p.db.Begin()
	gitHost := p.gitRepo.Get(ctx, tx)
	if gitHost == nil {
		tx.Rollback()
		return nil, error2.New(code.ErrDataNotExist)
	}
	client, err := git2.GetClient(git2.Gitlab, gitHost.Token, gitHost.Host)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	group, err := p.groupRepo.Get(p.db, req.GroupID)
	if err != nil {
		return nil, err
	}

	pmodel := &models.Project{
		ID:          id.StringUUID(),
		ProjectID:   req.ID,
		ProjectName: req.Name,
		RepoUrl:     genRepoURl(gitHost.KnownHosts, group.GroupName, req.Name),
		Alias:       req.Alias,
		Describe:    req.Describe,
		Language:    req.Language,
		Version:     req.Version,
		Status:      models.ProjectSuccessStatus,
		GroupID:     req.GroupID,
		UserID:      req.UserID,
		CreatedBy:   req.UserID,
		UpdatedBy:   req.UserID,
		CreatedAt:   time.NowUnix(),
		UpdatedAt:   time.NowUnix(),
	}

	if req.ID == 0 {
		project, err := client.CreateProject(ctx, req.Name, group.GroupID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		pmodel.ProjectID = project.ID
	}

	err = p.projectRepo.Insert(tx, pmodel)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &CreateProjectResp{
		ID:        pmodel.ID,
		RepoUrl:   pmodel.RepoUrl,
		CreatedBy: pmodel.CreatedBy,
		CreatedAt: pmodel.CreatedAt,
	}, nil
}

func genRepoURl(host, group, project string) string {
	return fmt.Sprintf("%s%s/%s.git", host, group, project)
}

type GetProjectByIDReq struct {
	ProjectID string `json:"-"`
}

type GetProjectByIDResp struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	Status      string `json:"state"`
	Language    string `json:"language"`
	Version     string `json:"version"`
	RepoUrl     string `json:"repoUrl"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
	CreatedBy   string `json:"creator"`
}

func (p *project) GetProjectByID(ctx context.Context, req *GetProjectByIDReq) (*GetProjectByIDResp, error) {
	project, err := p.projectRepo.Get(p.db, req.ProjectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	return &GetProjectByIDResp{
		ID:        project.ID,
		Name:      project.ProjectName,
		Alias:     project.Alias,
		Status:    models.ProjectStatus[project.Status],
		Language:  project.Language,
		Version:   project.Version,
		RepoUrl:   project.RepoUrl,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
		CreatedBy: project.CreatedBy,
	}, nil
}

type Project struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	Status      string `json:"state"`
	Language    string `json:"language"`
	Version     string `json:"version"`
	RepoUrl     string `json:"repoUrl"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
	CreatedBy   string `json:"creator"`
}

type GetListReq struct {
	GroupID string `json:"-"`
	Alias   string `json:"alias"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}

type GetListResp struct {
	Projects []*Project `json:"projects"`
	Count    int64      `json:"count"`
}

func (p *project) GetList(ctx context.Context, req *GetListReq) (*GetListResp, error) {
	projects, count, err := p.projectRepo.GetByGroup(p.db, req.Alias, req.GroupID, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	resp := &GetListResp{
		Projects: make([]*Project, 0, len(projects)),
		Count:    count,
	}
	for _, project := range projects {
		resp.Projects = append(resp.Projects, &Project{
			ID:          project.ID,
			Name:        project.ProjectName,
			Alias:       project.Alias,
			Status:      models.ProjectStatus[project.Status],
			Description: project.Describe,
			Language:    project.Language,
			Version:     project.Version,
			RepoUrl:     project.RepoUrl,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
			CreatedBy:   project.CreatedBy,
		})
	}
	return resp, nil
}

type DelProjectReq struct {
	ProjectID string
}

type DelProjectResp struct {
}

func (p *project) DelProject(ctx context.Context, req *DelProjectReq) (*DelProjectResp, error) {
	err := p.projectRepo.Del(p.db, req.ProjectID)
	if err != nil {
		return nil, err
	}
	return &DelProjectResp{}, nil
}

type UpdateDescribeReq struct {
	ProjectID string `json:"-"`
	Describe  string `json:"description"`
}

type UpdateDescribeResp struct {
}

func (p *project) UpdateDescribe(ctx context.Context, req *UpdateDescribeReq) (*UpdateDescribeResp, error) {
	err := p.projectRepo.UpdDescribe(p.db, req.ProjectID, req.Describe)
	if err != nil {
		return nil, err
	}
	return &UpdateDescribeResp{}, nil
}

type ListGITProjectsReq struct {
	GroupID string `json:"groupID"`
}

type ListGITProjectsResp struct {
	Projects []GroupProjects `json:"projects"`
}

type GroupProjects struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (p *project) ListGITProjects(ctx context.Context, req *ListGITProjectsReq) (*ListGITProjectsResp, error) {
	gitHost := p.gitRepo.Get(ctx, p.db)
	if gitHost == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}

	group, err := p.groupRepo.Get(p.db, req.GroupID)
	if err != nil {
		return nil, err
	}

	client, err := git2.GetClient(git2.Gitlab, gitHost.Token, gitHost.Host)
	if err != nil {
		return nil, err
	}

	projects, err := client.GetGroupProjects(ctx, group.GroupID)
	if err != nil {
		return nil, err
	}

	ret := make([]GroupProjects, 0, len(projects))
	for _, v := range projects {
		ret = append(ret, GroupProjects{
			ID:   v.ID,
			Name: v.Name,
		})
	}
	return &ListGITProjectsResp{
		Projects: ret,
	}, nil
}
