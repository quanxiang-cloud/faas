package models

import "gorm.io/gorm"

const (
	ProjectSuccessStatus = 1
	ProjectFailureStatus = -1
	ProjectUnknownStatus = 0
)

var ProjectStatus = map[int]string{
	ProjectSuccessStatus: "True",
	ProjectFailureStatus: "False",
	ProjectUnknownStatus: "Unknown",
}

type Project struct {
	ID          string
	GroupID     string
	ProjectID   int
	Alias       string
	ProjectName string
	RepoUrl     string
	Language    string
	Version     string
	Status      int
	Describe    string
	UserID      string
	CreatedAt   int64
	UpdatedAt   int64
	CreatedBy   string
	UpdatedBy   string
	DeletedBy   string
}

type ProjectRepo interface {
	Insert(db *gorm.DB, project *Project) error
	Del(db *gorm.DB, id string) error
	Get(db *gorm.DB, id string) (*Project, error)
	GetByGroup(db *gorm.DB, alias, group string, page, limit int) ([]*Project, int64, error)
	UpdDescribe(db *gorm.DB, id, describe string) error
}
