package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// TestFSToolBasic 测试文件系统工具基础功能
func TestFSToolBasic(t *testing.T) {
	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "fs-tool-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建FSTool实例
	tool, err := NewFSTool(tempDir)
	if err != nil {
		t.Fatalf("创建文件系统工具失败: %v", err)
	}

	ctx := context.Background()

	// 测试1: 创建目录
	t.Run("创建目录", func(t *testing.T) {
		execution := ToolExecution{
			ToolID: tool.GetDefinition().ID,
			Parameters: map[string]interface{}{
				"operation": "mkdir",
				"path":      "testdir",
			},
		}

		result, err := tool.Execute(ctx, execution)
		if err != nil {
			t.Errorf("创建目录失败: %v", err)
		}
		if !result.Success {
			t.Errorf("创建目录应该成功: %v", result.Error)
		}

		// 验证目录是否创建
		dirPath := filepath.Join(tempDir, "testdir")
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Error("目录未创建")
		}
	})

	// 测试2: 写入文件
	t.Run("写入文件", func(t *testing.T) {
		execution := ToolExecution{
			ToolID: tool.GetDefinition().ID,
			Parameters: map[string]interface{}{
				"operation": "write",
				"path":      "test.txt",
				"content":   "Hello, World!",
			},
		}

		result, err := tool.Execute(ctx, execution)
		if err != nil {
			t.Errorf("写入文件失败: %v", err)
		}
		if !result.Success {
			t.Errorf("写入文件应该成功: %v", result.Error)
		}

		// 验证文件是否创建
		filePath := filepath.Join(tempDir, "test.txt")
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("读取文件失败: %v", err)
		}
		if string(content) != "Hello, World!" {
			t.Errorf("文件内容不匹配: 期望 'Hello, World!', 实际 '%s'", string(content))
		}
	})

	// 测试3: 读取文件
	t.Run("读取文件", func(t *testing.T) {
		execution := ToolExecution{
			ToolID: tool.GetDefinition().ID,
			Parameters: map[string]interface{}{
				"operation": "read",
				"path":      "test.txt",
			},
		}

		result, err := tool.Execute(ctx, execution)
		if err != nil {
			t.Errorf("读取文件失败: %v", err)
		}
		if !result.Success {
			t.Errorf("读取文件应该成功: %v", result.Error)
		}

		// 验证读取结果
		output, ok := result.Output.(map[string]interface{})
		if !ok {
			t.Error("输出格式不正确")
		}
		if content, ok := output["content"].(string); !ok || content != "Hello, World!" {
			t.Errorf("读取内容不匹配: %v", output)
		}
	})

	// 测试4: 列出目录
	t.Run("列出目录", func(t *testing.T) {
		execution := ToolExecution{
			ToolID: tool.GetDefinition().ID,
			Parameters: map[string]interface{}{
				"operation": "list",
				"path":      ".",
			},
		}

		result, err := tool.Execute(ctx, execution)
		if err != nil {
			t.Errorf("列出目录失败: %v", err)
		}
		if !result.Success {
			t.Errorf("列出目录应该成功: %v", result.Error)
		}

		// 验证列表结果
		output, ok := result.Output.(map[string]interface{})
		if !ok {
			t.Error("输出格式不正确")
		}

		// 检查文件数量
		if count, ok := output["count"].(int); !ok || count != 2 {
			t.Errorf("目录文件数量不正确: 期望2, 实际 %v", output["count"])
		}

		// 检查files字段是否存在
		if files, ok := output["files"]; !ok || files == nil {
			t.Error("files字段不存在")
		}
	})

	// 测试5: 检查文件是否存在
	t.Run("检查文件是否存在", func(t *testing.T) {
		execution := ToolExecution{
			ToolID: tool.GetDefinition().ID,
			Parameters: map[string]interface{}{
				"operation": "exists",
				"path":      "test.txt",
			},
		}

		result, err := tool.Execute(ctx, execution)
		if err != nil {
			t.Errorf("检查文件存在失败: %v", err)
		}
		if !result.Success {
			t.Errorf("检查文件存在应该成功: %v", result.Error)
		}

		// 验证存在结果
		output, ok := result.Output.(map[string]interface{})
		if !ok {
			t.Error("输出格式不正确")
		}
		if exists, ok := output["exists"].(bool); !ok || !exists {
			t.Errorf("文件应该存在: %v", output)
		}
	})

	// 测试6: 复制文件
	t.Run("复制文件", func(t *testing.T) {
		execution := ToolExecution{
			ToolID: tool.GetDefinition().ID,
			Parameters: map[string]interface{}{
				"operation":   "copy",
				"path":        "test.txt",
				"target_path": "test_copy.txt",
			},
		}

		result, err := tool.Execute(ctx, execution)
		if err != nil {
			t.Errorf("复制文件失败: %v", err)
		}
		if !result.Success {
			t.Errorf("复制文件应该成功: %v", result.Error)
		}

		// 验证复制结果
		copyPath := filepath.Join(tempDir, "test_copy.txt")
		if _, err := os.Stat(copyPath); os.IsNotExist(err) {
			t.Error("复制文件未创建")
		}
	})

	// 测试7: 移动文件
	t.Run("移动文件", func(t *testing.T) {
		execution := ToolExecution{
			ToolID: tool.GetDefinition().ID,
			Parameters: map[string]interface{}{
				"operation":   "move",
				"path":        "test_copy.txt",
				"target_path": "test_moved.txt",
			},
		}

		result, err := tool.Execute(ctx, execution)
		if err != nil {
			t.Errorf("移动文件失败: %v", err)
		}
		if !result.Success {
			t.Errorf("移动文件应该成功: %v", result.Error)
		}

		// 验证移动结果
		movedPath := filepath.Join(tempDir, "test_moved.txt")
		if _, err := os.Stat(movedPath); os.IsNotExist(err) {
			t.Error("移动后的文件不存在")
		}
		oldPath := filepath.Join(tempDir, "test_copy.txt")
		if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
			t.Error("原始文件应该被移动")
		}
	})

	// 测试8: 删除文件
	t.Run("删除文件", func(t *testing.T) {
		execution := ToolExecution{
			ToolID: tool.GetDefinition().ID,
			Parameters: map[string]interface{}{
				"operation": "delete",
				"path":      "test_moved.txt",
			},
		}

		result, err := tool.Execute(ctx, execution)
		if err != nil {
			t.Errorf("删除文件失败: %v", err)
		}
		if !result.Success {
			t.Errorf("删除文件应该成功: %v", result.Error)
		}

		// 验证删除结果
		deletedPath := filepath.Join(tempDir, "test_moved.txt")
		if _, err := os.Stat(deletedPath); !os.IsNotExist(err) {
			t.Error("文件应该被删除")
		}
	})

	// 测试9: 验证参数
	t.Run("验证参数", func(t *testing.T) {
		err := tool.ValidateParameters(map[string]interface{}{
			"operation": "read",
			"path":      "test.txt",
		})
		if err != nil {
			t.Errorf("参数验证失败: %v", err)
		}
	})

	// 测试10: 工具基本信息
	t.Run("工具基本信息", func(t *testing.T) {
		def := tool.GetDefinition()
		if def.ID != "fs_tool" {
			t.Errorf("工具ID不正确: 期望 %s, 实际 %s", "fs_tool", def.ID)
		}
		if def.Name != "File System Tool" {
			t.Errorf("工具名称不正确: 期望 %s, 实际 %s", "File System Tool", def.Name)
		}
		if def.Category != ToolCategoryFileSystem {
			t.Errorf("工具类别不正确: 期望 %s, 实际 %s", ToolCategoryFileSystem, def.Category)
		}
		if !tool.IsEnabled() {
			t.Error("工具应该被启用")
		}
	})
}

