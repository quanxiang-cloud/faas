package models

import "gorm.io/gorm"

const (
	ProjectSuccessStatus = 1
)

var ProjectStatus = map[int]string{
	ProjectSuccessStatus: "成功",
}

type Project struct {
	ID          string
	GroupID     string
	ProjectID   int
	Alias       string
	ProjectName string
	Language    string
	Version     string
	Status      int
	Describe    string
	UserID      string
	CreatedAt   int64
	UpdatedAt   int64
	DeletedAt   int64
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