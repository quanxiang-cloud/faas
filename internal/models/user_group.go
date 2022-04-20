package models

import "gorm.io/gorm"

type UserGroup struct {
	ID        string
	UserID    string
	GitID     int
	GroupID   string
	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
	CreatedBy string
	UpdatedBy string
	DeletedBy string
}

type UserGroupRepo interface {
	Insert(db *gorm.DB, userGroup *UserGroup) error
	GetByUserID(db *gorm.DB, userID string) (*UserGroup, error)
	GetByUserGroup(db *gorm.DB, userID, groupID string) (*UserGroup, error)
}
