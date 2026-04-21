package memory

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// WorkingMemory 工作记忆（L1）- 会话级上下文存储
type WorkingMemory struct {
	sessionID string
	context   context.Context
	data      map[string]any
	mu        sync.RWMutex
	createdAt time.Time
	updatedAt time.Time
	ttl       time.Duration
}

// WorkingMemoryManager 工作记忆管理器
type WorkingMemoryManager struct {
	sessions    map[string]*WorkingMemory
	mu          sync.RWMutex
	maxSessions int
	defaultTTL  time.Duration
}

// NewWorkingMemory 创建新的工作记忆
func NewWorkingMemory(sessionID string, ctx context.Context) *WorkingMemory {
	now := time.Now()
	return &WorkingMemory{
		sessionID: sessionID,
		context:   ctx,
		data:      make(map[string]any),
		createdAt: now,
		updatedAt: now,
		ttl:       time.Hour, // 默认1小时TTL
	}
}

// NewWorkingMemorySimple 创建简单的工作记忆（无上下文）
func NewWorkingMemorySimple() *WorkingMemoryManager {
	return DefaultWorkingMemoryManager()
}

// NewWorkingMemoryManager 创建工作记忆管理器
func NewWorkingMemoryManager(maxSessions int, defaultTTL time.Duration) *WorkingMemoryManager {
	return &WorkingMemoryManager{
		sessions:    make(map[string]*WorkingMemory),
		maxSessions: maxSessions,
		defaultTTL:  defaultTTL,
	}
}

// GetOrCreateSession 获取或创建会话
func (wm *WorkingMemoryManager) GetOrCreateSession(sessionID string, ctx context.Context) *WorkingMemory {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// 清理过期会话
	wm.cleanupExpiredSessions()

	// 查找现有会话
	if session, exists := wm.sessions[sessionID]; exists {
		session.updatedAt = time.Now()
		return session
	}

	// 创建新会话
	session := NewWorkingMemory(sessionID, ctx)
	session.ttl = wm.defaultTTL

	// 如果达到最大会话数，清理最旧的会话
	if len(wm.sessions) >= wm.maxSessions {
		wm.removeOldestSession()
	}

	wm.sessions[sessionID] = session
	return session
}

// GetSession 获取会话
func (wm *WorkingMemoryManager) GetSession(sessionID string) (*WorkingMemory, bool) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	session, exists := wm.sessions[sessionID]
	if exists && !session.IsExpired() {
		return session, true
	}

	return nil, false
}

// DeleteSession 删除会话
func (wm *WorkingMemoryManager) DeleteSession(sessionID string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	delete(wm.sessions, sessionID)
}

// ListSessions 列出所有会话
func (wm *WorkingMemoryManager) ListSessions() []string {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	var sessionIDs []string
	for sessionID, session := range wm.sessions {
		if !session.IsExpired() {
			sessionIDs = append(sessionIDs, sessionID)
		}
	}

	return sessionIDs
}

// Set 设置工作记忆数据
func (wm *WorkingMemory) Set(key string, value any) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.data[key] = value
	wm.updatedAt = time.Now()
}

// Get 获取工作记忆数据
func (wm *WorkingMemory) Get(key string) (any, bool) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	value, exists := wm.data[key]
	return value, exists
}

// GetString 获取字符串类型数据
func (wm *WorkingMemory) GetString(key string) (string, bool) {
	value, exists := wm.Get(key)
	if !exists {
		return "", false
	}

	if str, ok := value.(string); ok {
		return str, true
	}

	return fmt.Sprintf("%v", value), true
}

// GetInt 获取整数类型数据
func (wm *WorkingMemory) GetInt(key string) (int, bool) {
	value, exists := wm.Get(key)
	if !exists {
		return 0, false
	}

	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

// GetBool 获取布尔类型数据
func (wm *WorkingMemory) GetBool(key string) (bool, bool) {
	value, exists := wm.Get(key)
	if !exists {
		return false, false
	}

	if b, ok := value.(bool); ok {
		return b, true
	}

	return false, false
}

// GetMap 获取map类型数据
func (wm *WorkingMemory) GetMap(key string) (map[string]any, bool) {
	value, exists := wm.Get(key)
	if !exists {
		return nil, false
	}

	if m, ok := value.(map[string]any); ok {
		return m, true
	}

	return nil, false
}

// GetSlice 获取切片类型数据
func (wm *WorkingMemory) GetSlice(key string) ([]any, bool) {
	value, exists := wm.Get(key)
	if !exists {
		return nil, false
	}

	if s, ok := value.([]any); ok {
		return s, true
	}

	return nil, false
}

// Delete 删除工作记忆数据
func (wm *WorkingMemory) Delete(key string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	delete(wm.data, key)
	wm.updatedAt = time.Now()
}

// Clear 清空工作记忆
func (wm *WorkingMemory) Clear() {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.data = make(map[string]any)
	wm.updatedAt = time.Now()
}

// GetAll 获取所有数据
func (wm *WorkingMemory) GetAll() map[string]any {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	// 返回副本
	result := make(map[string]any)
	for k, v := range wm.data {
		result[k] = v
	}

	return result
}

// Keys 获取所有键
func (wm *WorkingMemory) Keys() []string {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	var keys []string
	for k := range wm.data {
		keys = append(keys, k)
	}

	return keys
}

// HasKey 检查键是否存在
func (wm *WorkingMemory) HasKey(key string) bool {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	_, exists := wm.data[key]
	return exists
}

// Size 获取数据大小
func (wm *WorkingMemory) Size() int {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	return len(wm.data)
}

// IsEmpty 检查是否为空
func (wm *WorkingMemory) IsEmpty() bool {
	return wm.Size() == 0
}

// SetTTL 设置TTL
func (wm *WorkingMemory) SetTTL(ttl time.Duration) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.ttl = ttl
}

