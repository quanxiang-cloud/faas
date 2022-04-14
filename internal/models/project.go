package models

import "gorm.io/gorm"

type Project struct {
	ID          string
	GroupID     int
	Project     string
	ProjectName string
	Describe    string
}

type ProjectRepo interface {
	Insert(db *gorm.DB, project *Project) error
	Delete(db *gorm.DB, id string) error
	Get(db *gorm.DB, id string) (*Project, error)
}
