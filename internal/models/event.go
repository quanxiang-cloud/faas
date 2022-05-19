package models

import (
	"gorm.io/gorm"
)

type Event struct {
	ID       string
	Name     string
	Type     int
	State    string
	Message  string
	CreateBy string
	CreateAt int64
	UpdateAt int64
	DeleteAt int64
}

func (Event) TableName() string {
	return "event"
}

type EventRepo interface {
	Insert(db *gorm.DB, event *Event) error
	Update(db *gorm.DB, event *Event) error
	Query(db *gorm.DB, id string) (*Event, error)
	QueryByName(db *gorm.DB, name string) (*Event, error)
	Delete(db *gorm.DB, id string) error
}
