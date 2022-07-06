package models

import (
	"context"
	"gorm.io/gorm"
)

type Git struct {
	ID                string `gorm:"column:id;type:varchar(64);PRIMARY_KEY" json:"id"`
	Host              string `gorm:"column:host;type:varchar(200);" json:"host"`
	KnownHosts        string `gorm:"column:known_hosts;type:text;" json:"knownHosts"`
	KeyScanKnownHosts string `gorm:"column:key_scan_known_hosts;type:text;" json:"keyScanKnownHosts"`
	SSH               string `gorm:"column:ssh;type:text;" json:"ssh"`
	Token             string `gorm:"column:token;type:text;" json:"token"`
	Name              string `gorm:"column:name;type:varchar(200);" json:"name"`

	CreatedAt int64  `gorm:"column:created_at;type:bigint; " json:"createdAt,omitempty" `
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint; " json:"updatedAt,omitempty" `
	DeletedAt int64  `gorm:"column:deleted_at;type:bigint; " json:"deletedAt,omitempty" `
	CreatedBy string `gorm:"column:created_by;type:varchar(64); " json:"createdBy,omitempty"` //创建者
	UpdatedBy string `gorm:"column:updated_by;type:varchar(64); " json:"updatedBy,omitempty"` //创建者
	DeletedBy string `gorm:"column:deleted_by;type:varchar(64); " json:"deletedBy,omitempty"` //删除者
	TenantID  string `gorm:"column:tenant_id;type:varchar(64); " json:"tenantID"`             //租户id
}

// TableName table name
func (Git) TableName() string {
	return "gits"
}

type GitRepo interface {
	Insert(ctx context.Context, tx *gorm.DB, data *Git) error
	Update(ctx context.Context, tx *gorm.DB, data *Git) error
	Delete(ctx context.Context, tx *gorm.DB, id ...string) error
	Get(ctx context.Context, db *gorm.DB) *Git
}
