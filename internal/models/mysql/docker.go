package mysql

import (
	"context"

	"gorm.io/gorm"

	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/faas/internal/models"
)

type dockerRepo struct {
}

func NewDockerRepo() models.DockerRepo {
	return &dockerRepo{}
}

func (g *dockerRepo) Insert(ctx context.Context, tx *gorm.DB, data *models.Docker) error {
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	data.Name = tenantID + "-docker"
	return tx.Create(data).Error
}

func (g *dockerRepo) Update(ctx context.Context, tx *gorm.DB, data *models.Docker) error {
	return tx.Model(data).Updates(data).Error
}

func (g *dockerRepo) Delete(ctx context.Context, tx *gorm.DB, id ...string) error {
	return tx.Where("id in (?)", id).Delete(&models.Git{}).Error
}

func (g *dockerRepo) Get(ctx context.Context, db *gorm.DB) *models.Docker {
	one := models.Docker{}
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
