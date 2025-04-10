package notify

import (
	"github.com/gin-gonic/gin"
	"takeout/internal/control/notify"
)

type NotifyRouter struct {
	notify *gin.RouterGroup
}

func NewNotifyRouter(group *gin.RouterGroup) *NotifyRouter {
	return &NotifyRouter{
		notify: group,
	}
}

func (r *NotifyRouter) RegisterRouters() {
	r.notifyRouter()
}

func (r *NotifyRouter) notifyRouter() {
	// 支付成功回调
	notifyController := notify.NewNotifyController()
	r.notify.POST("/pay", notifyController.PaySuccess)
	r.notify.POST("/refund", notifyController.RefundSuccess)
}
