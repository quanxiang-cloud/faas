package mysql

import (
	"context"
	"github.com/quanxiang-cloud/faas/internal/models"
	page2 "github.com/quanxiang-cloud/faas/pkg/page"
	"gorm.io/gorm"
)

type functionRepo struct {
}

func NewFunctionRepo() models.FunctionRepo {
	return &functionRepo{}
}

func (g *functionRepo) Insert(ctx context.Context, tx *gorm.DB, data *models.Function) error {
	return tx.Create(data).Error
}

func (g *functionRepo) Update(ctx context.Context, tx *gorm.DB, data *models.Function) error {
	return tx.Model(data).Updates(data).Error
}

func (g *functionRepo) Delete(ctx context.Context, tx *gorm.DB, id string) error {
	return tx.Where("id=?", id).Delete(&models.Function{}).Error
}

func (g *functionRepo) Get(ctx context.Context, db *gorm.DB, id string) *models.Function {
	one := models.Function{}
	affected := db.Where("id=?", id).Find(&one).RowsAffected
	if affected == 1 {
		return &one
	}
	return nil
}

func (g *functionRepo) Search(ctx context.Context, db *gorm.DB, projectID, groupID string, page, limit int) ([]models.Function, int64) {
	functions := make([]models.Function, 0)
	db = db.Where("group_id=? and project_id=?", groupID, projectID)
	var num int64
	db.Model(&models.Function{}).Count(&num)
	newPage := page2.NewPage(page, limit, num)

	db = db.Limit(newPage.PageSize).Offset(newPage.StartIndex).Order("created_at DESC")
	affected := db.Find(&functions).RowsAffected
	if affected > 0 {
		return functions, num
	}
	return nil, 0
}

func (g *functionRepo) GetByName(ctx context.Context, db *gorm.DB, name string) *models.Function {
	one := models.Function{}
	affected := db.Where("name=?", name).Find(&one).RowsAffected
	if affected == 1 {
		return &one
	}
	return nil
}

func (g *functionRepo) GetByResourceRef(ctx context.Context, db *gorm.DB, resourceRef string) *models.Function {
	one := models.Function{}
	affected := db.Where("resource_ref=?", resourceRef).Find(&one).RowsAffected
	if affected == 1 {
		return &one
	}
	return nil
}

func (g *functionRepo) UpdateDescribe(ctx context.Context, tx *gorm.DB, data *models.Function) error {
	return tx.Model(data).Updates(data).Update("describe", data.Describe).Error
}
