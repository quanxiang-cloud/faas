package logic

import (
	"context"
	"time"

	"github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/faas/internal/models"
	"github.com/quanxiang-cloud/faas/internal/models/mysql"
	"github.com/quanxiang-cloud/faas/pkg/basic/event"
	"gorm.io/gorm"
)

type EventRecordor interface {
	Save(msg *event.MsgBus) error
	GetEvent(ctx context.Context, req *GetEventReq) (*GetEventResp, error)
}

type eventRecordor struct {
	db        *gorm.DB
	eventRepo models.EventRepo
}

func NewEvent(db *gorm.DB) EventRecordor {
	return &eventRecordor{
		db:        db,
		eventRepo: mysql.NewEventRepo(),
	}
}

func (e *eventRecordor) Save(msg *event.MsgBus) error {
	b := event.Convert(msg)
	event := &models.Event{
		ID:       id.ShortID(0),
		Name:     b.Name,
		Type:     int(msg.Type),
		State:    b.State,
		CreateAt: time.Now().Unix(),
	}

	msg.Data = event.ID
	return e.eventRepo.Insert(e.db, event)
}

type GetEventReq struct {
	ID string `json:"id"`
}

type GetEventResp struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     int    `json:"type"`
	State    string `json:"state"`
	CreateAt int64  `json:"createAt"`
}

func (e *eventRecordor) GetEvent(ctx context.Context, req *GetEventReq) (*GetEventResp, error) {
	event, err := e.eventRepo.Query(e.db, req.ID)
	if err != nil {
		return nil, err
	}

	return &GetEventResp{
		ID:       event.ID,
		Name:     event.Name,
		Type:     event.Type,
		State:    event.State,
		CreateAt: event.CreateAt,
	}, nil
}
