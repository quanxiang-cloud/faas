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

type GroupService interface {
	CreateGroup(ctx context.Context, req *CreateGroupReq) (*CreateGroupResp, error)
	BindingGroup(ctx context.Context, req *BindingGroupReq) (*BindingGroupResp, error)
	CheckGroup(ctx context.Context, req *CheckGroupReq) (*CheckGroupResp, error)
	CheckMember(ctx context.Context, req *CheckMemberReq) (*CheckMemberResp, error)
}

type groupService struct {
	db            *gorm.DB
	userRepo      models.UserRepo
	groupRepo     models.GroupRepo
	gitRepo       models.GitRepo
	userGroupRepo models.UserGroupRepo
}

func NewGroupService(ctx context.Context, db *gorm.DB, conf *config.Config) GroupService {
	return &groupService{
		db:            db,
		userRepo:      mysql.NewUserRepo(),
		groupRepo:     mysql.NewGroupRepo(),
		gitRepo:       mysql.NewGitRepo(),
		userGroupRepo: mysql.NewUserGroupRepo(),
	}
}

type CreateGroupReq struct {
	AppID  string `json:"appID"`
	Group  string `json:"group"`
	UserID string `json:"-"`
}

type CreateGroupResp struct {
}

func (g *groupService) CreateGroup(ctx context.Context, req *CreateGroupReq) (*CreateGroupResp, error) {
	tx := g.db.Begin()
	// get the git host in tenant
	gitHost := g.gitRepo.Get(ctx, tx)
	if gitHost == nil {
		// TODO
		tx.Rollback()
		return nil, nil
	}
	// get the git admin client
	client, err := git2.GetClient(git2.Gitlab, gitHost.Token, gitHost.Host)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// create the group
	group, err := client.CreateGroup(ctx, req.Group, req.Group)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	groupInfo, err := g.groupRepo.GetByName(tx, req.Group)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if groupInfo == nil {
		// group information is stored
		err = g.groupRepo.Insert(tx, &models.Group{
			ID:        id.StringUUID(),
			GroupName: group.Name,
			GroupID:   group.ID,
			Describe:  group.Description,
			CreatedBy: req.UserID,
			UpdatedBy: req.UserID,
			CreatedAt: time.NowUnix(),
			UpdatedAt: time.NowUnix(),
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	// TODO return group id
	return &CreateGroupResp{}, nil
}

// BindingGroupReq BindingGroupReq
type BindingGroupReq struct {
	UserID  string `json:"-"`
	GroupID string `json:"-"`
}

type BindingGroupResp struct {
}

func (g *groupService) BindingGroup(ctx context.Context, req *BindingGroupReq) (*BindingGroupResp, error) {
	tx := g.db.Begin()
	gitHost := g.gitRepo.Get(ctx, tx)
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
	user, err := g.userRepo.GetByUserID(tx, req.UserID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	group, err := g.groupRepo.Get(g.db, req.GroupID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = client.AddGroupMember(ctx, group.GroupID, user.GitID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// create associations between user and groups
	err = g.userGroupRepo.Insert(tx, &models.UserGroup{
		ID:        id.StringUUID(),
		UserID:    req.UserID,
		GitID:     user.GitID,
		GroupID:   req.GroupID,
		CreatedBy: req.UserID,
		UpdatedBy: req.UserID,
		CreatedAt: time.NowUnix(),
		UpdatedAt: time.NowUnix(),
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &BindingGroupResp{}, nil
}

type CheckGroupReq struct {
	Group string `json:"group"`
}

type CheckGroupResp struct {
	GroupID string `json:"groupID"`
}

func (g *groupService) CheckGroup(ctx context.Context, req *CheckGroupReq) (*CheckGroupResp, error) {
	group, err := g.groupRepo.GetByName(g.db, req.Group)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return &CheckGroupResp{
			GroupID: "",
		}, nil
	}
	return &CheckGroupResp{
		GroupID: group.ID,
	}, nil
}

type CheckMemberReq struct {
	Group  string `json:"group"`
	UserID string `json:"-"`
}

type CheckMemberResp struct {
	IsMember bool `json:"isMember"`
}

func (g *groupService) CheckMember(ctx context.Context, req *CheckMemberReq) (*CheckMemberResp, error) {
	group, err := g.groupRepo.GetByName(g.db, req.Group)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return &CheckMemberResp{
			IsMember: false,
		}, nil
	}
	userGroup, err := g.userGroupRepo.GetByUserGroup(g.db, req.UserID, group.ID)
	if err != nil {
		return nil, err
	}
	if userGroup == nil {
		return &CheckMemberResp{
			IsMember: false,
		}, nil
	}
	return &CheckMemberResp{
		IsMember: true,
	}, nil
}
