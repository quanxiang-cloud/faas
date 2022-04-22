package models

import "gorm.io/gorm"

type Group struct {
	ID        string
	GroupID   int
	GroupName string
	Describe  string
	CreatedAt int64
	UpdatedAt int64
	CreatedBy string
	UpdatedBy string
	DeletedBy string
	AppID     string
}

type GroupRepo interface {
	Insert(db *gorm.DB, group *Group) error
	Del(db *gorm.DB, id string) error
	Get(db *gorm.DB, id string) (*Group, error)
	GetByName(db *gorm.DB, name string) (*Group, error)
	GetByApp(db *gorm.DB, appID string) (*Group, error)
}
