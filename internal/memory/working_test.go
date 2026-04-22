package memory

import (
	"context"
	"testing"
	"time"
)

// TestWorkingMemoryBasic 测试工作记忆基础功能
func TestWorkingMemoryBasic(t *testing.T) {
	ctx := context.Background()
	sessionID := "test-session"
	wm := NewWorkingMemory(sessionID, ctx)

	// 测试1: 设置和获取字符串
	t.Run("设置和获取字符串", func(t *testing.T) {
		wm.Set("name", "Alice")

		value, exists := wm.Get("name")
		if !exists {
			t.Error("键应该存在")
		}
		if value != "Alice" {
			t.Errorf("值不正确: 期望 %s, 实际 %v", "Alice", value)
		}

		// 测试 GetString
		strValue, ok := wm.GetString("name")
		if !ok {
			t.Error("获取字符串应该成功")
		}
		if strValue != "Alice" {
			t.Errorf("字符串值不正确: 期望 %s, 实际 %s", "Alice", strValue)
		}
	})

	// 测试2: 设置和获取整数
	t.Run("设置和获取整数", func(t *testing.T) {
		wm.Set("count", 42)

		intValue, ok := wm.GetInt("count")
		if !ok {
			t.Error("获取整数应该成功")
		}
		if intValue != 42 {
			t.Errorf("整数值不正确: 期望 %d, 实际 %d", 42, intValue)
		}
	})

	// 测试3: 设置和获取布尔值
	t.Run("设置和获取布尔值", func(t *testing.T) {
		wm.Set("active", true)

		boolValue, ok := wm.GetBool("active")
		if !ok {
			t.Error("获取布尔值应该成功")
		}
		if !boolValue {
			t.Errorf("布尔值不正确: 期望 %t, 实际 %t", true, boolValue)
		}
	})

	// 测试4: 设置和获取Map
	t.Run("设置和获取Map", func(t *testing.T) {
		testMap := map[string]any{
			"key1": "value1",
			"key2": 123,
		}
		wm.Set("config", testMap)

		mapValue, ok := wm.GetMap("config")
		if !ok {
			t.Error("获取Map应该成功")
		}
		if mapValue["key1"] != "value1" {
			t.Errorf("Map值不正确: 期望 %s, 实际 %v", "value1", mapValue["key1"])
		}
	})

	// 测试5: 检查存在性
	t.Run("检查存在性", func(t *testing.T) {
		key := "exists-test"
		value := "test"

		if wm.HasKey(key) {
			t.Error("键不应该存在")
		}

		wm.Set(key, value)
		if !wm.HasKey(key) {
			t.Error("键应该存在")
		}

		wm.Delete(key)
		if wm.HasKey(key) {
			t.Error("删除后键不应该存在")
		}
	})

	// 测试6: 获取所有键
	t.Run("获取所有键", func(t *testing.T) {
		keys := wm.Keys()
		// 应该有之前设置的值
		if len(keys) < 3 {
			t.Errorf("键数量不正确: 期望至少3个, 实际 %d", len(keys))
		}
	})

	// 测试7: 获取所有数据
	t.Run("获取所有数据", func(t *testing.T) {
		allData := wm.GetAll()
		if len(allData) == 0 {
			t.Error("应该有数据")
		}
	})

	// 测试8: 清空数据
	t.Run("清空数据", func(t *testing.T) {
		wm.Clear()

		size := wm.Size()
		if size != 0 {
			t.Errorf("清空后大小应该为0，实际 %d", size)
		}
		if !wm.IsEmpty() {
			t.Error("清空后应该为空")
		}
	})

	// 测试9: 会话信息
	t.Run("会话信息", func(t *testing.T) {
		// 测试会话ID
		id := wm.GetSessionID()
		if id != sessionID {
			t.Errorf("会话ID不正确: 期望 %s, 实际 %s", sessionID, id)
		}

		// 测试上下文
		ctx := wm.GetContext()
		if ctx == nil {
			t.Error("上下文不应该为nil")
		}

		// 测试创建时间
		createdAt := wm.GetCreatedAt()
		if createdAt.IsZero() {
			t.Error("创建时间不应该为零")
		}

		// 测试更新时间
		updatedAt := wm.GetUpdatedAt()
		if updatedAt.IsZero() {
			t.Error("更新时间不应该为零")
		}
	})

	// 测试10: TTL功能
	t.Run("TTL功能", func(t *testing.T) {
		wm.SetTTL(5 * time.Second)
		ttl := wm.GetTTL()
		if ttl != 5*time.Second {
			t.Errorf("TTL不正确: 期望 %v, 实际 %v", 5*time.Second, ttl)
		}

		if wm.IsExpired() {
			t.Error("刚设置的TTL不应该过期")
		}

		remaining := wm.TimeToLive()
		if remaining <= 0 || remaining > 5*time.Second {
			t.Errorf("剩余时间不合理: %v", remaining)
		}
	})

	// 测试11: 更新上下文
	t.Run("更新上下文", func(t *testing.T) {
		newCtx := context.WithValue(context.Background(), "key", "value")
		wm.UpdateContext(newCtx)

		ctx := wm.GetContext()
		if ctx == nil {
			t.Error("上下文不应该为nil")
		}
	})

	// 测试12: 转换为Map
	t.Run("转换为Map", func(t *testing.T) {
		// 先设置一些数据
		wm.Set("key1", "value1")
		wm.Set("key2", 42)

		data := wm.ToMap()
		if data == nil {
			t.Error("ToMap不应该返回nil")
		}
	})
}

