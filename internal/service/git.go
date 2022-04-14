package service

import (
	"context"
	"faas/pkg/git"
)

type GITService interface {
	CreateUser(context.Context, *CreateUserReq) (*CreateUserResp, error)
	GetUser(context.Context, *GetUserReq) (*GetUserResp, error)

	AddSSHKey(context.Context, *AddSSHKeyReq) (*AddSSHKeyResp, error)
	ListSSHKey(context.Context, *ListSSHKeyReq) (*ListSSHKeyResp, error)
	DelSSHKey(context.Context, *DelSSHKeyReq) (*DelSSHKeyResp, error)

	CreateGroup(context.Context, *CreateGroupReq) (*CreateGroupResp, error)
	GetGroupByName(context.Context, *GetGroupByNameReq) (*GetGroupByNameResp, error)
	AddGroupMember(context.Context, *AddGroupMemberReq) (*AddGroupMemberResp, error)
	RemoveGroupMember(context.Context, *RemoveGroupMemberReq) (*RemoveGroupMemberResp, error)

	CreateProject(context.Context, *CreateProjectReq) (*CreateProjectResp, error)
	GetProjectByName(context.Context, *GetProjectByNameReq) (*GetProjectByNameResp, error)
	GetProjectByID(context.Context, *GetProjectByIDReq) (*GetProjectByIDResp, error)
}

type gitService struct {
	git git.GIT
}

func NewGIT() GITService {
	return &gitService{}
}

type CreateUserReq struct {
	UserID   string `json:"-"`
	UserName string `json:"-"`

	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type CreateUserResp struct {
}

func (g *gitService) CreateUser(context.Context, *CreateUserReq) (*CreateUserResp, error) {

}

type GetUserReq struct {
	UserName string `json:"-"`
}

type GetUserResp struct {
	ID      string `json:"id"`
	UserGID int    `json:"userGid"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	State   string `json:"state"`
}

func (g *gitService) GetUser(context.Context, *GetUserReq) (*GetUserResp, error) {

}

type AddSSHKeyReq struct {
	// UserGid string `json:"userGid"`
	Title string `json:"title"`
	Key   string `json:"key"`
}

type AddSSHKeyResp struct {
}

func (g *gitService) AddSSHKey(context.Context, *AddSSHKeyReq) (*AddSSHKeyResp, error) {

}

type ListSSHKeyReq struct {
	// UserGID string `json:"userGid"`
}

type ListSSHKeyResp struct {
	KeyID    string `json:"keyID"`
	Title    string `json:"title"`
	CreateAt string `json:"createAt"`
}

func (g *gitService) ListSSHKey(context.Context, *ListSSHKeyReq) (*ListSSHKeyResp, error) {

}

type DelSSHKeyReq struct {
	// UserGid string `json:"userGid"`
	KeyID string `json:"keyID"`
}

type DelSSHKeyResp struct {
}

func (g *gitService) DelSSHKey(context.Context, *DelSSHKeyReq) (*DelSSHKeyResp, error) {

}

type CreateGroupReq struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type CreateGroupResp struct {
}

func (g *gitService) CreateGroup(context.Context, *CreateGroupReq) (*CreateGroupResp, error) {

}

type GetGroupByNameReq struct {
	FullName string `json:"fullName"`
}

type GetGroupByNameResp struct {
}

func (g *gitService) GetGroupByName(context.Context, *GetGroupByNameReq) (*GetGroupByNameResp, error) {

}

type AddGroupMemberReq struct {
	// UserGID int `json:"userGid"`
	GroupID int `json:"groupID"`
}

type AddGroupMemberResp struct {
}

func (g *gitService) AddGroupMember(context.Context, *AddGroupMemberReq) (*AddGroupMemberResp, error) {

}

type RemoveGroupMemberReq struct {
	GroupID int `json:"groupid"`
	UserGID int `json:"userGid"`
}

type RemoveGroupMemberResp struct {
}

func (g *gitService) RemoveGroupMember(context.Context, *RemoveGroupMemberReq) (*RemoveGroupMemberResp, error) {

}

type CreateProjectReq struct {
	NamespaceID int    `json:"namespaceID"`
	Name        string `json:"name"`
}

type CreateProjectResp struct {
}

func (g *gitService) CreateProject(context.Context, *CreateProjectReq) (*CreateProjectResp, error) {

}

type GetProjectByNameReq struct {
	FullName string `json:"fullName"`
}

type GetProjectByNameResp struct {
}

func (g *gitService) GetProjectByName(context.Context, *GetProjectByNameReq) (*GetProjectByNameResp, error) {

}

type GetProjectByIDReq struct {
	ProjectID string `json:"projectID"`
}

type GetProjectByIDResp struct {
}

func (g *gitService) GetProjectByID(context.Context, *GetProjectByIDReq) (*GetProjectByIDResp, error) {

}
