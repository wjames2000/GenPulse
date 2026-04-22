package memory

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestEpisodicMemoryStatistics 测试统计功能
func TestEpisodicMemoryStatistics(t *testing.T) {
	// 创建临时数据库
	tempDir, err := os.MkdirTemp("", "episodic-stats-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "stats-test.db")
	em, err := NewEpisodicMemory(dbPath)
	if err != nil {
		t.Fatalf("创建情节记忆失败: %v", err)
	}
	defer em.Close()

	// 添加测试数据
	records := []*MemoryRecord{
		{
			ID:          "stats-1",
			SessionID:   "session-1",
			TaskType:    "code_generation",
			Description: "记录1",
			Content:     "内容1",
			Category:    "backend",
			Importance:  0.9,
			Success:     true,
			CreatedAt:   time.Now(),
		},
		{
			ID:          "stats-2",
			SessionID:   "session-1",
			TaskType:    "code_review",
			Description: "记录2",
			Content:     "内容2",
			Category:    "backend",
			Importance:  0.8,
			Success:     true,
			CreatedAt:   time.Now(),
		},
		{
			ID:          "stats-3",
			SessionID:   "session-2",
			TaskType:    "bug_fix",
			Description: "记录3",
			Content:     "内容3",
			Category:    "frontend",
			Importance:  0.7,
			Success:     false,
			ErrorType:   "runtime_error",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "stats-4",
			SessionID:   "session-2",
			TaskType:    "documentation",
			Description: "记录4",
			Content:     "内容4",
			Category:    "docs",
			Importance:  0.6,
			Success:     true,
			CreatedAt:   time.Now().Add(-48 * time.Hour),
		},
	}

	for _, record := range records {
		if err := em.Store(record); err != nil {
			t.Errorf("存储记录 %s 失败: %v", record.ID, err)
		}
	}

	// 测试1: 获取基本统计
	t.Run("基本统计", func(t *testing.T) {
		stats, err := em.GetStats()
		if err != nil {
			t.Errorf("获取统计信息失败: %v", err)
		}

		// 验证统计信息
		if totalCount, ok := stats["total_count"].(int); !ok || totalCount != 4 {
			t.Errorf("总记录数不正确: 期望 4, 实际 %v", stats["total_count"])
		}

		if successCount, ok := stats["success_count"].(int); !ok || successCount != 3 {
			t.Errorf("成功记录数不正确: 期望 3, 实际 %v", stats["success_count"])
		}

		if failedCount, ok := stats["failed_count"].(int); !ok || failedCount != 1 {
			t.Errorf("失败记录数不正确: 期望 1, 实际 %v", stats["failed_count"])
		}

		if successRate, ok := stats["success_rate"].(float64); !ok || successRate != 0.75 {
			t.Errorf("成功率不正确: 期望 0.75, 实际 %v", stats["success_rate"])
		}

		if avgImportance, ok := stats["avg_importance"].(float64); !ok || avgImportance < 0.6 || avgImportance > 0.9 {
			t.Errorf("平均重要性不正确: 期望 0.6-0.9, 实际 %v", stats["avg_importance"])
		}
	})

	// 测试2: 按会话获取记录
	t.Run("按会话获取记录", func(t *testing.T) {
		session1Records, err := em.GetBySession("session-1")
		if err != nil {
			t.Errorf("按会话获取记录失败: %v", err)
		}

		if len(session1Records) != 2 {
			t.Errorf("session-1 记录数不正确: 期望 2, 实际 %d", len(session1Records))
		}

		session2Records, err := em.GetBySession("session-2")
		if err != nil {
			t.Errorf("按会话获取记录失败: %v", err)
		}

		if len(session2Records) != 2 {
			t.Errorf("session-2 记录数不正确: 期望 2, 实际 %d", len(session2Records))
		}
	})

	// 测试3: 按任务类型获取记录
	t.Run("按任务类型获取记录", func(t *testing.T) {
		codeGenRecords, err := em.GetByTaskType("code_generation", 10)
		if err != nil {
			t.Errorf("按任务类型获取记录失败: %v", err)
		}

		if len(codeGenRecords) != 1 {
			t.Errorf("code_generation 记录数不正确: 期望 1, 实际 %d", len(codeGenRecords))
		}

		codeReviewRecords, err := em.GetByTaskType("code_review", 10)
		if err != nil {
			t.Errorf("按任务类型获取记录失败: %v", err)
		}

		if len(codeReviewRecords) != 1 {
			t.Errorf("code_review 记录数不正确: 期望 1, 实际 %d", len(codeReviewRecords))
		}
	})

	// 测试4: 获取最近记录
	t.Run("获取最近记录", func(t *testing.T) {
		recentRecords, err := em.GetRecent(2)
		if err != nil {
			t.Errorf("获取最近记录失败: %v", err)
		}

		if len(recentRecords) != 2 {
			t.Errorf("最近记录数不正确: 期望 2, 实际 %d", len(recentRecords))
		}

		// 验证记录是按创建时间降序排列的
		if len(recentRecords) >= 2 {
			if recentRecords[0].CreatedAt.Before(recentRecords[1].CreatedAt) {
				t.Error("最近记录应该按创建时间降序排列")
			}
		}
	})

	// 测试5: 获取最常访问记录
	t.Run("获取最常访问记录", func(t *testing.T) {
		// 先访问一些记录以增加访问计数
		for i := 0; i < 3; i++ {
			if _, err := em.Get("stats-1"); err != nil {
				t.Errorf("获取记录失败: %v", err)
			}
		}
		for i := 0; i < 2; i++ {
			if _, err := em.Get("stats-2"); err != nil {
				t.Errorf("获取记录失败: %v", err)
			}
		}

		mostAccessed, err := em.GetMostAccessed(2)
		if err != nil {
			t.Errorf("获取最常访问记录失败: %v", err)
		}

		if len(mostAccessed) >= 1 {
			if mostAccessed[0].ID != "stats-1" {
				t.Errorf("最常访问的记录应该是 stats-1, 实际是 %s", mostAccessed[0].ID)
			}
		}
	})

	// 测试6: 按重要性获取记录
	t.Run("按重要性获取记录", func(t *testing.T) {
		importantRecords, err := em.GetByImportance(0.8, 10)
		if err != nil {
			t.Errorf("按重要性获取记录失败: %v", err)
		}

		if len(importantRecords) != 2 {
			t.Errorf("重要性>=0.8的记录数不正确: 期望 2, 实际 %d", len(importantRecords))
		}

		for _, record := range importantRecords {
			if record.Importance < 0.8 {
				t.Errorf("记录 %s 的重要性 %f 低于最小值 0.8", record.ID, record.Importance)
			}
		}
	})

	// 测试7: 获取成功记录
	t.Run("获取成功记录", func(t *testing.T) {
		successfulRecords, err := em.GetSuccessful(10)
		if err != nil {
			t.Errorf("获取成功记录失败: %v", err)
		}

		if len(successfulRecords) != 3 {
			t.Errorf("成功记录数不正确: 期望 3, 实际 %d", len(successfulRecords))
		}

		for _, record := range successfulRecords {
			if !record.Success {
				t.Errorf("记录 %s 应该标记为成功", record.ID)
			}
		}
	})

	// 测试8: 获取失败记录
	t.Run("获取失败记录", func(t *testing.T) {
		failedRecords, err := em.GetFailed(10)
		if err != nil {
			t.Errorf("获取失败记录失败: %v", err)
		}

		if len(failedRecords) != 1 {
			t.Errorf("失败记录数不正确: 期望 1, 实际 %d", len(failedRecords))
		}

		for _, record := range failedRecords {
			if record.Success {
				t.Errorf("记录 %s 应该标记为失败", record.ID)
			}
		}
	})

	// 测试9: 清理旧记录
	t.Run("清理旧记录", func(t *testing.T) {
		// 清理48小时前且重要性<0.7的记录
		deletedCount, err := em.Cleanup(2, 0.7)
		if err != nil {
			t.Errorf("清理记录失败: %v", err)
		}

		// 应该删除 stats-4 (48小时前，重要性0.6)
		if deletedCount != 1 {
			t.Errorf("清理记录数不正确: 期望 1, 实际 %d", deletedCount)
		}

		// 验证记录已被删除
		_, err = em.Get("stats-4")
		if err == nil {
			t.Error("记录 stats-4 应该已被删除")
		}

		// 验证其他记录仍然存在
		_, err = em.Get("stats-1")
		if err != nil {
			t.Error("记录 stats-1 应该仍然存在")
		}
	})
}