// TestWorkingMemoryMerge 测试合并功能
func TestWorkingMemoryMerge(t *testing.T) {
	ctx := context.Background()

	wm1 := NewWorkingMemory("session-1", ctx)
	wm1.Set("shared-key", "value-from-1")
	wm1.Set("unique-key-1", "only-in-1")

	wm2 := NewWorkingMemory("session-2", ctx)
	wm2.Set("shared-key", "value-from-2")
	wm2.Set("unique-key-2", "only-in-2")

	// 测试合并到wm1，不覆盖已有值
	t.Run("合并不覆盖", func(t *testing.T) {
		wm1.Merge(wm2, false)

		// 共享键应该保留原值（不覆盖）
		value, exists := wm1.Get("shared-key")
		if !exists {
			t.Error("shared-key 应该存在")
		}
		if value != "value-from-1" {
			t.Errorf("合并后应该保留原值: 期望 %s, 实际 %v", "value-from-1", value)
		}

		// 新键应该被添加
		value, exists = wm1.Get("unique-key-2")
		if !exists {
			t.Error("unique-key-2 应该被合并")
		}
		if value != "only-in-2" {
			t.Errorf("新键值不正确: 期望 %s, 实际 %v", "only-in-2", value)
		}
	})

	// 测试合并到新工作记忆，覆盖已有值
	t.Run("合并覆盖", func(t *testing.T) {
		wm3 := NewWorkingMemory("session-3", ctx)
		wm3.Set("shared-key", "original")

		wm4 := NewWorkingMemory("session-4", ctx)
		wm4.Set("shared-key", "new-value")

		wm3.Merge(wm4, true)

		// 共享键应该被覆盖
		value, exists := wm3.Get("shared-key")
		if !exists {
			t.Error("shared-key 应该存在")
		}
		if value != "new-value" {
			t.Errorf("合并后应该被覆盖: 期望 %s, 实际 %v", "new-value", value)
		}
	})
}

