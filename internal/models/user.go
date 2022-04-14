package models

import "gorm.io/gorm"

type User struct {
	ID     string
	UserID string
	GitID  int
}

type UserRepo interface {
	Insert(db *gorm.DB, user *User) error
	Delete(db *gorm.DB, id string) error
	Get(db *gorm.DB, id string) (*User, error)
}
