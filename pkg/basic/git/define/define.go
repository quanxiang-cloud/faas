package define

import "time"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	State    string `json:"state"`
}

type SSHKey struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	CreatedAt *time.Time `json:"created_at"`
}

type Group struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	Path                 string `json:"path"`
	Description          string `json:"description"`
	MembershipLock       bool   `json:"membership_lock"`
	Visibility           string `json:"visibility"`
	LFSEnabled           bool   `json:"lfs_enabled"`
	AvatarURL            string `json:"avatar_url"`
	WebURL               string `json:"web_url"`
	RequestAccessEnabled bool   `json:"request_access_enabled"`
	FullName             string `json:"full_name"`
	FullPath             string `json:"full_path"`
}

type Project struct {
	ID             int        `json:"id"`
	Description    string     `json:"description"`
	Public         bool       `json:"public"`
	Visibility     string     `json:"visibility"`
	Name           string     `json:"name"`
	Path           string     `json:"path"`
	SSHURLToRepo   string     `json:"sshUrl"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	LastActivityAt *time.Time `json:"last_activity_at,omitempty"`
	CreatorID      int        `json:"creator_id"`
}
