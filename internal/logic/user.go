package logic

import (
	"context"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/faas/internal/models"
	"github.com/quanxiang-cloud/faas/internal/models/mysql"
	"github.com/quanxiang-cloud/faas/pkg/basic/define/code"
	git2 "github.com/quanxiang-cloud/faas/pkg/basic/git"
	"github.com/quanxiang-cloud/faas/pkg/config"
	"github.com/quanxiang-cloud/organizations/pkg/client"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(ctx context.Context, req *CreateUserReq) (*CreateUserResp, error)
	GetUser(ctx context.Context, req *GetUserReq) (*GetUserResp, error)
	CheckUser(ctx context.Context, req *CheckUserReq) (*CheckUserResp, error)
}

type userSerice struct {
	db         *gorm.DB
	userRepo   models.UserRepo
	gitRepo    models.GitRepo
	userClient client.User
}

func NewUserService(ctx context.Context, db *gorm.DB, conf *config.Config) UserService {
	return &userSerice{
		db:         db,
		userRepo:   mysql.NewUserRepo(),
		gitRepo:    mysql.NewGitRepo(),
		userClient: client.NewUser(conf.InternalNet),
	}
}

type CreateUserReq struct {
	Account string `json:"account"`
	UserID  string `json:"-"`
}

type CreateUserResp struct {
}

func (u *userSerice) CreateUser(ctx context.Context, req *CreateUserReq) (*CreateUserResp, error) {
	tx := u.db.Begin()
	userExist, err := u.userRepo.GetByUserID(u.db, req.UserID)
	if err != nil {
		return nil, err
	}
	if userExist != nil {
		return nil, error2.New(code.ErrDataExist)
	}
	gitHost := u.gitRepo.Get(ctx, tx)
	if gitHost == nil {
		tx.Rollback()
		return nil, error2.New(code.ErrDataNotExist)
	}
	gitClient, err := git2.GetClient(git2.Gitlab, gitHost.Token, gitHost.Host)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	user, err := gitClient.GetUser(ctx, req.Account)
	if err != nil || user == nil {
		tx.Rollback()
		return nil, error2.New(code.ErrDataNotExist)
	}

	token, err := gitClient.CreateUserToken(ctx, user.ID)
	if err != nil || token == "" {
		tx.Rollback()
		return nil, error2.New(code.ErrDataNotExist)
	}
	err = u.userRepo.Insert(tx, &models.User{
		ID:        id.StringUUID(),
		UserID:    req.UserID,
		GitName:   req.Account,
		GitID:     user.ID,
		Token:     token,
		CreatedAt: time.NowUnix(),
		UpdatedAt: time.NowUnix(),
		CreatedBy: req.UserID,
		UpdatedBy: req.UserID,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &CreateUserResp{}, nil
}

type GetUserReq struct {
	UserID string `json:"-"`
}

type GetUserResp struct {
	ID       int    `json:"id"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	State    string `json:"state"`
}

func (u *userSerice) GetUser(ctx context.Context, req *GetUserReq) (*GetUserResp, error) {
	gitHost := u.gitRepo.Get(ctx, u.db)
	if gitHost == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	user, err := u.userRepo.GetByUserID(u.db, req.UserID)
	if err != nil {
		return nil, err
	}
	gitClient, err := git2.GetClient(git2.Gitlab, gitHost.Token, gitHost.Host)
	if err != nil {
		return nil, err
	}
	gitUser, err := gitClient.GetUser(ctx, user.GitName)
	if err != nil {
		return nil, err
	}
	return &GetUserResp{
		ID:       gitUser.ID,
		UserName: gitUser.Username,
		Name:     gitUser.Name,
		Email:    gitUser.Email,
		State:    gitUser.State,
	}, nil
}

type CheckUserReq struct {
	UserID string `json:"-"`
}

type CheckUserResp struct {
	UserAccount string `json:"userAccount"`
}

func (u *userSerice) CheckUser(ctx context.Context, req *CheckUserReq) (*CheckUserResp, error) {
	gitHost := u.gitRepo.Get(ctx, u.db)
	if gitHost == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	gitClient, err := git2.GetClient(git2.Gitlab, gitHost.Token, gitHost.Host)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.GetByUserID(u.db, req.UserID)
	if err != nil || user == nil {
		return &CheckUserResp{
			UserAccount: "",
		}, nil
	}
	logger.Logger.Info(user.GitName)
	gitUser, err := gitClient.GetUser(ctx, user.GitName)
	if err != nil || gitUser == nil {
		return &CheckUserResp{
			UserAccount: "",
		}, nil
	}
	return &CheckUserResp{
		UserAccount: gitUser.Username,
	}, nil
}
