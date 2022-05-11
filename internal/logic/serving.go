package logic

import (
	"context"
	"encoding/json"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/faas/internal/models"
	"github.com/quanxiang-cloud/faas/internal/models/mysql"
	"github.com/quanxiang-cloud/faas/pkg/basic/define/code"
	"github.com/quanxiang-cloud/faas/pkg/basic/k8s"
	"github.com/quanxiang-cloud/faas/pkg/config"
	"gorm.io/gorm"
)

type Serving interface {
	Serve(ctx context.Context, req *ServeReq) (*ServeResp, error)
	OffLine(ctx context.Context, req *OffLineReq) (*OffLineResp, error)
}

type serving struct {
	log          logger.AdaptedLogger
	db           *gorm.DB
	conf         *config.Config
	functionRepo models.FunctionRepo
	dockerRepo   models.DockerRepo
	projectRepo  models.ProjectRepo
	groupRepo    models.GroupRepo
	k8sc         k8s.Client
}

func NewServing(db *gorm.DB, conf *config.Config) Serving {
	return &serving{
		db:           db,
		conf:         conf,
		functionRepo: mysql.NewFunctionRepo(),
		dockerRepo:   mysql.NewDockerRepo(),
		projectRepo:  mysql.NewProjectRepo(),
		groupRepo:    mysql.NewGroupRepo(),
		k8sc:         k8s.NewClient("serving"),
	}
}

type ServeReq struct {
	ID string `json:"id"`
}

type ServeResp struct {
}

func (s *serving) Serve(ctx context.Context, req *ServeReq) (*ServeResp, error) {
	fn := s.functionRepo.Get(ctx, s.db, req.ID)
	if fn.Status != int(StatusOK) && fn.Status != int(StatusOffline) {
		return nil, error2.New(code.ErrDataIllegal)
	}

	project, err := s.projectRepo.Get(s.db, fn.ProjectID)
	if err != nil {
		return nil, err
	}

	group, err := s.groupRepo.Get(s.db, fn.GroupID)
	if err != nil {
		return nil, err
	}

	imageRepo := s.dockerRepo.Get(ctx, s.db)
	if imageRepo == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}

	tx := s.db.Begin()
	fn.Status = int(StatusServing)
	if err := s.functionRepo.Update(ctx, tx, fn); err != nil {
		tx.Rollback()
		return nil, err
	}

	env := make(map[string]string)
	if fn.Env != "" {
		json.Unmarshal([]byte(fn.Env), &env)
	}

	err = s.k8sc.CreateServing(ctx, &k8s.Function{
		Version:   fn.Version,
		Project:   project.ProjectName,
		GroupName: group.GroupName,
		Docker: &k8s.Docker{
			NameSpace: imageRepo.NameSpace,
			Name:      imageRepo.Name,
			Host:      imageRepo.Host,
		},
		ENV: env,
	})
	if err != nil {
		tx.Rollback()
		s.log.Error(err, "create serving failed")
		return nil, err
	}
	tx.Commit()
	return &ServeResp{}, nil
}

type OffLineReq struct {
	ID string
}

type OffLineResp struct {
}

func (s *serving) OffLine(ctx context.Context, req *OffLineReq) (*OffLineResp, error) {
	fn := s.functionRepo.Get(ctx, s.db, req.ID)
	if fn.Status != int(StatusOnline) {
		return nil, error2.New(code.ErrDataIllegal)
	}

	fn.Status = int(StatusOffline)
	if err := s.functionRepo.Update(ctx, s.db, fn); err != nil {
		return nil, err
	}

	project, err := s.projectRepo.Get(s.db, fn.ProjectID)
	if err != nil {
		return nil, err
	}

	group, err := s.groupRepo.Get(s.db, fn.GroupID)
	if err != nil {
		return nil, err
	}

	tx := s.db.Begin()
	fn.Status = int(StatusOffline)
	if err := s.functionRepo.Update(ctx, tx, fn); err != nil {
		tx.Rollback()
		return nil, err
	}

	err = s.k8sc.DelServing(ctx, &k8s.Function{
		Version:   fn.Version,
		Project:   project.ProjectName,
		GroupName: group.GroupName,
	})
	if err != nil {
		tx.Rollback()
		s.log.Error(err, "delete serving failed")
		return nil, err
	}

	tx.Commit()
	return &OffLineResp{}, nil
}
