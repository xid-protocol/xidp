package v1

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xid-protocol/xidp/db/repositories"
	"github.com/xid-protocol/xidp/internal"
	"github.com/xid-protocol/xidp/protocols"
)

func GetXidInfo(c *gin.Context) {
	xid := c.Param("xid")
	path := c.Param("path")
	//如果path以info开头
	if strings.HasPrefix(path, "info") {
		xidInfoRepository := repositories.NewXidInfoRepository()
		xidInfo, err := xidInfoRepository.FindOneByXidAndPath(context.Background(), xid, path)
		if err != nil {
			c.JSON(200, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"xid":  xid,
			"info": xidInfo,
		})
		return
	}
	xidInfoRepository := repositories.NewXidInfoRepository()
	xidInfo, err := xidInfoRepository.FindOneByXidAndPath(context.Background(), xid, path)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"xid":  xid,
		"info": xidInfo,
	})
}

func Getxid(c *gin.Context) {

	var req map[string]interface{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}

	id := req["id"].(string)

	xid := protocols.GenerateXid(id)

	c.JSON(200, gin.H{
		"xid": xid,
		"id":  id,
	})

	// //获取get请求参数
	// username := c.Query("username")
	// source := c.Query("source")
	// //如果不为空
	// if username != "" {
	// 	xid, err := internal.GetXid(username, source)
	// 	if err != nil {
	// 		c.JSON(200, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	c.JSON(200, gin.H{})
	// } else {
	// 	c.JSON(200, gin.H{"error": "username and source are required"})
	// }

	// var req map[string]interface{}
	// err := c.BindJSON(&req)
	// if err != nil {
	// 	c.JSON(200, gin.H{"error": err.Error()})
	// 	return
	// }
	// if req["path"] == nil {
	// 	c.JSON(200, gin.H{"error": "path is required"})
	// 	return
	// }

}

//传入的参数示例
// {
// 	"id": "1234567890",
// 	"info": {
// 		"type": "email",
// 		"encryption": false
// 	},
// 	"metadata": {
// 		"operation": "create",
// 		"path": "info",
// 		"contentType": "application/json"
// 	},
// 	"payload": {
// 		"name": "John Doe",
// 		"email": "john.doe@example.com"
// 	}
// }

func CreateXID(c *gin.Context) {
	var req map[string]interface{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// 验证并获取必需的 plainID 字段
	ID, ok := req["id"].(string)
	if !ok || ID == "" {
		c.JSON(400, gin.H{"error": "ID is required and must be a non-empty string"})
		return
	}

	// 验证并获取 info
	info, err := internal.ConvertXIDInfo(req["info"].(map[string]interface{}))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// metadata conversion
	metaMap, _ := req["metadata"].(map[string]interface{})
	meta, err := internal.MapToMetadata(metaMap)
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

	XID := protocols.NewXID(&info, &meta, payload)

	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"xid": XID})
}

// func CreateXidId(c *gin.Context) {
// 	//获取body里的text参数，json格式

// 	var req map[string]interface{}
// 	err := c.BindJSON(&req)
// 	if err != nil {
// 		c.JSON(200, gin.H{"error": err.Error()})
// 		return
// 	}

// 	id := common.GenerateId(req["text"].(string))

// 	c.JSON(200, gin.H{
// 		"xid":  id,
// 		"text": req["text"].(string),
// 	})
// }
