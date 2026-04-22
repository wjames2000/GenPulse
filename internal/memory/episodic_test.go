package memory

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestEpisodicMemoryBasic 测试情节记忆基础功能
func TestEpisodicMemoryBasic(t *testing.T) {
	// 创建临时数据库文件
	tempDir, err := os.MkdirTemp("", "episodic-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")

	// 创建情节记忆实例
	em, err := NewEpisodicMemory(dbPath)
	if err != nil {
		t.Fatalf("创建情节记忆失败: %v", err)
	}
	defer em.Close()

	// 测试1: 数据库初始化
	t.Run("数据库初始化", func(t *testing.T) {
		// 验证主表已创建
		var count int
		query := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='memories'"
		err := em.db.QueryRow(query).Scan(&count)
		if err != nil {
			t.Errorf("查询表 memories 失败: %v", err)
		}
		if count != 1 {
			t.Error("表 memories 应该存在")
		}

		// 验证索引已创建
		indexes := []string{"idx_memories_session_id", "idx_memories_task_id", "idx_memories_created_at"}
		for _, index := range indexes {
			query := "SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name=?"
			err := em.db.QueryRow(query, index).Scan(&count)
			if err != nil {
				t.Errorf("查询索引 %s 失败: %v", index, err)
			}
			if count != 1 {
				t.Errorf("索引 %s 应该存在", index)
			}
		}
	})

	// 测试2: 存储记忆记录
	t.Run("存储记忆记录", func(t *testing.T) {
		record := &MemoryRecord{
			ID:          "test-record-1",
			SessionID:   "test-session-1",
			TaskID:      "test-task-1",
			TaskType:    "code_generation",
			Description: "测试记忆记录",
			Content:     "这是一个测试记忆记录的内容",
			Metadata: map[string]any{
				"language":   "go",
				"complexity": "simple",
			},
			Tags:        []string{"test", "go", "memory"},
			Category:    "development",
			Importance:  0.7,
			Success:     true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			AccessedAt:  time.Now(),
			AccessCount: 0,
		}

		err := em.Store(record)
		if err != nil {
			t.Errorf("存储记忆记录失败: %v", err)
		}

		// 验证记录已存储
		retrieved, err := em.Get("test-record-1")
		if err != nil {
			t.Errorf("获取记忆记录失败: %v", err)
		}
		if retrieved == nil {
			t.Error("获取的记录不应该为nil")
		}
		if retrieved.ID != "test-record-1" {
			t.Errorf("记录ID不正确: 期望 %s, 实际 %s", "test-record-1", retrieved.ID)
		}
		if retrieved.Description != "测试记忆记录" {
			t.Errorf("记录描述不正确: 期望 %s, 实际 %s", "测试记忆记录", retrieved.Description)
		}
		if retrieved.Importance != 0.7 {
			t.Errorf("记录重要性不正确: 期望 %f, 实际 %f", 0.7, retrieved.Importance)
		}
		if !retrieved.Success {
			t.Error("记录应该标记为成功")
		}
		if len(retrieved.Tags) != 3 {
			t.Errorf("记录标签数量不正确: 期望 %d, 实际 %d", 3, len(retrieved.Tags))
		}
	})

	// 测试3: 存储多个记录
	t.Run("存储多个记录", func(t *testing.T) {
		records := []*MemoryRecord{
			{
				ID:          "test-record-2",
				SessionID:   "test-session-1",
				TaskID:      "test-task-2",
				TaskType:    "code_review",
				Description: "代码审查记录",
				Content:     "审查Go代码的质量",
				Metadata:    map[string]any{"lines": 50},
				Tags:        []string{"review", "go"},
				Category:    "quality",
				Importance:  0.8,
				Success:     true,
				CreatedAt:   time.Now(),
			},
			{
				ID:          "test-record-3",
				SessionID:   "test-session-2",
				TaskID:      "test-task-3",
				TaskType:    "error_fix",
				Description: "错误修复记录",
				Content:     "修复空指针异常",
				Metadata:    map[string]any{"error_type": "nil_pointer"},
				Tags:        []string{"bug", "fix"},
				Category:    "maintenance",
				Importance:  0.9,
				Success:     true,
				CreatedAt:   time.Now(),
			},
			{
				ID:           "test-record-4",
				SessionID:    "test-session-2",
				TaskID:       "test-task-4",
				TaskType:     "testing",
				Description:  "测试失败记录",
				Content:      "单元测试失败",
				Metadata:     map[string]any{"test_count": 10, "failed": 2},
				Tags:         []string{"test", "failure"},
				Category:     "testing",
				Importance:   0.6,
				Success:      false,
				ErrorType:    "test_failure",
				ErrorMessage: "断言失败: 期望 5, 实际 3",
				CreatedAt:    time.Now(),
			},
		}

		for _, record := range records {
			err := em.Store(record)
			if err != nil {
				t.Errorf("存储记录 %s 失败: %v", record.ID, err)
			}
		}

		// 验证记录数量
		stats, err := em.GetStats()
		if err != nil {
			t.Errorf("获取统计信息失败: %v", err)
		}
		if totalCount, ok := stats["total_count"].(int); !ok || totalCount != 4 { // 包括之前存储的test-record-1
			t.Errorf("记录数量不正确: 期望 %d, 实际 %v", 4, stats["total_count"])
		}
	})

	// 测试4: 获取不存在的记录
	t.Run("获取不存在的记录", func(t *testing.T) {
		record, err := em.Get("non-existent-record")
		if err == nil {
			t.Error("获取不存在的记录应该返回错误")
		}
		if record != nil {
			t.Error("不存在的记录应该返回nil")
		}
	})

	// 测试5: 更新记录
	t.Run("更新记录", func(t *testing.T) {
		// 获取现有记录
		record, err := em.Get("test-record-1")
		if err != nil {
			t.Fatalf("获取记录失败: %v", err)
		}
		if record == nil {
			t.Fatal("记录不应该为nil")
		}

		// 修改记录
		record.Description = "更新后的描述"
		record.Importance = 0.9
		record.UpdatedAt = time.Now()

		err = em.Update(record)
		if err != nil {
			t.Errorf("更新记录失败: %v", err)
		}

		// 验证更新
		updated, err := em.Get("test-record-1")
		if err != nil {
			t.Errorf("获取更新后的记录失败: %v", err)
		}
		if updated.Description != "更新后的描述" {
			t.Errorf("更新后的描述不正确: 期望 %s, 实际 %s", "更新后的描述", updated.Description)
		}
		if updated.Importance != 0.9 {
			t.Errorf("更新后的重要性不正确: 期望 %f, 实际 %f", 0.9, updated.Importance)
		}
	})

	// 测试6: 删除记录
	t.Run("删除记录", func(t *testing.T) {
		// 删除记录
		err := em.Delete("test-record-2")
		if err != nil {
			t.Errorf("删除记录失败: %v", err)
		}

		// 验证记录已删除
		record, err := em.Get("test-record-2")
		if err == nil {
			t.Error("获取已删除的记录应该返回错误")
		}
		if record != nil {
			t.Error("已删除的记录应该返回nil")
		}

		// 删除不存在的记录（可能会返回错误，也可能不会，取决于实现）
		err = em.Delete("non-existent-record")
		_ = err // 忽略删除不存在的记录的返回值
	})

	// 测试7: 按会话获取记录
	t.Run("按会话获取记录", func(t *testing.T) {
		// 获取test-session-1的所有记录
		records, err := em.GetBySession("test-session-1")
		if err != nil {
			t.Errorf("按会话获取记录失败: %v", err)
		}

		// test-session-1应该有1个记录（test-record-1，test-record-2已被删除）
		if len(records) != 1 {
			t.Errorf("会话记录数量不正确: 期望 %d, 实际 %d", 1, len(records))
		}
		if records[0].SessionID != "test-session-1" {
			t.Errorf("记录会话ID不正确: 期望 %s, 实际 %s", "test-session-1", records[0].SessionID)
		}

		// 获取test-session-2的所有记录
		records, err = em.GetBySession("test-session-2")
		if err != nil {
			t.Errorf("按会话获取记录失败: %v", err)
		}

		// test-session-2应该有2个记录
		if len(records) != 2 {
			t.Errorf("会话记录数量不正确: 期望 %d, 实际 %d", 2, len(records))
		}
	})

	// 测试8: 按任务类型获取记录
	t.Run("按任务类型获取记录", func(t *testing.T) {
		records, err := em.GetByTaskType("testing", 10)
		if err != nil {
			t.Errorf("按任务类型获取记录失败: %v", err)
		}

		if len(records) != 1 {
			t.Errorf("任务类型记录数量不正确: 期望 %d, 实际 %d", 1, len(records))
		}

		if records[0].ID != "test-record-4" {
			t.Errorf("任务类型记录ID不正确: 期望 %s, 实际 %s", "test-record-4", records[0].ID)
		}
	})

	// 测试9: 获取最近记录
	t.Run("获取最近记录", func(t *testing.T) {
		records, err := em.GetRecent(2)
		if err != nil {
			t.Errorf("获取最近记录失败: %v", err)
		}

		if len(records) != 2 {
			t.Errorf("最近记录数量不正确: 期望 %d, 实际 %d", 2, len(records))
		}

		// 验证记录按创建时间倒序排列
		if records[0].CreatedAt.Before(records[1].CreatedAt) {
			t.Error("记录应该按创建时间倒序排列")
		}
	})

	// 测试10: 获取成功记录
	t.Run("获取成功记录", func(t *testing.T) {
		records, err := em.GetSuccessful(10)
		if err != nil {
			t.Errorf("获取成功记录失败: %v", err)
		}

		// 应该有2个成功记录（test-record-1, test-record-3）
		if len(records) != 2 {
			t.Errorf("成功记录数量不正确: 期望 %d, 实际 %d", 2, len(records))
		}

		for _, record := range records {
			if !record.Success {
				t.Errorf("记录 %s 应该标记为成功", record.ID)
			}
		}
	})

	// 测试11: 获取失败记录
	t.Run("获取失败记录", func(t *testing.T) {
		records, err := em.GetFailed(10)
		if err != nil {
			t.Errorf("获取失败记录失败: %v", err)
		}

		// 应该有1个失败记录（test-record-4）
		if len(records) != 1 {
			t.Errorf("失败记录数量不正确: 期望 %d, 实际 %d", 1, len(records))
		}

		for _, record := range records {
			if record.Success {
				t.Errorf("记录 %s 应该标记为失败", record.ID)
			}
		}
	})
}
