package history

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestHistoryStorage(t *testing.T) {
	// 创建临时数据库文件
	tmpfile, err := os.CreateTemp("", "history-test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// 创建历史存储
	storage, err := NewHistoryStorage(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to create history storage: %v", err)
	}
	defer storage.Close()

	ctx := context.Background()
	now := time.Now()

	// 测试创建记录
	record := &ExecutionRecord{
		ID:            "test-id-001",
		Name:          "测试记录",
		Description:   "这是一个测试记录",
		Status:        "completed",
		StartTime:     now.Add(-1 * time.Hour),
		EndTime:       now,
		Duration:      1 * time.Hour,
		AgentCount:    3,
		ToolCallCount: 15,
		TokenUsage:    5000,
		CostEstimate:  0.25,
		Metadata: map[string]interface{}{
			"test": true,
		},
		Tags: []string{"test", "unit"},
	}

	err = storage.CreateRecord(ctx, record)
	if err != nil {
		t.Fatalf("Failed to create record: %v", err)
	}

	// 测试获取记录
	retrieved, err := storage.GetRecord(ctx, "test-id-001")
	if err != nil {
		t.Fatalf("Failed to get record: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Record should exist")
	}

	if retrieved.Name != "测试记录" {
		t.Errorf("Expected name '测试记录', got %s", retrieved.Name)
	}

	if retrieved.AgentCount != 3 {
		t.Errorf("Expected agent count 3, got %d", retrieved.AgentCount)
	}

	// 测试查询记录
	query := ExecutionQuery{
		Limit: 10,
	}

	records, total, err := storage.QueryRecords(ctx, query)
	if err != nil {
		t.Fatalf("Failed to query records: %v", err)
	}

	if total != 1 {
		t.Errorf("Expected 1 record total, got %d", total)
	}

	if len(records) != 1 {
		t.Errorf("Expected 1 record in results, got %d", len(records))
	}

	// 测试更新记录
	updates := map[string]interface{}{
		"status":          "failed",
		"tool_call_count": 20,
		"token_usage":     6000,
	}

	err = storage.UpdateRecord(ctx, "test-id-001", updates)
	if err != nil {
		t.Fatalf("Failed to update record: %v", err)
	}

	// 验证更新
	updated, err := storage.GetRecord(ctx, "test-id-001")
	if err != nil {
		t.Fatalf("Failed to get updated record: %v", err)
	}

	if updated.Status != "failed" {
		t.Errorf("Expected status 'failed', got %s", updated.Status)
	}

	if updated.ToolCallCount != 20 {
		t.Errorf("Expected tool call count 20, got %d", updated.ToolCallCount)
	}

	// 测试删除记录
	err = storage.DeleteRecord(ctx, "test-id-001")
	if err != nil {
		t.Fatalf("Failed to delete record: %v", err)
	}

	// 验证删除
	deleted, err := storage.GetRecord(ctx, "test-id-001")
	if err != nil {
		t.Fatalf("Failed to check deleted record: %v", err)
	}

	if deleted != nil {
		t.Error("Record should be deleted")
	}

	// 测试统计信息
	stats, err := storage.GetStatistics(ctx, now.Add(-24*time.Hour), now)
	if err != nil {
		t.Fatalf("Failed to get statistics: %v", err)
	}

	if totalRecords, ok := stats["total_records"].(int64); !ok || totalRecords != 0 {
		t.Errorf("Expected 0 total records after deletion, got %v", stats["total_records"])
	}
}

func TestHistoryService(t *testing.T) {
	// 创建临时数据库文件
	tmpfile, err := os.CreateTemp("", "history-service-test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// 创建历史存储
	storage, err := NewHistoryStorage(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to create history storage: %v", err)
	}
	defer storage.Close()

	// 创建历史服务
	service := NewHistoryServiceWithStorage(storage, nil)

	ctx := context.Background()
	now := time.Now()

	// 测试创建记录
	record := &ExecutionRecord{
		Name:          "服务测试记录",
		Description:   "测试历史服务功能",
		Status:        "running",
		StartTime:     now,
		AgentCount:    2,
		ToolCallCount: 5,
		TokenUsage:    2000,
		Tags:          []string{"service", "test"},
	}

	created, err := service.CreateRecord(ctx, record)
	if err != nil {
		t.Fatalf("Failed to create record via service: %v", err)
	}

	if created.ID == "" {
		t.Error("Record ID should be generated")
	}

	// 测试查询记录
	query := ExecutionQuery{
		Statuses: []string{"running"},
		Limit:    10,
	}

	_, total, err := service.QueryRecords(ctx, query)
	if err != nil {
		t.Fatalf("Failed to query records via service: %v", err)
	}

	if total != 1 {
		t.Errorf("Expected 1 record total, got %d", total)
	}

	// 测试更新记录
	updates := map[string]interface{}{
		"status":      "completed",
		"end_time":    now.Add(30 * time.Minute),
		"token_usage": 3000,
	}

	updated, err := service.UpdateRecord(ctx, created.ID, updates)
	if err != nil {
		t.Fatalf("Failed to update record via service: %v", err)
	}

	if updated.Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", updated.Status)
	}

	if updated.Duration != 30*time.Minute {
		t.Errorf("Expected duration 30m, got %v", updated.Duration)
	}

	// 测试获取统计信息
	stats, err := service.GetStatistics(ctx, now.Add(-1*time.Hour), now.Add(1*time.Hour))
	if err != nil {
		t.Fatalf("Failed to get statistics via service: %v", err)
	}

	if successRate, ok := stats["success_rate"].(float64); !ok || successRate != 100.0 {
		t.Errorf("Expected 100%% success rate, got %v", stats["success_rate"])
	}
}
