package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/faas/internal/logic"
	"gorm.io/gorm"
)

type Event struct {
	e logic.EventRecordor
}

func newEvent(db *gorm.DB) *Event {
	return &Event{
		e: logic.NewEvent(db),
	}
}

func (e *Event) GetEvent(c *gin.Context) {
	req := &logic.GetEventReq{}
	req.ID = c.Param("id")
	resp.Format(e.e.GetEvent(header.MutateContext(c), req)).Context(c)
}
