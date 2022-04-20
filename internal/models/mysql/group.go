package mysql

import (
	"github.com/quanxiang-cloud/faas/internal/models"
	"gorm.io/gorm"
)

type groupRepo struct {
}

func NewGroupRepo() models.GroupRepo {
	return &groupRepo{}
}

func (g *groupRepo) getTable(db *gorm.DB) *gorm.DB {
	return db.Table("t_group")
}

func (g *groupRepo) Insert(db *gorm.DB, group *models.Group) error {
	return g.getTable(db).Create(group).Error
}

func (g *groupRepo) GetByName(db *gorm.DB, name string) (*models.Group, error) {
	res := &models.Group{}
	err := g.getTable(db).Where("group_name = ?", name).Find(&res).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (g *groupRepo) Del(db *gorm.DB, id string) error {
	return g.getTable(db).Where("id = ?", id).Delete(&models.Group{}).Error
}

func (g *groupRepo) Get(db *gorm.DB, id string) (*models.Group, error) {
	res := &models.Group{}
	err := g.getTable(db).Where("id = ?", id).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}