// TestFSToolErrorCases 测试文件系统工具错误情况
func TestFSToolErrorCases(t *testing.T) {
	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "fs-tool-error-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建FSTool实例
	tool, err := NewFSTool(tempDir)
	if err != nil {
		t.Fatalf("创建文件系统工具失败: %v", err)
	}

	ctx := context.Background()

	// 测试1: 缺少操作参数
	t.Run("缺少操作参数", func(t *testing.T) {
		execution := ToolExecution{
			ToolID:     tool.GetDefinition().ID,
			Parameters: map[string]interface{}{},
		}

		result, err := tool.Execute(ctx, execution)
		if err != nil {
			t.Errorf("执行失败: %v", err)
		}
		if result.Success {
			t.Error("缺少操作参数应该失败")
		}
	})

	// 测试2: 无效的操作类型
	t.Run("无效的操作类型", func(t *testing.T) {
		execution := ToolExecution{
			ToolID: tool.GetDefinition().ID,
			Parameters: map[string]interface{}{
				"operation": "invalid_op",
				"path":      "test.txt",
			},
		}

		result, err := tool.Execute(ctx, execution)
		if err != nil {
			t.Errorf("执行失败: %v", err)
		}
		if result.Success {
			t.Error("无效操作类型应该失败")
		}
	})

	// 测试3: 读取不存在的文件
	t.Run("读取不存在的文件", func(t *testing.T) {
		execution := ToolExecution{
			ToolID: tool.GetDefinition().ID,
			Parameters: map[string]interface{}{
				"operation": "read",
				"path":      "non-existent.txt",
			},
		}

		result, err := tool.Execute(ctx, execution)
		if err != nil {
			t.Errorf("执行失败: %v", err)
		}
		if result.Success {
			t.Error("读取不存在的文件应该失败")
		}
	})

	// 测试4: 写入到无效路径
	t.Run("写入到无效路径", func(t *testing.T) {
		execution := ToolExecution{
			ToolID: tool.GetDefinition().ID,
			Parameters: map[string]interface{}{
				"operation": "write",
				"path":      "/root/system/file.txt", // 尝试写入系统目录
				"content":   "test",
			},
		}

		result, err := tool.Execute(ctx, execution)
		if err != nil {
			t.Errorf("执行失败: %v", err)
		}
		if result.Success {
			t.Error("写入到系统目录应该失败")
		}
	})

	// 测试5: 删除不存在的文件
	t.Run("删除不存在的文件", func(t *testing.T) {
		execution := ToolExecution{
			ToolID: tool.GetDefinition().ID,
			Parameters: map[string]interface{}{
				"operation": "delete",
				"path":      "non-existent.txt",
			},
		}

		result, err := tool.Execute(ctx, execution)
		if err != nil {
			t.Errorf("执行失败: %v", err)
		}
		if result.Success {
			t.Error("删除不存在的文件应该失败")
		}
	})
}

