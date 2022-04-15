package mysql

import (
	"context"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
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
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if tenantID == "" {
		db = db.Where("tenant_id=? or tenant_id is null", tenantID)
	} else {
		db = db.Where("tenant_id=?", tenantID)
	}
	affected := db.Where("id=?", id).Find(&one).RowsAffected
	if affected == 1 {
		return &one
	}
	return nil
}
