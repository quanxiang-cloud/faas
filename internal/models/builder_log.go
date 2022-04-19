package models

import (
	"context"
	"time"
)

// BuilderLog BuilderLog
type BuilderLog struct {
	BaseLog
	Kubernetes `json:"kubernetes"`
}

// BuilderLogRepo BuilderLogRepo
type BuilderLogRepo interface {
	Search(ctx context.Context, ID string, time time.Time, page, size int) ([]*BuilderLog, int64, error)
}
