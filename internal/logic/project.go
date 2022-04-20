package logic

import (
	"context"
	"github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/faas/internal/models"
	"github.com/quanxiang-cloud/faas/internal/models/mysql"
	"github.com/quanxiang-cloud/faas/pkg/config"
	git2 "github.com/quanxiang-cloud/faas/pkg/git"
	"gorm.io/gorm"
)

type ProjectService interface {
	CreateProject(ctx context.Context, req *CreateProjectReq) (*CreateProjectResp, error)
	GetProjectByID(ctx context.Context, req *GetProjectByIDReq) (*GetProjectByIDResp, error)
	GetList(ctx context.Context, req *GetListReq) (*GetListResp, error)
	DelProject(ctx context.Context, req *DelProjectReq) (*DelProjectResp, error)
	UpdateDescribe(ctx context.Context, req *UpdateDescribeReq) (*UpdateDescribeResp, error)
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
	Name     string `json:"name"`
	Alias    string `json:"alias"`
	Language string `json:"language"`
	Version  string `json:"version"`
	UserID   string `json:"-"`
}

type CreateProjectResp struct {
}

func (p *project) CreateProject(ctx context.Context, req *CreateProjectReq) (*CreateProjectResp, error) {
	tx := p.db.Begin()
	gitHost := p.gitRepo.Get(ctx, tx)
	if gitHost == nil {
		// TODO return err
		tx.Rollback()
		return nil, nil
	}
	client, err := git2.GetClient(git2.Gitlab, gitHost.Token, gitHost.Host)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	group, err := p.groupRepo.Get(p.db, req.GroupID)
	project, err := client.CreateProject(ctx, req.Name, group.GroupID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = p.projectRepo.Insert(tx, &models.Project{
		ID:          id.StringUUID(),
		ProjectID:   project.ID,
		ProjectName: project.Name,
		Alias:       req.Alias,
		Describe:    project.Description,
		Status:      models.ProjectSuccessStatus,
		GroupID:     req.GroupID,
		UserID:      req.UserID,
		CreatedBy:   req.UserID,
		UpdatedBy:   req.UserID,
		CreatedAt:   time.NowUnix(),
		UpdatedAt:   time.NowUnix(),
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &CreateProjectResp{}, nil
}

type GetProjectByIDReq struct {
	ProjectID string `json:"projectID"`
}

type GetProjectByIDResp struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	Status      string `json:"status"`
	Language    string `json:"language"`
	Version     string `json:"version"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
	CreatedBy   string `json:"creator"`
}

func (p *project) GetProjectByID(ctx context.Context, req *GetProjectByIDReq) (*GetProjectByIDResp, error) {
	project, err := p.projectRepo.Get(p.db, req.ProjectID)
	if err != nil {
		return nil, err
	}

	return &GetProjectByIDResp{
		ID:        project.ID,
		Name:      project.ProjectName,
		Alias:     project.Alias,
		Status:    models.ProjectStatus[project.Status],
		Language:  project.Language,
		Version:   project.Version,
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
	Status      string `json:"status"`
	Language    string `json:"language"`
	Version     string `json:"version"`
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
			ID:        project.ID,
			Name:      project.ProjectName,
			Alias:     project.Alias,
			Status:    models.ProjectStatus[project.Status],
			Language:  project.Language,
			Version:   project.Version,
			CreatedAt: project.CreatedAt,
			UpdatedAt: project.UpdatedAt,
			CreatedBy: project.CreatedBy,
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
	Describe  string `json:"describe"`
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
