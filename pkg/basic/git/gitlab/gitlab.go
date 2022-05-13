package gitlab

import (
	"context"
	"fmt"
	"net/http"

	"github.com/quanxiang-cloud/faas/pkg/basic/git/define"

	"github.com/xanzy/go-gitlab"
)

// TODO: init repo
// TODO: bind project
type GIT struct {
	git *gitlab.Client
}

func New(config *define.Config) (*GIT, error) {
	git, err := gitlab.NewClient(config.Token,
		gitlab.WithBaseURL(config.BaseURL),
	)
	if err != nil {
		return nil, err
	}

	return &GIT{
		git: git,
	}, nil
}

func (g *GIT) CreateUser(ctx context.Context, email, password, username, name string, skipConfirmation bool) (*define.User, error) {

	user, _, err := g.git.Users.CreateUser(&gitlab.CreateUserOptions{
		Email:            &email,
		Username:         &username,
		Name:             &name,
		Password:         &password,
		SkipConfirmation: &skipConfirmation,
	})
	if err != nil {
		return nil, err
	}
	return &define.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
		State:    user.State,
	}, nil
}

func (g *GIT) AddSSHKeyForUser(ctx context.Context, userID int, title, key string) error {
	_, _, err := g.git.Users.AddSSHKeyForUser(userID, &gitlab.AddSSHKeyOptions{
		Title: &title,
		Key:   &key,
	})
	return err
}

func (g *GIT) ListSSHKey(ctx context.Context, userID int) ([]*define.SSHKey, error) {
	sshKey, _, err := g.git.Users.ListSSHKeysForUser(userID, &gitlab.ListSSHKeysForUserOptions{})
	if err != nil {
		return nil, err
	}

	keys := make([]*define.SSHKey, 0, len(sshKey))
	for _, key := range sshKey {
		keys = append(keys, &define.SSHKey{
			ID:        key.ID,
			Title:     key.Title,
			CreatedAt: key.CreatedAt,
		})
	}
	return keys, nil
}

func (g *GIT) DeleteSSHKeyFromUser(ctx context.Context, userID, keyID int) error {
	_, err := g.git.Users.DeleteSSHKeyForUser(userID, keyID)
	return err
}

func (g *GIT) GetUser(ctx context.Context, username string) (*define.User, error) {
	users, _, err := g.git.Users.ListUsers(&gitlab.ListUsersOptions{
		Username: &username,
	})
	if err != nil {
		return nil, err
	}
	switch len(users) {
	case 0:
		return nil, nil
	case 1:
		user := users[0]
		return &define.User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Name:     user.Name,
			State:    user.State,
		}, nil
	}

	return nil, fmt.Errorf("user conflict")
}

func (g *GIT) CreateGroup(ctx context.Context, name, path string) (*define.Group, error) {
	visibility := gitlab.PrivateVisibility
	group, _, err := g.git.Groups.CreateGroup(&gitlab.CreateGroupOptions{
		Name:       &name,
		Path:       &path,
		Visibility: &visibility,
	})
	if err != nil {
		return nil, err
	}

	_ = group

	return &define.Group{
		ID:             group.ID,
		Name:           group.Name,
		Path:           group.Path,
		Description:    group.Description,
		MembershipLock: group.MembershipLock,
		Visibility:     string(group.Visibility),
		FullName:       group.FullName,
		FullPath:       group.FullPath,
	}, nil
}

func (g *GIT) ListGroup(ctx context.Context) ([]*define.Group, error) {
	groups, _, err := g.git.Groups.ListGroups(&gitlab.ListGroupsOptions{})
	if err, ok := err.(*gitlab.ErrorResponse); ok {
		if err.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
	}
	resp := make([]*define.Group, 0, len(groups))
	for _, group := range groups {
		resp = append(resp, &define.Group{
			ID:             group.ID,
			Name:           group.Name,
			Path:           group.Path,
			Description:    group.Description,
			MembershipLock: group.MembershipLock,
			Visibility:     string(group.Visibility),
			FullName:       group.FullName,
			FullPath:       group.FullPath,
		})
	}
	return resp, nil
}

func (g *GIT) GetGroupByName(ctx context.Context, name string) (*define.Group, error) {
	group, _, err := g.git.Groups.GetGroup(name, &gitlab.GetGroupOptions{})
	if err, ok := err.(*gitlab.ErrorResponse); ok {
		if err.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
	}
	return &define.Group{
		ID:             group.ID,
		Name:           group.Name,
		Path:           group.Path,
		Description:    group.Description,
		MembershipLock: group.MembershipLock,
		Visibility:     string(group.Visibility),
		FullName:       group.FullName,
		FullPath:       group.FullPath,
	}, err
}

