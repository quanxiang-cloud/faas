package git

import (
	"context"

	"github.com/quanxiang-cloud/faas/pkg/basic/git/define"
)

type GIT interface {
	CreateUser(ctx context.Context, email, password, username, name string, skipConfirmation bool) (*define.User, error)
	GetUser(ctx context.Context, username string) (*define.User, error)
	CreateUserToken(ctx context.Context, uid int) (string, error)

	AddSSHKeyForUser(ctx context.Context, userID int, title, key string) error
	ListSSHKey(ctx context.Context, userID int) ([]*define.SSHKey, error)
	DeleteSSHKeyFromUser(ctx context.Context, userID, keyID int) error

	CreateGroup(ctx context.Context, name, path string) (*define.Group, error)
	GetGroupByName(ctx context.Context, name string) (*define.Group, error)
	GetGroupByID(ctx context.Context, gid int) (*define.Group, error)
	ListGroup(ctx context.Context) ([]*define.Group, error)
	AddGroupMember(ctx context.Context, gid, userID int) error
	RemoveGroupMember(ctx context.Context, gid, userID int) error

	CreateProject(ctx context.Context, name string, namespaceID int) (*define.Project, error)
	GetProjectByName(ctx context.Context, name string) (*define.Project, error)
	GetProjectByID(ctx context.Context, id int) (*define.Project, error)
	GetGroupProjects(ctx context.Context, gid interface{}) ([]*define.Project, error)

	CreateFile(ctx context.Context, pid int, fullName string, content string) error
}
