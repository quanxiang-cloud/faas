package models

import (
	"context"
	"gorm.io/gorm"
)

type Function struct {
	ID        string `gorm:"column:id;type:varchar(64);PRIMARY_KEY" json:"id"`
	GroupName string `gorm:"column:group_name;type:varchar(200);" json:"groupName"`
	Project   string `gorm:"column:project;type:varchar(200);" json:"project"`
	Version   string `gorm:"column:version;type:varchar(200);" json:"version"`
	Language  string `gorm:"column:language;type:varchar(200);" json:"language"`
	Status    int    `gorm:"column:status;type:varchar(200);" json:"status"`
	Env       string `gorm:"column:env;type:text;" json:"env"`

	CreatedAt int64  `gorm:"column:created_at;type:bigint; " json:"createdAt,omitempty" `
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint; " json:"updatedAt,omitempty" `
	DeletedAt int64  `gorm:"column:deleted_at;type:bigint; " json:"deletedAt,omitempty" `
	CreatedBy string `gorm:"column:created_by;type:varchar(64); " json:"createdBy,omitempty"` //创建者
	UpdatedBy string `gorm:"column:updated_by;type:varchar(64); " json:"updatedBy,omitempty"` //创建者
	DeletedBy string `gorm:"column:deleted_by;type:varchar(64); " json:"deletedBy,omitempty"` //删除者
	TenantID  string `gorm:"column:tenant_id;type:varchar(64); " json:"tenantID"`             //租户id
}

// TableName table name
func (Function) TableName() string {
	return "functions"
}

type FunctionRepo interface {
	Insert(ctx context.Context, tx *gorm.DB, data *Function) error
	Update(ctx context.Context, tx *gorm.DB, data *Function) error
	Delete(ctx context.Context, tx *gorm.DB, id ...string) error
	Get(ctx context.Context, db *gorm.DB, id string) *Function
}
