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
			// NewXID
			xidGroup.POST("/create", v1.CreateXID)
			// 通过id获取xid
			xidGroup.POST("/get", v1.Getxid)
			// 通过xid获取info
			xidGroup.GET("/:xid/info/*path", v1.GetXidInfo)

		}
		notifyGroup := apiv1Group.Group("/notify")
		{
			notifyGroup.POST("/lark", func(c *gin.Context) {
				message := c.Query("message")
				v1.SendToLark(message)
			})
		}

		protocolGroup := apiv1Group.Group("/protocols")
		{
			externalAttackSurface := protocolGroup.Group("/attack-surface")
			externalAttackSurface.GET("/get", v1.GetAttackSurface)
			// externalAttackSurface.POST("/create", v1.CreateExternalAttackSurface)
			// externalAttackSurface.POST("/update", v1.UpdateExternalAttackSurface)
			// externalAttackSurface.POST("/delete", v1.DeleteExternalAttackSurface)

		}
	}
}
