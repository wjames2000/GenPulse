package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"GenPulse/internal/memory"
)

func main() {
	fmt.Println("=== 三层记忆架构演示程序 ===")

	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "memory_demo_*")
	if err != nil {
		log.Fatal("创建临时目录失败:", err)
	}
	defer os.RemoveAll(tempDir)

	fmt.Printf("使用临时目录: %s\n", tempDir)

	// 1. 初始化三层记忆
	fmt.Println("\n1. 初始化三层记忆架构...")

	// L1: 工作记忆
	wm := memory.NewWorkingMemorySimple()

	// L2: 情节记忆
	episodicDBPath := filepath.Join(tempDir, "episodic.db")
	em, err := memory.NewEpisodicMemory(episodicDBPath)
	if err != nil {
		log.Fatal("初始化情节记忆失败:", err)
	}
	defer em.Close()

	// L3: 语义记忆
	semanticDir := filepath.Join(tempDir, "semantic")
	sm, err := memory.NewSemanticMemory(semanticDir)
	if err != nil {
		log.Fatal("初始化语义记忆失败:", err)
	}

	// 记忆检索引擎
	se := memory.NewSearchEngine(wm, em, sm)

	// 自动更新管理器
	aum := memory.NewAutoUpdateManager(wm, em, sm, se)

	fmt.Println("✓ 三层记忆架构初始化完成")

	// 2. 演示工作记忆（L1）
	fmt.Println("\n2. 演示工作记忆（L1）...")

	// 创建会话
	sessionID := "demo_session_001"
	ctx := context.Background()
	session := wm.GetOrCreateSession(sessionID, ctx)

	// 存储工作记忆数据
	session.Set("current_task", "实现用户登录功能")
	session.Set("programming_language", "Go")
	session.Set("framework", "Gin")

	// 读取工作记忆
	task, _ := session.Get("current_task")
	fmt.Printf("当前任务: %v\n", task)

	// 获取会话数据
	sessionData := session.GetAll()
	fmt.Printf("会话数据数量: %d\n", len(sessionData))

	// 3. 演示情节记忆（L2）
	fmt.Println("\n3. 演示情节记忆（L2）...")

	// 创建一些记忆记录
	records := []*memory.MemoryRecord{
		{
			ID:          "task_001",
			SessionID:   sessionID,
			TaskID:      "implement_login",
			TaskType:    "代码实现",
			Description: "实现用户登录功能",
			Content:     "使用JWT token实现用户认证，包含密码加密和会话管理",
			Metadata: map[string]any{
				"complexity": "medium",
				"priority":   "high",
			},
			Tags:       []string{"authentication", "security", "jwt"},
			Category:   "后端开发",
			Importance: 0.8,
			Success:    true,
			CreatedAt:  time.Now().Add(-2 * time.Hour),
			UpdatedAt:  time.Now().Add(-2 * time.Hour),
			AccessedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:          "task_002",
			SessionID:   sessionID,
			TaskID:      "debug_payment",
			TaskType:    "问题调试",
			Description: "修复支付接口超时问题",
			Content:     "发现支付接口在高并发下超时，优化数据库查询和增加缓存",
			Metadata: map[string]any{
				"complexity": "high",
				"priority":   "critical",
			},
			Tags:       []string{"debugging", "performance", "database"},
			Category:   "后端开发",
			Importance: 0.9,
			Success:    true,
			CreatedAt:  time.Now().Add(-1 * time.Hour),
			UpdatedAt:  time.Now().Add(-1 * time.Hour),
			AccessedAt: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:          "task_003",
			SessionID:   sessionID,
			TaskID:      "design_api",
			TaskType:    "系统设计",
			Description: "设计商品管理API",
			Content:     "设计RESTful API用于商品CRUD操作，包含图片上传和分类管理",
			Metadata: map[string]any{
				"complexity": "low",
				"priority":   "medium",
			},
			Tags:       []string{"api-design", "restful", "documentation"},
			Category:   "系统设计",
			Importance: 0.6,
			Success:    true,
			CreatedAt:  time.Now().Add(-30 * time.Minute),
			UpdatedAt:  time.Now().Add(-30 * time.Minute),
			AccessedAt: time.Now().Add(-30 * time.Minute),
		},
	}

	// 存储记录
	for _, record := range records {
		if err := em.Store(record); err != nil {
			log.Printf("存储记录失败 %s: %v", record.ID, err)
		} else {
			fmt.Printf("✓ 存储记录: %s\n", record.Description)
		}
	}

	// 搜索记录
	fmt.Println("\n搜索'支付'相关记录:")
	searchQuery := &memory.SearchQuery{
		Query:  "支付",
		Limit:  5,
		SortBy: "relevance",
	}

	searchResults, err := em.Search(searchQuery)
	if err != nil {
		log.Printf("搜索失败: %v", err)
	} else {
		fmt.Printf("找到 %d 条记录:\n", len(searchResults))
		for i, result := range searchResults {
			fmt.Printf("  %d. %s (相关度: %.1f%%)\n", i+1, result.Record.Description, result.Relevance*100)
		}
	}

	// 获取统计信息
	stats, err := em.GetStats()
	if err != nil {
		log.Printf("获取统计失败: %v", err)
	} else {
		fmt.Printf("\n情节记忆统计:\n")
		fmt.Printf("  总记录数: %v\n", stats["total_count"])
		fmt.Printf("  成功率: %.1f%%\n", stats["success_rate"].(float64)*100)
		fmt.Printf("  平均重要性: %.1f\n", stats["avg_importance"].(float64))
	}

	// 4. 演示语义记忆（L3）
	fmt.Println("\n4. 演示语义记忆（L3）...")

	// 获取用户画像
	profile, err := sm.GetUserProfile()
	if err != nil {
		log.Printf("获取用户画像失败: %v", err)
	} else {
		fmt.Printf("用户ID: %s\n", profile.UserID)
		fmt.Printf("技能数量: %d\n", len(profile.Skills))
		fmt.Printf("知识领域: %d个\n", len(profile.KnowledgeAreas))
	}

	// 获取记忆内容
	memoryContent, err := sm.GetMemoryContent()
	if err != nil {
		log.Printf("获取记忆内容失败: %v", err)
	} else {
		fmt.Printf("记忆内容长度: %d 字符\n", len(memoryContent))
	}

	// 添加学习事件
	learningEvent := &memory.LearningEvent{
		EventID:     "learn_001",
		TaskType:    "代码实现",
		Description: "实现JWT认证",
		Insight:     "JWT token需要设置合理的过期时间，避免安全风险",
		Impact:      "positive",
		Confidence:  0.8,
		Tags:        []string{"security", "authentication", "jwt"},
		CreatedAt:   time.Now(),
	}

	if err := sm.AddLearningEvent(learningEvent); err != nil {
		log.Printf("添加学习事件失败: %v", err)
	} else {
		fmt.Println("✓ 添加学习事件成功")
	}

	// 获取任务建议
	advice, err := sm.GetTaskAdvice("代码实现")
	if err != nil {
		log.Printf("获取任务建议失败: %v", err)
	} else {
		fmt.Printf("\n代码实现任务建议:\n%s\n", advice)
	}

	// 5. 演示记忆检索引擎
	fmt.Println("\n5. 演示记忆检索引擎...")

	// 搜索所有层记忆
	query := &memory.MemoryQuery{
		Query:     "认证",
		TaskType:  "代码实现",
		SessionID: sessionID,
		Limit:     10,
		IncludeL1: true,
		IncludeL2: true,
		IncludeL3: true,
	}

	response, err := se.Search(ctx, query)
	if err != nil {
		log.Printf("记忆搜索失败: %v", err)
	} else {
		fmt.Printf("搜索查询: %s\n", query.Query)
		fmt.Printf("找到 %d 条结果:\n", len(response.Results))

		for i, result := range response.Results {
			// 截断长内容
			content := result.Content
			if len(content) > 100 {
				content = content[:100] + "..."
			}

			fmt.Printf("  %d. [%s] %s (相关度: %.1f%%)\n",
				i+1, result.Source, content, result.Relevance*100)
		}

		fmt.Printf("\n搜索统计:\n")
		for source, stat := range response.Stats {
			if sourceMap, ok := stat.(map[string]any); ok {
				if count, ok := sourceMap["count"].(int); ok {
					fmt.Printf("  %s: %d 条结果\n", source, count)
				}
			}
		}
	}

	// 6. 演示自动更新机制
	fmt.Println("\n6. 演示自动更新机制...")

	// 模拟任务结果
	taskResult := &memory.TaskResult{
		TaskID:      "task_004",
		SessionID:   sessionID,
		TaskType:    "代码实现",
		Description: "实现用户注册功能",
		Content:     "实现用户注册API，包含邮箱验证和密码强度检查",
		Success:     true,
		Duration:    45 * time.Minute,
		Importance:  0.7,
		Tags:        []string{"authentication", "validation", "api"},
		Category:    "后端开发",
		Metadata: map[string]any{
			"lines_of_code": 150,
			"tests_written": true,
		},
		Insights: []string{
			"密码强度检查很重要，可以防止弱密码",
			"邮箱验证可以防止虚假注册",
			"API响应应该包含清晰的错误信息",
		},
	}

	// 记录任务结果
	if err := aum.RecordTaskResult(ctx, taskResult); err != nil {
		log.Printf("记录任务结果失败: %v", err)
	} else {
		fmt.Println("✓ 任务结果已自动记录到三层记忆")

		// 验证记录 - 由于ID是动态生成的，我们检查是否有新记录
		recent, err := em.GetRecent(1)
		if err != nil {
			fmt.Println("记录验证: 获取最近记录失败")
		} else if len(recent) > 0 {
			fmt.Printf("记录验证: 最新记录是 '%s'\n", recent[0].Description)
		}
	}

	// 7. 演示批量操作
	fmt.Println("\n7. 演示批量操作...")

	// 获取最近记录
	recentMemories, err := em.GetRecent(5)
	if err != nil {
		log.Printf("获取最近记录失败: %v", err)
	} else {
		fmt.Printf("最近 %d 条记录:\n", len(recentMemories))
		for i, mem := range recentMemories {
			fmt.Printf("  %d. %s (%s)\n", i+1, mem.Description, mem.TaskType)
		}
	}

	// 获取成功记录
	successfulMemories, err := em.GetSuccessful(3)
	if err != nil {
		log.Printf("获取成功记录失败: %v", err)
	} else {
		fmt.Printf("\n成功记录 (前 %d 条):\n", len(successfulMemories))
		for i, mem := range successfulMemories {
			fmt.Printf("  %d. %s (重要性: %.1f)\n", i+1, mem.Description, mem.Importance)
		}
	}

	// 8. 清理演示
	fmt.Println("\n8. 清理演示...")

	// 获取更新统计
	updateStats := aum.GetUpdateStats()
	fmt.Printf("自动更新配置: %+v\n", updateStats["config"])

	// 清理旧记录（演示，实际上不会清理新记录）
	cleaned, err := aum.CleanupOldRecords(1, 0.1) // 清理1天前重要性<0.1的记录
	if err != nil {
		log.Printf("清理旧记录失败: %v", err)
	} else {
		fmt.Printf("清理了 %d 条旧记录\n", cleaned)
	}

	fmt.Println("\n=== 演示完成 ===")
	fmt.Println("三层记忆架构功能:")
	fmt.Println("✓ L1 工作记忆: 会话级上下文存储")
	fmt.Println("✓ L2 情节记忆: SQLite + FTS5全文检索")
	fmt.Println("✓ L3 语义记忆: 文件系统用户画像")
	fmt.Println("✓ 记忆检索引擎: 多层检索 (L3→L2→L1)")
	fmt.Println("✓ 自动更新机制: 任务完成后自动记录经验")
}
