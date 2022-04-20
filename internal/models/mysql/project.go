package mysql

import (
	"github.com/quanxiang-cloud/faas/internal/models"
	"gorm.io/gorm"
)

type projectRepo struct {
}

func NewProjectRepo() models.ProjectRepo {
	return &projectRepo{}
}

func (g *projectRepo) getTable(db *gorm.DB) *gorm.DB {
	return db.Table("t_project")
}

func (g *projectRepo) Insert(db *gorm.DB, project *models.Project) error {
	return g.getTable(db).Create(project).Error
}

func (g *projectRepo) Del(db *gorm.DB, id string) error {
	return g.getTable(db).Where("id = ?", id).Delete(&models.Project{}).Error
}

func (g *projectRepo) Get(db *gorm.DB, id string) (*models.Project, error) {
	res := &models.Project{}
	err := g.getTable(db).Where("id = ?", id).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (g *projectRepo) GetByGroup(db *gorm.DB, alias, groupID string, page, limit int) ([]*models.Project, int64, error) {
	res := make([]*models.Project, 0)
	db = g.getTable(db)
	if alias != "" {
		db = db.Where("alias like ?", "%"+alias+"%")
	}
	db = db.Where("group_id = ?", groupID)
	db = db.Order("updated_at desc")
	var count int64
	db = db.Count(&count)
	db.Limit(limit).Offset((page - 1) * limit)
	err := db.Find(&res).Error
	if err != nil {
		return nil, 0, err
	}
	return res, count, nil
}

func (g *projectRepo) UpdDescribe(db *gorm.DB, id, describe string) error {
	return g.getTable(db).Where("id = ?", id).Update("describe", describe).Error
}
