package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xid-protocol/xidp/common"
	"github.com/xid-protocol/xidp/protocols"
)

// func GetXid(c *gin.Context) {
// 	//获取get请求参数
// 	username := c.Query("username")
// 	source := c.Query("source")
// 	//如果不为空
// 	if username != "" {
// 		xid, err := internal.GetXid(username, source)
// 		if err != nil {
// 			c.JSON(200, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(200, gin.H{})
// 	} else {
// 		c.JSON(200, gin.H{"error": "username and source are required"})
// 	}

// 	var req map[string]interface{}
// 	err := c.BindJSON(&req)
// 	if err != nil {
// 		c.JSON(200, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if req["path"] == nil {
// 		c.JSON(200, gin.H{"error": "path is required"})
// 		return
// 	}

// }

func CreateXid(c *gin.Context) {
	//获取body里的text参数，json格式

	var req map[string]interface{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}

	id := common.GenerateId(req["text"].(string))

	// type validation
	typeStr, ok := req["type"].(string)
	if !ok {
		c.JSON(400, gin.H{"error": "type must be a string"})
		return
	}
	xidType := protocols.ConvertXIDType(typeStr)

	// metadata conversion
	metaMap, _ := req["metadata"].(map[string]interface{})
	meta, err := protocols.MapToMetadata(metaMap)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// payload
	payload := req["payload"].(map[string]interface{})
	if payload == nil {
		c.JSON(400, gin.H{"error": "payload is required"})
		return
	}

	xid := protocols.NewXID(id, xidType, meta, payload)

	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"xid": xid})
}