// TestFSToolConcurrent 测试并发执行
func TestFSToolConcurrent(t *testing.T) {
	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "fs-tool-concurrent-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建FSTool实例
	tool, err := NewFSTool(tempDir)
	if err != nil {
		t.Fatalf("创建文件系统工具失败: %v", err)
	}

	ctx := context.Background()
	errors := make(chan error, 10)

	// 并发创建多个文件
	for i := 0; i < 5; i++ {
		go func(index int) {
			execution := ToolExecution{
				ToolID: tool.GetDefinition().ID,
				Parameters: map[string]interface{}{
					"operation": "write",
					"path":      filepath.Join("concurrent", "file_"+string(rune('A'+index))+".txt"),
					"content":   "Content for file " + string(rune('A'+index)),
				},
			}

			result, err := tool.Execute(ctx, execution)
			if err != nil {
				errors <- err
				return
			}
			if !result.Success {
				errors <- err
				return
			}
			errors <- nil
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 5; i++ {
		if err := <-errors; err != nil {
			t.Errorf("并发执行失败: %v", err)
		}
	}

	// 验证所有文件都创建了
	for i := 0; i < 5; i++ {
		filePath := filepath.Join(tempDir, "concurrent", "file_"+string(rune('A'+i))+".txt")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("文件未创建: %s", filePath)
		}
	}
}
