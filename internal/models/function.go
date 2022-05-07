package models

import (
	"context"

	"gorm.io/gorm"
)

type Function struct {
	ID          string `gorm:"column:id;type:varchar(64);PRIMARY_KEY" json:"id"`
	GroupID     string `gorm:"column:group_id;type:varchar(64);" json:"groupID"`
	ProjectID   string `gorm:"column:project_id;type:varchar(64);" json:"projectID"`
	Version     string `gorm:"column:version;type:varchar(200);" json:"version"`
	Describe    string `gorm:"column:describe;type:text;" json:"describe"`
	Status      int    `gorm:"column:status;type:varchar(200);" json:"status"`
	DocStatus   int
	Env         string `gorm:"column:env;type:text;" json:"env"`
	ResourceRef string `gorm:"column:resource_ref;type:varchar(200);" json:"resourceRef"`
	Name        string `gorm:"column:name;type:varchar(200);" json:"name"`

	CreatedAt int64  `gorm:"column:created_at;type:bigint; " json:"createdAt,omitempty" `
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint; " json:"updatedAt,omitempty" `
	DeletedAt int64  `gorm:"column:deleted_at;type:bigint; " json:"deletedAt,omitempty" `
	CreatedBy string `gorm:"column:created_by;type:varchar(64); " json:"createdBy,omitempty"` //创建者
	UpdatedBy string `gorm:"column:updated_by;type:varchar(64); " json:"updatedBy,omitempty"` //创建者
	DeletedBy string `gorm:"column:deleted_by;type:varchar(64); " json:"deletedBy,omitempty"` //删除者
	TenantID  string `gorm:"column:tenant_id;type:varchar(64); " json:"tenantID"`             //租户id
	BuiltAt   int64  `gorm:"column:built_at;type:bigint; " json:"builtAt,omitempty" `
}

// TableName table name
func (Function) TableName() string {
	return "functions"
}

type FunctionRepo interface {
	Insert(ctx context.Context, tx *gorm.DB, data *Function) error
	Update(ctx context.Context, tx *gorm.DB, data *Function) error
	UpdateDescribe(ctx context.Context, tx *gorm.DB, data *Function) error
	Delete(ctx context.Context, tx *gorm.DB, id string) error
	Get(ctx context.Context, db *gorm.DB, id string) *Function
	Search(ctx context.Context, db *gorm.DB, projectID, groupID string, page, limit int) ([]Function, int64)
	GetByName(ctx context.Context, db *gorm.DB, name string) *Function
	GetByResourceRef(ctx context.Context, db *gorm.DB, resourceRef string) *Function
}
