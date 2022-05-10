package logic

import (
	"context"

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

type GroupService interface {
	CreateGroup(ctx context.Context, req *CreateGroupReq) (*CreateGroupResp, error)
	AddGroupMember(ctx context.Context, req *AddGroupMemberReq) (*AddGroupMemberResp, error)
	CheckGroup(ctx context.Context, req *CheckGroupReq) (*CheckGroupResp, error)
	CheckMember(ctx context.Context, req *CheckMemberReq) (*CheckMemberResp, error)
	BindGroup(ctx context.Context, req *BindGroupReq) (*BindGroupResp, error)
	ListGroup(ctx context.Context, req *ListGroupReq) (*ListGroupResp, error)
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

// CreateGroupReq CreateGroupReq
type CreateGroupReq struct {
	AppID    string `json:"appID"`
	Group    string `json:"name"`
	Describe string `json:"describe"`
	UserID   string `json:"-"`
}

// CreateGroupResp CreateGroupResp
type CreateGroupResp struct {
	GroupID string `json:"groupID"`
}

func (g *groupService) CreateGroup(ctx context.Context, req *CreateGroupReq) (*CreateGroupResp, error) {
	tx := g.db.Begin()
	groupExist, err := g.groupRepo.GetByApp(g.db, req.AppID)
	if err != nil {
		return nil, err
	}
	if groupExist != nil {
		return nil, error2.New(code.ErrDataExist)
	}

	// get the git host in tenant
	gitHost := g.gitRepo.Get(ctx, tx)
	if gitHost == nil {
		tx.Rollback()
		return nil, error2.New(code.ErrDataNotExist)
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
	groupInfo := &models.Group{
		ID:        id.StringUUID(),
		GroupName: group.Name,
		GroupID:   group.ID,
		Describe:  req.Describe,
		AppID:     req.AppID,
		CreatedBy: req.UserID,
		UpdatedBy: req.UserID,
		CreatedAt: time.NowUnix(),
		UpdatedAt: time.NowUnix(),
	}
	// group information is stored
	err = g.groupRepo.Insert(tx, groupInfo)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &CreateGroupResp{
		GroupID: groupInfo.ID,
	}, nil
}

// AddGroupMemberReq AddGroupMemberReq
type AddGroupMemberReq struct {
	UserID  string `json:"-"`
	GroupID string `json:"-"`
}

// AddGroupMemberResp AddGroupMemberResp
type AddGroupMemberResp struct {
}

func (g *groupService) AddGroupMember(ctx context.Context, req *AddGroupMemberReq) (*AddGroupMemberResp, error) {
	tx := g.db.Begin()
	groupExist, err := g.userGroupRepo.GetByUserGroup(g.db, req.UserID, req.GroupID)
	if err != nil {
		return nil, err
	}
	if groupExist != nil {
		return nil, error2.New(code.ErrDataExist)
	}
	gitHost := g.gitRepo.Get(ctx, tx)
	if gitHost == nil {
		tx.Rollback()
		return nil, error2.New(code.ErrDataNotExist)
	}
	user, err := g.userRepo.GetByUserID(tx, req.UserID)
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
	return &AddGroupMemberResp{}, nil
}

type CheckGroupReq struct {
	AppID string `form:"appID"`
}

type CheckGroupResp struct {
	GroupID string `json:"groupID"`
}

func (g *groupService) CheckGroup(ctx context.Context, req *CheckGroupReq) (*CheckGroupResp, error) {
	group, err := g.groupRepo.GetByApp(g.db, req.AppID)
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
	GroupID string `form:"groupID"`
	AppID   string `form:"appID"`
	UserID  string `json:"-"`
}

type CheckMemberResp struct {
	IsMember bool `json:"isMember"`
}

func (g *groupService) CheckMember(ctx context.Context, req *CheckMemberReq) (*CheckMemberResp, error) {
	group, err := g.groupRepo.Get(g.db, req.GroupID)
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

// BindGroupReq BindGroupReq
type BindGroupReq struct {
	GID      int    `json:"gid"`
	Describe string `json:"describe"`
	AppID    string `json:"appID"`
	UserID   string `json:"-"`
}

// BindGroupResp BindGroupResp
type BindGroupResp struct {
	GroupID string `json:"groupID"`
}

func (g *groupService) BindGroup(ctx context.Context, req *BindGroupReq) (*BindGroupResp, error) {
	tx := g.db.Begin()
	groupExist, err := g.groupRepo.GetByApp(g.db, req.AppID)
	if err != nil {
		return nil, err
	}
	if groupExist != nil {
		return nil, error2.New(code.ErrDataExist)
	}
	// get the git host in tenant
	gitHost := g.gitRepo.Get(ctx, tx)
	if gitHost == nil {
		tx.Rollback()
		return nil, error2.New(code.ErrDataNotExist)
	}
	// get the git admin client
	client, err := git2.GetClient(git2.Gitlab, gitHost.Token, gitHost.Host)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	group, err := client.GetGroupByID(ctx, req.GID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	groupInfo := &models.Group{
		ID:        id.StringUUID(),
		GroupName: group.Name,
		GroupID:   group.ID,
		Describe:  group.Description,
		AppID:     req.AppID,
		CreatedBy: req.UserID,
		UpdatedBy: req.UserID,
		CreatedAt: time.NowUnix(),
		UpdatedAt: time.NowUnix(),
	}
	// group information is stored
	err = g.groupRepo.Insert(tx, groupInfo)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &BindGroupResp{
		GroupID: groupInfo.ID,
	}, nil
}

// ListGroupReq ListGroupReq
type ListGroupReq struct {
	UserID string `json:"-"`
}

// ListGroupResp ListGroupResp
type ListGroupResp struct {
	Groups []*GroupVO `json:"groups"`
}

// GroupVO GroupVO
type GroupVO struct {
	GID      int    `json:"gid"`
	Name     string `json:"name"`
	Describe string `json:"describe"`
}

func (g *groupService) ListGroup(ctx context.Context, req *ListGroupReq) (*ListGroupResp, error) {
	gitHost := g.gitRepo.Get(ctx, g.db)
	if gitHost == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	user, err := g.userRepo.GetByUserID(g.db, req.UserID)
	if err != nil {
		return nil, err
	}
	// get the git admin client
	client, err := git2.GetClient(git2.Gitlab, user.Token, gitHost.Host)
	if err != nil {
		return nil, err
	}
	groups, err := client.ListGroup(ctx)
	resp := &ListGroupResp{
		Groups: make([]*GroupVO, 0, len(groups)),
	}
	for _, group := range groups {
		resp.Groups = append(resp.Groups, &GroupVO{
			GID:      group.ID,
			Name:     group.Name,
			Describe: group.Description,
		})
	}
	return resp, nil
}
