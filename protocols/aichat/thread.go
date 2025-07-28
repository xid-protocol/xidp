package aichat

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	threadManager *ThreadManager
	once          sync.Once
)

// Backend -> Frontend 输出消息
type ChatEvent struct {
	ThreadID string `json:"threadID"`
	Agent    string `json:"agent,omitempty"`
	MsgID    string `json:"msgID,omitempty"` // 用于标识消息的唯一ID
	Content  string `json:"content,omitempty"`
	Type     string `json:"type"`
}

// ThreadManager 管理每个线程的独立通道
type ThreadManager struct {
	channels    map[string]chan ChatEvent     // threadID -> 响应通道
	ctxMap      map[string]context.Context    // threadID -> 上下文
	cancelMap   map[string]context.CancelFunc // threadID -> 取消函数
	status      map[string]string             // threadID -> 状态 (running/cancelled/completed)
	mu          sync.RWMutex
	ChatRequest map[string]ChatRequest
}

// var threadManager = &ThreadManager{
// 	channels:    make(map[string]chan SSEEvent),
// 	connections: make(map[string]*gin.Context),
// 	contexts:    make(map[string]context.CancelFunc),
// 	status:      make(map[string]string),
// }

func newThreadManager() *ThreadManager {
	return &ThreadManager{
		channels:    make(map[string]chan ChatEvent),
		cancelMap:   make(map[string]context.CancelFunc),
		ctxMap:      make(map[string]context.Context),
		status:      make(map[string]string),
		ChatRequest: make(map[string]ChatRequest),
	}
}

func ThreadMan() *ThreadManager {
	once.Do(func() {
		threadManager = newThreadManager()
	})
	return threadManager
}

// StartNewThread 开始新的thread
func (tm *ThreadManager) CreateThread(ctx context.Context, chatRequest ChatRequest) (chan ChatEvent, string) {
	threadID := uuid.NewString()

	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 创建响应通道
	ch := make(chan ChatEvent, 50)

	// 注册thread
	tm.channels[threadID] = ch
	tm.ctxMap[threadID] = ctx
	tm.status[threadID] = "running"
	tm.ChatRequest[threadID] = chatRequest

	return ch, threadID
}

// CancelThread 取消指定thread
func (tm *ThreadManager) CancelThread(threadID string) bool {
	tm.mu.RLock()
	cancel, ok := tm.cancelMap[threadID]
	tm.mu.RUnlock()
	if !ok {
		return false // 没找到这个 thread
	}

	// 触发 ctx.Done()，让后台 goroutine 立刻退出
	cancel()

	// （可选）推送一条取消事件，前端立即显示
	tm.SendToThread(threadID, ChatEvent{
		ThreadID: threadID,
		Agent:    "system",
		Content:  "用户已取消",
		Type:     "cancelled",
	})
	return true
}

// CompleteThread 标记thread完成
func (tm *ThreadManager) CompleteThread(threadID string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.status[threadID] == "running" {
		tm.status[threadID] = "completed"
	}
}

// CleanupThread 清理thread资源
func (tm *ThreadManager) CleanupThread(threadID string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if cancel, exists := tm.cancelMap[threadID]; exists {
		cancel()
		delete(tm.cancelMap, threadID)
	}

	if ch, exists := tm.channels[threadID]; exists {
		close(ch)
		delete(tm.channels, threadID)
	}

	delete(tm.status, threadID)
}

// SendToThread 发送消息到指定thread（检查是否已取消）
func (tm *ThreadManager) SendToThread(threadID string, event ChatEvent) bool {
	tm.mu.RLock()
	ch, exists := tm.channels[threadID]
	status := tm.status[threadID]
	tm.mu.RUnlock()

	if !exists || status == "cancelled" {
		return false
	}

	select {
	case ch <- event:
		return true
	case <-time.After(1 * time.Second):
		return false
	}
}

// GetThreadStatus 获取thread状态
func (tm *ThreadManager) GetThreadStatus(threadID string) string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.status[threadID]
}