// TestWorkingMemoryManager 测试工作记忆管理器
func TestWorkingMemoryManager(t *testing.T) {
	// 创建管理器
	manager := NewWorkingMemoryManager(5, 30*time.Minute)

	// 测试1: 获取或创建会话
	t.Run("获取或创建会话", func(t *testing.T) {
		sessionID := "test-session-manager"
		ctx := context.Background()

		wm := manager.GetOrCreateSession(sessionID, ctx)
		if wm == nil {
			t.Error("工作记忆不应该为nil")
		}
		if wm.GetSessionID() != sessionID {
			t.Errorf("会话ID不正确: 期望 %s, 实际 %s", sessionID, wm.GetSessionID())
		}

		// 验证会话已添加
		session, exists := manager.GetSession(sessionID)
		if !exists {
			t.Error("会话应该存在")
		}
		if session != wm {
			t.Error("获取的会话应该与创建的相同")
		}
	})

	// 测试2: 重复获取或创建会话
	t.Run("重复获取或创建会话", func(t *testing.T) {
		sessionID := "duplicate-session"
		ctx := context.Background()

		// 第一次创建
		wm1 := manager.GetOrCreateSession(sessionID, ctx)
		if wm1 == nil {
			t.Errorf("第一次创建会话失败")
		}

		// 第二次获取或创建（应该返回相同的实例）
		wm2 := manager.GetOrCreateSession(sessionID, ctx)
		if wm2 == nil {
			t.Errorf("第二次获取或创建会话失败")
		}
		if wm1 != wm2 {
			t.Error("重复获取或创建应该返回相同的实例")
		}
	})

	// 测试3: 获取不存在的会话
	t.Run("获取不存在的会话", func(t *testing.T) {
		_, exists := manager.GetSession("non-existent-session")
		if exists {
			t.Error("不存在的会话应该返回false")
		}
	})

	// 测试4: 删除会话
	t.Run("删除会话", func(t *testing.T) {
		sessionID := "delete-session"
		ctx := context.Background()

		// 创建会话
		wm := manager.GetOrCreateSession(sessionID, ctx)
		if wm == nil {
			t.Errorf("创建会话失败")
		}

		// 验证会话存在
		_, exists := manager.GetSession(sessionID)
		if !exists {
			t.Error("会话应该存在")
		}

		// 删除会话
		manager.DeleteSession(sessionID)

		// 验证会话已删除
		_, exists = manager.GetSession(sessionID)
		if exists {
			t.Error("删除后会话不应该存在")
		}
	})

	// 测试5: 列出会话
	t.Run("列出会话", func(t *testing.T) {
		sessions := manager.ListSessions()
		// 应该有之前创建的会话
		if len(sessions) < 2 {
			t.Errorf("会话数量不正确: 期望至少2个, 实际 %d", len(sessions))
		}
	})

	// 测试6: 会话限制
	t.Run("会话限制", func(t *testing.T) {
		limitedManager := NewWorkingMemoryManager(3, 30*time.Minute)
		ctx := context.Background()

		// 创建3个会话
		for i := 0; i < 3; i++ {
			sessionID := "limited-session-" + string(rune('A'+i))
			wm := limitedManager.GetOrCreateSession(sessionID, ctx)
			if wm == nil {
				t.Errorf("创建会话 %s 失败", sessionID)
			}
		}

		// 验证会话数量为3
		if len(limitedManager.ListSessions()) != 3 {
			t.Errorf("会话数量不正确: 期望 3, 实际 %d", len(limitedManager.ListSessions()))
		}
	})

	// 测试7: 管理器统计
	t.Run("管理器统计", func(t *testing.T) {
		stats := manager.Stats()
		if stats == nil {
			t.Error("统计信息不应该为nil")
		}
		if totalSessions, ok := stats["total_sessions"].(int); ok {
			if totalSessions < 2 {
				t.Errorf("总会话数不正确: 期望至少2, 实际 %d", totalSessions)
			}
		}
	})
}

