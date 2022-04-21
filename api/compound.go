package restful

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/faas/internal/logic"
	"net/http"
)

// Compound Compound
type Compound struct {
	ps logic.PubSub
}

// NewCompoundAPI NewCompoundAPI
func NewCompoundAPI(ctx context.Context, redisClient redis.UniversalClient) *Compound {
	return &Compound{
		ps: logic.NewPubSub(ctx, redisClient),
	}
}

// Subscribe Subscribe
func (cu *Compound) Subscribe(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetHeader("User-Id")
	req := &logic.SubscribeReq{
		UserID: userID,
	}

	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp.Format(cu.ps.Subscribe(ctx, req)).Context(c)
}
