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

type Docker interface {
	Create(c context.Context, r *CreateDockerRequest) (*CreateDockerResponse, error)
	Update(c context.Context, r *UpdateDockerRequest) (*UpdateDockerResponse, error)
	Delete(c context.Context, r *DeleteDockerRequest) (*DeleteDockerResponse, error)
	Get(c context.Context, r *GetDockerRequest) (*GetDockerResponse, error)
}

type docker struct {
	db         *gorm.DB
	dockerRepo models.DockerRepo
}

func NewDocker(c context.Context, db *gorm.DB) Docker {
	return &docker{
		db:         db,
		dockerRepo: mysql.NewDockerRepo(),
	}
}

type CreateDockerRequest struct {
	Host     string `json:"host"`
	UserName string `json:"userName"`
	Secret   string `json:"secret"`
}
type CreateDockerResponse struct {
}

func (g docker) Create(c context.Context, r *CreateDockerRequest) (*CreateDockerResponse, error) {
	data := &models.Docker{}
	data.ID = id.ShortID(0)
	data.UserName = r.UserName
	data.Host = r.Host
	data.Secret = r.Secret
	data.NameSpace = "faas/"
	unix := time.NowUnix()
	data.CreatedAt = unix
	data.UpdatedAt = unix
	return &CreateDockerResponse{}, g.dockerRepo.Insert(c, g.db, data)
}

type UpdateDockerRequest struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
	Host     string `json:"host"`
	Secret   string `json:"secret"`
}
type UpdateDockerResponse struct {
}

func (g docker) Update(c context.Context, r *UpdateDockerRequest) (*UpdateDockerResponse, error) {
	data := g.dockerRepo.Get(c, g.db)
	if data == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	data.Secret = r.Secret
	data.UserName = r.UserName
	data.Host = r.Host
	unix := time.NowUnix()
	data.UpdatedAt = unix
	return &UpdateDockerResponse{}, g.dockerRepo.Update(c, g.db, data)
}

type DeleteDockerRequest struct {
	ID string `json:"id"`
}
type DeleteDockerResponse struct {
}

func (g docker) Delete(c context.Context, r *DeleteDockerRequest) (*DeleteDockerResponse, error) {
	data := g.dockerRepo.Get(c, g.db)
	if data == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	return &DeleteDockerResponse{}, g.dockerRepo.Delete(c, g.db, r.ID)
}

type GetDockerRequest struct {
}
type GetDockerResponse struct {
	ID        string `json:"id"`
	UserName  string `json:"userName"`
	Host      string `json:"host"`
	Secret    string `json:"secret"`
	UpdatedAt int64  `json:"updatedAt"`
}

func (g docker) Get(c context.Context, r *GetDockerRequest) (*GetDockerResponse, error) {
	data := g.dockerRepo.Get(c, g.db)
	if data == nil {
		return nil, nil
	}
	res := &GetDockerResponse{
		ID:        data.ID,
		UserName:  data.UserName,
		Secret:    data.Secret,
		Host:      data.Host,
		UpdatedAt: data.UpdatedAt,
	}
	return res, nil

}
