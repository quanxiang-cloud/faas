package mysql

import (
	"github.com/quanxiang-cloud/faas/internal/models"
	"gorm.io/gorm"
)

type userGroupRepo struct {
}

func NewUserGroupRepo() models.UserGroupRepo {
	return &userGroupRepo{}
}

func (u *userGroupRepo) getTable(db *gorm.DB) *gorm.DB {
	return db.Table("user_group")
}

func (u *userGroupRepo) Insert(db *gorm.DB, userGroup *models.UserGroup) error {
	return u.getTable(db).Create(userGroup).Error
}

func (u *userGroupRepo) GetByUserID(db *gorm.DB, userID string) (*models.UserGroup, error) {
	res := &models.UserGroup{}
	db = u.getTable(db).Where("user_id = ?", userID).Find(&res)
	err := db.Error
	if err != nil {
		return nil, err
	}
	if db.RowsAffected <= 0 {
		return nil, nil
	}
	return res, nil
}

func (u *userGroupRepo) GetByUserGroup(db *gorm.DB, userID, groupID string) (*models.UserGroup, error) {
	res := &models.UserGroup{}
	db = u.getTable(db).Where("user_id = ?", userID).Where("group_id=?", groupID).Find(&res)
	err := db.Error
	if err != nil {
		return nil, err
	}
	if db.RowsAffected <= 0 {
		return nil, nil
	}
	return res, nil
}
