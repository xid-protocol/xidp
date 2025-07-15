package v1

import (
	"crypto/sha256"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xid-protocol/xidp/protocols"
	"github.com/xid-protocol/xidp/protocols/whitelist"
)

type CreateWhitelistRequest struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func CreateWhitelist(c *gin.Context) {
	var req CreateWhitelistRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if req.Type == "awsOpenPort" {
		value, ok := req.Value.(whitelist.AWSOpenPort)
		if !ok {
			c.JSON(400, gin.H{"error": "Invalid value format for awsOpenPort type"})
			return
		}

		hash := sha256.Sum256([]byte(fmt.Sprintf("%+v", value)))
		sha256Value := fmt.Sprintf("%x", hash)
		xid := protocols.GenerateXid(value.InstanceID)

		wl := whitelist.Whitelist{
			Type:        req.Type,
			Value:       value,
			Sha256Value: sha256Value,
		}

		whitelist, err := whitelist.NewWhitelist(xid, req.Type, wl)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"whitelist": whitelist})
		return
	}

	c.JSON(400, gin.H{"error": "Unsupported type"})
}