// GetTTL 获取TTL
func (wm *WorkingMemory) GetTTL() time.Duration {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	return wm.ttl
}

// IsExpired 检查是否过期
func (wm *WorkingMemory) IsExpired() bool {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	return time.Since(wm.updatedAt) > wm.ttl
}

// TimeToLive 获取剩余生存时间
func (wm *WorkingMemory) TimeToLive() time.Duration {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	elapsed := time.Since(wm.updatedAt)
	if elapsed > wm.ttl {
		return 0
	}

	return wm.ttl - elapsed
}

// GetSessionID 获取会话ID
func (wm *WorkingMemory) GetSessionID() string {
	return wm.sessionID
}

// GetContext 获取上下文
func (wm *WorkingMemory) GetContext() context.Context {
	return wm.context
}

// GetCreatedAt 获取创建时间
func (wm *WorkingMemory) GetCreatedAt() time.Time {
	return wm.createdAt
}

// GetUpdatedAt 获取更新时间
func (wm *WorkingMemory) GetUpdatedAt() time.Time {
	return wm.updatedAt
}

// UpdateContext 更新上下文
func (wm *WorkingMemory) UpdateContext(ctx context.Context) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.context = ctx
	wm.updatedAt = time.Now()
}

// Merge 合并其他工作记忆数据
func (wm *WorkingMemory) Merge(other *WorkingMemory, overwrite bool) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	other.mu.RLock()
	defer other.mu.RUnlock()

	for k, v := range other.data {
		if overwrite || !wm.HasKey(k) {
			wm.data[k] = v
		}
	}

	wm.updatedAt = time.Now()
}

// ToMap 转换为map
func (wm *WorkingMemory) ToMap() map[string]any {
	result := make(map[string]any)

	result["session_id"] = wm.sessionID
	result["created_at"] = wm.createdAt
	result["updated_at"] = wm.updatedAt
	result["ttl"] = wm.ttl.String()
	result["size"] = wm.Size()
	result["data"] = wm.GetAll()

	return result
}

// cleanupExpiredSessions 清理过期会话
func (wm *WorkingMemoryManager) cleanupExpiredSessions() {
	var expiredSessions []string

	for sessionID, session := range wm.sessions {
		if session.IsExpired() {
			expiredSessions = append(expiredSessions, sessionID)
		}
	}

	for _, sessionID := range expiredSessions {
		delete(wm.sessions, sessionID)
	}
}

// removeOldestSession 移除最旧的会话
func (wm *WorkingMemoryManager) removeOldestSession() {
	var oldestSessionID string
	var oldestTime time.Time

	for sessionID, session := range wm.sessions {
		if oldestSessionID == "" || session.updatedAt.Before(oldestTime) {
			oldestSessionID = sessionID
			oldestTime = session.updatedAt
		}
	}

	if oldestSessionID != "" {
		delete(wm.sessions, oldestSessionID)
	}
}

// Stats 获取统计信息
func (wm *WorkingMemoryManager) Stats() map[string]any {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	stats := make(map[string]any)
	stats["total_sessions"] = len(wm.sessions)
	stats["max_sessions"] = wm.maxSessions
	stats["default_ttl"] = wm.defaultTTL.String()

	// 计算活动会话数
	activeSessions := 0
	totalSize := 0
	for _, session := range wm.sessions {
		if !session.IsExpired() {
			activeSessions++
			totalSize += session.Size()
		}
	}

	stats["active_sessions"] = activeSessions
	stats["total_data_size"] = totalSize
	stats["avg_data_per_session"] = 0
	if activeSessions > 0 {
		stats["avg_data_per_session"] = totalSize / activeSessions
	}

	return stats
}

// DefaultWorkingMemoryManager 默认工作记忆管理器
func DefaultWorkingMemoryManager() *WorkingMemoryManager {
	return NewWorkingMemoryManager(100, time.Hour)
}