func (g *GIT) GetGroupByID(ctx context.Context, gid int) (*define.Group, error) {
	group, _, err := g.git.Groups.GetGroup(gid, &gitlab.GetGroupOptions{})
	if err, ok := err.(*gitlab.ErrorResponse); ok {
		if err.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
	}
	return &define.Group{
		ID:             group.ID,
		Name:           group.Name,
		Path:           group.Path,
		Description:    group.Description,
		MembershipLock: group.MembershipLock,
		Visibility:     string(group.Visibility),
		FullName:       group.FullName,
		FullPath:       group.FullPath,
	}, err
}

func (g *GIT) AddGroupMember(ctx context.Context, gid, userID int) error {
	accessLevel := gitlab.MaintainerPermissions
	_, _, err := g.git.GroupMembers.AddGroupMember(gid, &gitlab.AddGroupMemberOptions{
		UserID:      &userID,
		AccessLevel: &accessLevel,
	})
	return err
}

func (g *GIT) RemoveGroupMember(ctx context.Context, gid, userID int) error {
	if _, err := g.git.GroupMembers.RemoveGroupMember(gid, userID); err != nil {
		return err
	}
	return nil
}

func (g *GIT) CreateProject(ctx context.Context, name string, namespaceID int) (*define.Project, error) {
	project, _, err := g.git.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:        &name,
		Path:        &name,
		NamespaceID: &namespaceID,
	})
	if err != nil {
		return nil, err
	}
	return &define.Project{
		ID:             project.ID,
		Description:    project.Description,
		Public:         project.Public,
		Visibility:     string(project.Visibility),
		Name:           project.Name,
		Path:           project.Path,
		SSHURLToRepo:   project.SSHURLToRepo,
		CreatedAt:      project.CreatedAt,
		LastActivityAt: project.LastActivityAt,
		CreatorID:      project.CreatorID,
	}, nil
}

func (g *GIT) GetProjectByName(ctx context.Context, name string) (*define.Project, error) {
	return g.getProject(ctx, name)
}

func (g *GIT) GetProjectByID(ctx context.Context, id int) (*define.Project, error) {
	return g.getProject(ctx, id)
}

func (g *GIT) getProject(ctx context.Context, pid interface{}) (*define.Project, error) {
	project, _, err := g.git.Projects.GetProject(pid, &gitlab.GetProjectOptions{})
	if err, ok := err.(*gitlab.ErrorResponse); ok {
		if err.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
	}
	if err != nil {
		return nil, err
	}

	return &define.Project{
		ID:             project.ID,
		Description:    project.Description,
		Public:         project.Public,
		Visibility:     string(project.Visibility),
		Name:           project.Name,
		Path:           project.Path,
		SSHURLToRepo:   project.SSHURLToRepo,
		CreatedAt:      project.CreatedAt,
		LastActivityAt: project.LastActivityAt,
		CreatorID:      project.CreatorID,
	}, nil
}

func (g *GIT) GetGroupProjects(ctx context.Context, gid interface{}) ([]*define.Project, error) {
	projects, _, err := g.git.Groups.ListGroupProjects(gid, &gitlab.ListGroupProjectsOptions{})
	if err, ok := err.(*gitlab.ErrorResponse); ok {
		if err.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
	}
	if err != nil {
		return nil, err
	}
	projectsResp := make([]*define.Project, 0, len(projects))
	for _, project := range projects {
		projectsResp = append(projectsResp, &define.Project{
			ID:             project.ID,
			Description:    project.Description,
			Public:         project.Public,
			Visibility:     string(project.Visibility),
			Name:           project.Name,
			Path:           project.Path,
			CreatedAt:      project.CreatedAt,
			LastActivityAt: project.LastActivityAt,
			CreatorID:      project.CreatorID,
		})
	}
	return projectsResp, nil
}

func (g *GIT) CreateUserToken(ctx context.Context, uid int) (string, error) {

	defaultTokenName := "qxp-op-tokeb"
	tokenScopes := make([]string, 0)
	tokenScopes = append(tokenScopes, "api")
	tokenScopes = append(tokenScopes, "read_user")
	token, _, err := g.git.Users.CreateImpersonationToken(uid, &gitlab.CreateImpersonationTokenOptions{
		Name:   &defaultTokenName,
		Scopes: &tokenScopes,
	})
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

func (g *GIT) CreateFile(ctx context.Context, pid int, fullName, content, branch, commit string) error {
	_, _, err := g.git.RepositoryFiles.CreateFile(pid, fullName, &gitlab.CreateFileOptions{
		Branch:        &branch,
		CommitMessage: &commit,
		Content:       &content,
	})
	return err
}
