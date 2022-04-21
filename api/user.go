package restful

import (
	"context"
	"github.com/gin-gonic/gin"
	error2 "github.com/quanxiang-cloud/cabin/error"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/faas/internal/logic"
	"github.com/quanxiang-cloud/faas/pkg/code"
	"github.com/quanxiang-cloud/faas/pkg/config"
	"gorm.io/gorm"
)

const (
	_userID = "User-Id"
)

// UserAPI UserAPI
type UserAPI struct {
	userService logic.UserService
}

// NewUserAPI NewUserAPI
func NewUserAPI(ctx context.Context, conf *config.Config, db *gorm.DB) *UserAPI {
	return &UserAPI{
		userService: logic.NewUserService(ctx, db, conf),
	}
}

// CreateUser CreateUser
func (u *UserAPI) CreateUser(c *gin.Context) {
	req := &logic.CreateUserReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	ctx := ginheader.MutateContext(c)
	req.UserID = c.GetHeader(_userID)
	resp.Format(u.userService.CreateUser(ctx, req)).Context(c)
}

func (u *UserAPI) GetUser(c *gin.Context) {
	req := &logic.GetUserReq{}
	ctx := ginheader.MutateContext(c)
	req.UserID = c.GetHeader(_userID)
	resp.Format(u.userService.GetUser(ctx, req)).Context(c)
}

func (u *UserAPI) CheckUser(c *gin.Context) {
	req := &logic.CheckUserReq{
		UserID: c.GetHeader(_userID),
	}
	ctx := ginheader.MutateContext(c)
	resp.Format(u.userService.CheckUser(ctx, req)).Context(c)
}
