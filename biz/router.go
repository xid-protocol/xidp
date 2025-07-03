package biz

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/xid-protocol/xidp/biz/handler/v1"
)

func RegisterRouter(r *gin.Engine) {
	apiv1Group := r.Group("/api/v1")
	{
		xidGroup := apiv1Group.Group("/xid")
		{
			xidGroup.GET("/get", v1.GetXid)
			xidGroup.POST("/create", v1.CreateXid)

		}
		notifyGroup := apiv1Group.Group("/notify")
		{
			notifyGroup.POST("/lark", v1.SendToLark)
		}
	}
}
