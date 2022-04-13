package mysql

import (
	"context"

	"gorm.io/gorm"

	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/faas/internal/models"
)

type gitRepo struct {
}

func NewGitRepo() models.GitRepo {
	return &gitRepo{}
}

func (g *gitRepo) Insert(ctx context.Context, tx *gorm.DB, data *models.Git) error {
	return tx.Create(data).Error
}

func (g *gitRepo) Update(ctx context.Context, tx *gorm.DB, data *models.Git) error {
	return tx.Model(data).Updates(data).Error
}

func (g *gitRepo) Delete(ctx context.Context, tx *gorm.DB, id ...string) error {
	return tx.Where("id in (?)", id).Delete(&models.Git{}).Error
}

func (g *gitRepo) Get(ctx context.Context, db *gorm.DB) *models.Git {
	one := models.Git{}
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if tenantID == "" {
		db = db.Where("tenant_id=? or tenant_id is null", tenantID)
	} else {
		db = db.Where("tenant_id=?", tenantID)
	}
	affected := db.Find(&one).RowsAffected
	if affected == 1 {
		return &one
	}
	return nil
}
