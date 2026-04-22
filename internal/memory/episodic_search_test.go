package memory

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestEpisodicMemorySearch 测试搜索功能
func TestEpisodicMemorySearch(t *testing.T) {
	// 创建临时数据库
	tempDir, err := os.MkdirTemp("", "episodic-search-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "search-test.db")
	em, err := NewEpisodicMemory(dbPath)
	if err != nil {
		t.Fatalf("创建情节记忆失败: %v", err)
	}
	defer em.Close()

	// 添加测试数据
	records := []*MemoryRecord{
		{
			ID:          "search-1",
			SessionID:   "session-a",
			TaskID:      "task-1",
			TaskType:    "code_generation",
			Description: "生成用户认证模块",
			Content:     "实现JWT认证和用户注册登录功能",
			Metadata:    map[string]any{"language": "go", "framework": "gin"},
			Tags:        []string{"auth", "security", "go", "api"},
			Category:    "backend",
			Importance:  0.9,
			Success:     true,
			CreatedAt:   time.Now(),
		},
		{
			ID:          "search-2",
			SessionID:   "session-a",
			TaskID:      "task-2",
			TaskType:    "database",
			Description: "设计用户数据库表",
			Content:     "创建users表和相关的索引",
			Metadata:    map[string]any{"database": "postgresql", "tables": 1},
			Tags:        []string{"database", "schema", "sql"},
			Category:    "database",
			Importance:  0.8,
			Success:     true,
			CreatedAt:   time.Now().Add(-time.Hour),
		},
		{
			ID:          "search-3",
			SessionID:   "session-b",
			TaskID:      "task-3",
			TaskType:    "frontend",
			Description: "创建登录页面",
			Content:     "使用React实现登录表单和验证",
			Metadata:    map[string]any{"framework": "react", "components": 3},
			Tags:        []string{"frontend", "react", "ui", "form"},
			Category:    "frontend",
			Importance:  0.7,
			Success:     true,
			CreatedAt:   time.Now().Add(-2 * time.Hour),
		},
		{
			ID:           "search-4",
			SessionID:    "session-b",
			TaskID:       "task-4",
			TaskType:     "error",
			Description:  "修复数据库连接错误",
			Content:      "解决PostgreSQL连接池耗尽问题",
			Metadata:     map[string]any{"error": "connection_pool", "fixed": true},
			Tags:         []string{"bug", "database", "fix"},
			Category:     "maintenance",
			Importance:   0.6,
			Success:      false,
			ErrorType:    "database_error",
			ErrorMessage: "连接池大小配置错误",
			CreatedAt:    time.Now().Add(-3 * time.Hour),
		},
	}

	for _, record := range records {
		err := em.Store(record)
		if err != nil {
			t.Errorf("存储记录 %s 失败: %v", record.ID, err)
		}
	}

	// 测试1: 全文搜索
	t.Run("全文搜索", func(t *testing.T) {
		query := &SearchQuery{
			Query: "auth",
			Limit: 10,
		}

		results, err := em.Search(query)
		if err != nil {
			t.Errorf("全文搜索失败: %v", err)
		}

		// 应该至少找到1个记录（search-1包含tags中的auth）
		if len(results) < 1 {
			t.Errorf("搜索结果数量不正确: 期望至少1个, 实际 %d", len(results))
		}

		// 验证结果包含相关记录
		foundSearch1 := false
		for _, result := range results {
			if result.Record.ID == "search-1" {
				foundSearch1 = true
			}
		}
		if !foundSearch1 {
			t.Error("搜索结果应该包含search-1")
		}
	})

	// 测试2: 按标签搜索
	t.Run("按标签搜索", func(t *testing.T) {
		query := &SearchQuery{
			Query: "database",
			Limit: 10,
		}

		results, err := em.Search(query)
		if err != nil {
			t.Errorf("按标签搜索失败: %v", err)
		}

		// 应该找到2个记录（search-2和search-4）
		if len(results) != 2 {
			t.Errorf("搜索结果数量不正确: 期望 %d, 实际 %d", 2, len(results))
		}
	})

	// 测试3: 按类别过滤
	t.Run("按类别过滤", func(t *testing.T) {
		query := &SearchQuery{
			Query: "",
			Filters: map[string]any{
				"category": "backend",
			},
			Limit: 10,
		}

		results, err := em.Search(query)
		if err != nil {
			t.Errorf("按类别过滤失败: %v", err)
		}

		// 应该找到1个记录（search-1）
		if len(results) != 1 {
			t.Errorf("过滤结果数量不正确: 期望 %d, 实际 %d", 1, len(results))
		}
		if results[0].Record.ID != "search-1" {
			t.Errorf("过滤结果不正确: 期望 %s, 实际 %s", "search-1", results[0].Record.ID)
		}
	})

	// 测试4: 按成功状态过滤
	t.Run("按成功状态过滤", func(t *testing.T) {
		query := &SearchQuery{
			Query: "",
			Filters: map[string]any{
				"success": true,
			},
			Limit: 10,
		}

		results, err := em.Search(query)
		if err != nil {
			t.Errorf("按成功状态过滤失败: %v", err)
		}

		// 应该找到3个成功记录
		if len(results) != 3 {
			t.Errorf("成功记录数量不正确: 期望 %d, 实际 %d", 3, len(results))
		}

		for _, result := range results {
			if !result.Record.Success {
				t.Errorf("记录 %s 应该标记为成功", result.Record.ID)
			}
		}
	})

	// 测试5: 组合过滤
	t.Run("组合过滤", func(t *testing.T) {
		query := &SearchQuery{
			Query: "",
			Filters: map[string]any{
				"session_id": "session-b",
				"success":    false,
			},
			Limit: 10,
		}

		results, err := em.Search(query)
		if err != nil {
			t.Errorf("组合过滤失败: %v", err)
		}

		// 应该找到1个记录（search-4）
		if len(results) != 1 {
			t.Errorf("组合过滤结果数量不正确: 期望 %d, 实际 %d", 1, len(results))
		}
		if results[0].Record.ID != "search-4" {
			t.Errorf("组合过滤结果不正确: 期望 %s, 实际 %s", "search-4", results[0].Record.ID)
		}
	})

	// 测试6: 按重要性排序
	t.Run("按重要性排序", func(t *testing.T) {
		query := &SearchQuery{
			Query:     "",
			Limit:     10,
			SortBy:    "importance",
			SortOrder: "desc",
		}

		results, err := em.Search(query)
		if err != nil {
			t.Errorf("按重要性排序失败: %v", err)
		}

		// 验证按重要性降序排列
		if len(results) < 2 {
			t.Errorf("排序结果数量太少: %d", len(results))
			return
		}

		for i := 0; i < len(results)-1; i++ {
			if results[i].Record.Importance < results[i+1].Record.Importance {
				t.Errorf("记录没有按重要性降序排列: %f < %f",
					results[i].Record.Importance, results[i+1].Record.Importance)
			}
		}
	})

	// 测试7: 按时间排序
	t.Run("按时间排序", func(t *testing.T) {
		query := &SearchQuery{
			Query:     "",
			Limit:     10,
			SortBy:    "created_at",
			SortOrder: "asc",
		}

		results, err := em.Search(query)
		if err != nil {
			t.Errorf("按时间排序失败: %v", err)
		}

		// 验证按创建时间升序排列
		if len(results) < 2 {
			t.Errorf("排序结果数量太少: %d", len(results))
			return
		}

		for i := 0; i < len(results)-1; i++ {
			if results[i].Record.CreatedAt.After(results[i+1].Record.CreatedAt) {
				t.Error("记录没有按创建时间升序排列")
			}
		}
	})

	// 测试8: 限制结果数量
	t.Run("限制结果数量", func(t *testing.T) {
		query := &SearchQuery{
			Query: "",
			Limit: 2,
		}

		results, err := em.Search(query)
		if err != nil {
			t.Errorf("限制结果数量失败: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("限制结果数量不正确: 期望 %d, 实际 %d", 2, len(results))
		}
	})

	// 测试9: 分页
	t.Run("分页", func(t *testing.T) {
		// 第一页
		query1 := &SearchQuery{
			Query:     "",
			Limit:     2,
			Offset:    0,
			SortBy:    "created_at",
			SortOrder: "desc",
		}

		results1, err := em.Search(query1)
		if err != nil {
			t.Errorf("第一页搜索失败: %v", err)
		}

		// 第二页
		query2 := &SearchQuery{
			Query:     "",
			Limit:     2,
			Offset:    2,
			SortBy:    "created_at",
			SortOrder: "desc",
		}

		results2, err := em.Search(query2)
		if err != nil {
			t.Errorf("第二页搜索失败: %v", err)
		}

		// 验证两页结果不重复
		if len(results1) == 2 && len(results2) == 2 {
			for _, r1 := range results1 {
				for _, r2 := range results2 {
					if r1.Record.ID == r2.Record.ID {
						t.Errorf("分页结果重复: %s", r1.Record.ID)
					}
				}
			}
		}
	})

	// 测试10: 最小重要性过滤
	t.Run("最小重要性过滤", func(t *testing.T) {
		query := &SearchQuery{
			Query:         "",
			MinImportance: 0.8,
			Limit:         10,
		}

		results, err := em.Search(query)
		if err != nil {
			t.Errorf("最小重要性过滤失败: %v", err)
		}

		// 应该找到重要性>=0.8的记录
		for _, result := range results {
			if result.Record.Importance < 0.8 {
				t.Errorf("记录 %s 的重要性 %f 低于最小值 0.8",
					result.Record.ID, result.Record.Importance)
			}
		}
	})
}
