package logic

import (
	"context"

	"gorm.io/gorm"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/faas/internal/models"
	"github.com/quanxiang-cloud/faas/internal/models/mysql"
	"github.com/quanxiang-cloud/faas/pkg/code"
)

type Git interface {
	Create(c context.Context, r *CreateGitRequest) (*CreateGitResponse, error)
	Update(c context.Context, r *UpdateGitRequest) (*UpdateGitResponse, error)
	Delete(c context.Context, r *DeleteGitRequest) (*DeleteGitResponse, error)
	Get(c context.Context, r *GetGitRequest) (*GetGitResponse, error)
}

type git struct {
	db      *gorm.DB
	gitRepo models.GitRepo
}

func NewGit(c context.Context, db *gorm.DB) Git {
	return &git{
		db:      db,
		gitRepo: mysql.NewGitRepo(),
	}
}

type CreateGitRequest struct {
	Host       string `json:"host"`
	KnownHosts string `json:"knownHosts"`
	Token      string `json:"token"`
	SSH        string `json:"ssh"`
}
type CreateGitResponse struct {
}

func (g *git) Create(c context.Context, r *CreateGitRequest) (*CreateGitResponse, error) {
	one := g.gitRepo.Get(c, g.db)
	if one != nil {
		return nil, error2.New(code.ErrDataExist)
	}
	data := &models.Git{}
	data.ID = id.ShortID(0)
	data.Host = r.Host
	data.Token = r.Token
	data.SSH = r.SSH
	data.KnownHosts = r.KnownHosts
	unix := time.NowUnix()
	data.CreatedAt = unix
	data.UpdatedAt = unix
	return &CreateGitResponse{}, g.gitRepo.Insert(c, g.db, data)
}

type UpdateGitRequest struct {
	ID    string `json:"id"`
	Host  string `json:"host"`
	Token string `json:"token"`
}
type UpdateGitResponse struct {
}

func (g *git) Update(c context.Context, r *UpdateGitRequest) (*UpdateGitResponse, error) {
	data := g.gitRepo.Get(c, g.db)
	if data == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	data.Host = r.Host
	data.Token = r.Token
	unix := time.NowUnix()
	data.UpdatedAt = unix
	return &UpdateGitResponse{}, g.gitRepo.Update(c, g.db, data)
}

type DeleteGitRequest struct {
	ID string `json:"id"`
}
type DeleteGitResponse struct {
}

func (g *git) Delete(c context.Context, r *DeleteGitRequest) (*DeleteGitResponse, error) {
	data := g.gitRepo.Get(c, g.db)
	if data == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	return &DeleteGitResponse{}, g.gitRepo.Delete(c, g.db, r.ID)
}

type GetGitRequest struct {
}
type GetGitResponse struct {
	ID        string `json:"id"`
	Host      string `json:"host"`
	Token     string `json:"token"`
	UpdatedAt int64  `json:"updatedAt"`
}

func (g *git) Get(c context.Context, r *GetGitRequest) (*GetGitResponse, error) {
	data := g.gitRepo.Get(c, g.db)
	if data == nil {
		return nil, nil
	}
	res := &GetGitResponse{
		ID:        data.ID,
		Host:      data.Host,
		Token:     data.Token,
		UpdatedAt: data.UpdatedAt,
	}
	return res, nil

}
