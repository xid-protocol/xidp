package aichat

import (
	"fmt"
	"strings"

	"github.com/colin-404/logx"
)

func (tm *ThreadManager) SendAgentStart(threadID string, agent string, content string) {
	ssEvent := ChatEvent{
		ThreadID: threadID,
		Agent:    agent,
		Content:  fmt.Sprintf("开始执行: %s", content),
		Type:     "agent_start",
	}
	tm.SendToThread(threadID, ssEvent)
	fmt.Println("=======================================")
	logx.Infof("SendAgentStart: %s, %s", agent, content)
	//发送节点信息

}

func (tm *ThreadManager) SendStart(threadID string) {
	ssEvent := ChatEvent{
		ThreadID: threadID,
		Agent:    "system",
		Content:  "开始处理您的请求...",
		Type:     "processing",
	}
	tm.SendToThread(threadID, ssEvent)

}

func (tm *ThreadManager) SendToolResult(threadID string, agent string, content string) {
	ssEvent := ChatEvent{
		ThreadID: threadID,
		Agent:    agent,
		Content:  content,
		Type:     "tool_result",
	}
	tm.SendToThread(threadID, ssEvent)

}

func (tm *ThreadManager) SendToolRuning(threadID string, agent string, content string) {
	ssEvent := ChatEvent{
		ThreadID: threadID,
		Agent:    agent,
		Content:  content,
		Type:     "tool_running",
	}
	tm.SendToThread(threadID, ssEvent)

}

func (tm *ThreadManager) SendReasoning(threadID string, agent string, content string) {
	ssEvent := ChatEvent{
		ThreadID: threadID,
		Agent:    agent,
		Content:  content,
		Type:     "reasoning",
	}
	tm.SendToThread(threadID, ssEvent)
	// fmt.Println("Sending reasoning:", ssEvent)

}

func (tm *ThreadManager) SendAnswer(threadID string, agent string, content string) {
	ssEvent := ChatEvent{
		ThreadID: threadID,
		Agent:    agent,
		Content:  content,
		Type:     "answer",
	}
	// fmt.Println(ssEvent)
	tm.SendToThread(threadID, ssEvent)

}

func (tm *ThreadManager) SendImage(threadID string, agent string, url string, description string) {
	//转为{"url":url, "description":description}
	content := fmt.Sprintf("{\"url\":%s, \"description\":%s}", url, description)
	content = strings.ReplaceAll(content, "\"", "\\\"")
	ssEvent := ChatEvent{
		ThreadID: threadID,
		Agent:    agent,
		Content:  content,
		Type:     "image",
	}
	// fmt.Println(ssEvent)
	tm.SendToThread(threadID, ssEvent)

}

func (tm *ThreadManager) SendEnd(threadID string) {
	ssEvent := ChatEvent{
		ThreadID: threadID,
		Agent:    "system",
		Content:  "end",
		Type:     "end",
	}
	tm.SendToThread(threadID, ssEvent)

}

func (tm *ThreadManager) SendError(threadID string, err error) {
	ssEvent := ChatEvent{
		ThreadID: threadID,
		Agent:    "system",
		Content:  fmt.Sprintf("执行错误: %v", err),
		Type:     "end",
	}
	tm.SendToThread(threadID, ssEvent)

}
