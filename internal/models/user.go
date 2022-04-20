package models

import "gorm.io/gorm"

type User struct {
	ID        string
	UserID    string
	GitID     int
	GitName   string
	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
	CreatedBy string
	UpdatedBy string
	DeletedBy string
}

type UserRepo interface {
	Insert(db *gorm.DB, user *User) error
	Delete(db *gorm.DB, id string) error
	Get(db *gorm.DB, id string) (*User, error)
	GetByUserID(db *gorm.DB, userID string) (*User, error)
}
