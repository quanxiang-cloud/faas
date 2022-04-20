package mysql

import (
	"github.com/quanxiang-cloud/faas/internal/models"
	"gorm.io/gorm"
)

type userRepo struct {
}

func NewUserRepo() models.UserRepo {
	return &userRepo{}
}

func (u *userRepo) getTable(db *gorm.DB) *gorm.DB {
	return db.Table("t_user")
}

func (u *userRepo) Insert(db *gorm.DB, user *models.User) error {
	return u.getTable(db).Create(user).Error
}

func (u *userRepo) Delete(db *gorm.DB, id string) error {
	return u.getTable(db).Where("id = ?", id).Delete(&models.User{}).Error
}

func (u *userRepo) Get(db *gorm.DB, id string) (*models.User, error) {
	res := &models.User{}
	err := u.getTable(db).Where("id = ?", id).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *userRepo) GetByUserID(db *gorm.DB, userID string) (*models.User, error) {
	res := &models.User{}
	err := u.getTable(db).Where("user_id = ?", userID).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
