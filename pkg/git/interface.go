package git

import (
	"context"
	"faas/pkg/git/define"
)

type GIT interface {
	CreateUser(ctx context.Context, email, password, username, name string, skipConfirmation bool) (*define.User, error)
	GetUser(ctx context.Context, username string) (*define.User, error)

	AddSSHKeyForUser(ctx context.Context, userID int, title, key string) error
	ListSSHKey(ctx context.Context, userID int) ([]*define.SSHKey, error)
	DeleteSSHKeyFromUser(ctx context.Context, userID, keyID int) error

	CreateGroup(ctx context.Context, name, path string) (*define.Group, error)
	GetGroupByName(ctx context.Context, name string) (*define.Group, error)
	AddGroupMember(ctx context.Context, gid, userID int) error
	RemoveGroupMember(ctx context.Context, gid, userID int) error

	CreateProject(ctx context.Context, name string, namespaceID int) (*define.Project, error)
	GetProjectByName(ctx context.Context, name string) (*define.Project, error)
	GetProjectByID(ctx context.Context, id int) (*define.Project, error)
}
