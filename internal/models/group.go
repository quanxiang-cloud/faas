package models

import "gorm.io/gorm"

type Group struct {
	ID        string
	GroupID   int
	GroupName string
	Describe  string
}

type GroupRepo interface {
	Insert(db *gorm.DB, group *Group) error
	Del(db *gorm.DB, id string) error
	Get(db *gorm.DB, id string) (*Group, error)
}
