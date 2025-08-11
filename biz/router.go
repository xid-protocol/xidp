package biz

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/xid-protocol/xidp/biz/handler/v1"
	"github.com/xid-protocol/xidp/protocols/task"
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
		//sha1
		apiv1Group.POST("/sha1", v1.CreateSHA1)

		//notify
		notifyGroup := apiv1Group.Group("/notify")
		{
			notifyGroup.POST("/lark", func(c *gin.Context) {
				message := c.Query("message")
				v1.SendToLark(message)
			})
		}

		protocolGroup := apiv1Group.Group("/protocols")
		{
			attackSurface := protocolGroup.Group("/attack-surface")
			attackSurface.GET("/list", v1.GetAttackSurface)

			//task protocols
			taskGroup := protocolGroup.Group("/task")
			taskGroup.POST("/create", task.CreateTaskHandler)

		}
		whitelistGroup := protocolGroup.Group("/whitelist")
		{
			whitelistGroup.POST("/create", v1.CreateWhitelist)
		}
	}
}
