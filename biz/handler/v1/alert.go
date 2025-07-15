package v1

import (
	"github.com/colin-404/logx"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/xid-protocol/xidp/common"
)

// Lark消息结构
type LarkMessage struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

func Notify(c *gin.Context) {
	//从body里面获取message
	var req map[string]string
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}
	if req["method"] == "lark_custom_bot" {
		SendToLark(req["message"])
	} else {
		c.JSON(200, gin.H{"error": "method not supported"})
		return
	}
}

func SendToLark(message string) {
	webhookURL := viper.GetString("Notify.lark_custom_bot_webhook")
	logx.Infof("SendToLark: %s", webhookURL)

	// 创建Lark消息结构
	// larkMessage := LarkMessage{
	// 	MsgType: "text",
	// }
	common.DoHttp("POST", webhookURL, message, map[string]string{})
	// larkMessage.Content.Text = message
}
