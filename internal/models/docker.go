package models

import (
	"context"
	"gorm.io/gorm"
)

type Docker struct {
	ID       string `gorm:"column:id;type:varchar(64);PRIMARY_KEY" json:"id"`
	Host     string `gorm:"column:host;type:varchar(200);" json:"host"`
	UserName string `gorm:"column:name;type:varchar(64);" json:"name"`
	Secret   string `gorm:"column:secret;type:text;" json:"secret"`

	CreatedAt int64  `gorm:"column:created_at;type:bigint; " json:"createdAt,omitempty" `
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint; " json:"updatedAt,omitempty" `
	DeletedAt int64  `gorm:"column:deleted_at;type:bigint; " json:"deletedAt,omitempty" `
	CreatedBy string `gorm:"column:created_by;type:varchar(64); " json:"createdBy,omitempty"` //创建者
	UpdatedBy string `gorm:"column:updated_by;type:varchar(64); " json:"updatedBy,omitempty"` //创建者
	DeletedBy string `gorm:"column:deleted_by;type:varchar(64); " json:"deletedBy,omitempty"` //删除者
	TenantID  string `gorm:"column:tenant_id;type:varchar(64); " json:"tenantID"`             //租户id
}

// TableName table name
func (Docker) TableName() string {
	return "docker"
}

type DockerRepo interface {
	Insert(ctx context.Context, tx *gorm.DB, data *Docker) error
	Update(ctx context.Context, tx *gorm.DB, data *Docker) error
	Delete(ctx context.Context, tx *gorm.DB, id ...string) error
	Get(ctx context.Context, db *gorm.DB) *Docker
}
