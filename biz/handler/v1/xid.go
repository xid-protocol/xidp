package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xid-protocol/xidp/common"
	"github.com/xid-protocol/xidp/internal"
)

func GetXid(c *gin.Context) {
	//获取get请求参数
	username := c.Query("username")
	source := c.Query("source")
	//如果不为空
	if username != "" {
		xid, err := internal.GetXid(username, source)
		if err != nil {
			c.JSON(200, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"xid": xid.Xid})
	} else {
		c.JSON(200, gin.H{"error": "username and source are required"})
	}

	var req map[string]interface{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}
	if req["path"] == nil {
		c.JSON(200, gin.H{"error": "path is required"})
		return
	}

}

func CreateXid(c *gin.Context) {
	//获取body里的text参数，json格式
	var req map[string]interface{}
	err := c.BindJSON(&req)
	xid := common.GenerateXid(req["text"].(string))
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}
	msg := map[string]interface{}{
		"name":    "xid-protocol",
		"xid":     xid,
		"version": "0.1.0",
	}
	c.JSON(200, gin.H{"msg": msg})
}
