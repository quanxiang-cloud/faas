package mysql

import (
	"context"
	"github.com/quanxiang-cloud/faas/internal/models"
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

func (g *functionRepo) Delete(ctx context.Context, tx *gorm.DB, id ...string) error {
	return tx.Where("id in (?)", id).Delete(&models.Function{}).Error
}

func (g *functionRepo) Get(ctx context.Context, db *gorm.DB, id string) *models.Function {
	one := models.Function{}
	affected := db.Where("id=?", id).Find(&one).RowsAffected
	if affected == 1 {
		return &one
	}
	return nil
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