// TestWorkingMemoryConcurrent 测试并发访问
func TestWorkingMemoryConcurrent(t *testing.T) {
	ctx := context.Background()
	wm := NewWorkingMemory("concurrent-session", ctx)

	// 并发写入
	numWorkers := 10
	done := make(chan bool, numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func(index int) {
			key := "key-" + string(rune('A'+index))
			wm.Set(key, index)
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < numWorkers; i++ {
		<-done
	}

	// 验证所有键都已写入
	keys := wm.Keys()
	if len(keys) != numWorkers {
		t.Errorf("并发写入后键数量不正确: 期望 %d, 实际 %d", numWorkers, len(keys))
	}

	// 并发读取
	for i := 0; i < numWorkers; i++ {
		go func(index int) {
			key := "key-" + string(rune('A'+index))
			_, exists := wm.Get(key)
			if !exists {
				t.Errorf("键 %s 应该存在", key)
			}
			done <- true
		}(i)
	}

	// 等待所有读取完成
	for i := 0; i < numWorkers; i++ {
		<-done
	}
}

// TestWorkingMemoryEdgeCases 测试边界情况
func TestWorkingMemoryEdgeCases(t *testing.T) {
	ctx := context.Background()
	wm := NewWorkingMemory("edge-case-session", ctx)

	// 测试1: 获取不存在的键
	t.Run("获取不存在的键", func(t *testing.T) {
		_, exists := wm.Get("non-existent")
		if exists {
			t.Error("不存在的键应该返回false")
		}
	})

	// 测试2: 删除不存在的键
	t.Run("删除不存在的键", func(t *testing.T) {
		// 不应该panic
		wm.Delete("non-existent")
	})

	// 测试3: 空工作记忆
	t.Run("空工作记忆", func(t *testing.T) {
		wm.Clear()
		if !wm.IsEmpty() {
			t.Error("清空后应该为空")
		}
		if wm.Size() != 0 {
			t.Errorf("清空后大小应该为0，实际 %d", wm.Size())
		}
	})

	// 测试4: 类型安全的获取方法
	t.Run("类型安全的获取方法", func(t *testing.T) {
		wm.Set("string-val", "test")
		wm.Set("int-val", 42)
		wm.Set("bool-val", true)
		wm.Set("map-val", map[string]any{"a": 1})
		wm.Set("slice-val", []any{1, 2, 3})

		// 错误类型获取（GetString总是返回true，因为它使用Sprintf做后备）
		_, ok := wm.GetInt("string-val")
		if ok {
			t.Error("从字符串键获取整数应该失败")
		}

		_, ok = wm.GetBool("string-val")
		if ok {
			t.Error("从字符串键获取布尔值应该失败")
		}

		_, ok = wm.GetMap("string-val")
		if ok {
			t.Error("从字符串键获取Map应该失败")
		}

		_, ok = wm.GetSlice("string-val")
		if ok {
			t.Error("从字符串键获取Slice应该失败")
		}
	})

	// 测试5: 覆盖已存在的值
	t.Run("覆盖已存在的值", func(t *testing.T) {
		wm.Set("override-key", "original")
		wm.Set("override-key", "updated")

		value, exists := wm.Get("override-key")
		if !exists {
			t.Error("键应该存在")
		}
		if value != "updated" {
			t.Errorf("值应该被覆盖: 期望 %s, 实际 %v", "updated", value)
		}
	})
}

// TestDefaultWorkingMemoryManager 测试默认管理器
func TestDefaultWorkingMemoryManager(t *testing.T) {
	manager := DefaultWorkingMemoryManager()
	if manager == nil {
		t.Error("默认管理器不应该为nil")
	}

	ctx := context.Background()
	wm := manager.GetOrCreateSession("default-session", ctx)
	if wm == nil {
		t.Error("获取或创建会话不应该为nil")
	}

	// 验证可以正常使用
	wm.Set("test", "value")
	value, exists := wm.Get("test")
	if !exists || value != "value" {
		t.Errorf("默认管理器使用失败: 期望 %s, 实际 %v", "value", value)
	}
}
