package tools

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"GenPulse/internal/utils"
)

// FSTool 文件系统工具
type FSTool struct {
	*BaseTool
	workspacePath string
}

// NewFSTool 创建文件系统工具
func NewFSTool(workspacePath string) (*FSTool, error) {
	// 确保工作区目录存在
	if err := os.MkdirAll(workspacePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workspace directory: %w", err)
	}

	definition := ToolDefinition{
		ID:          "fs_tool",
		Name:        "File System Tool",
		Description: "提供文件系统操作功能，包括读取、写入、列出、删除文件等",
		Category:    ToolCategoryFileSystem,
		Version:     "1.0.0",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"operation": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"read", "write", "list", "delete", "mkdir", "exists", "copy", "move"},
					"description": "操作类型",
				},
				"path": map[string]interface{}{
					"type":        "string",
					"description": "文件或目录路径",
				},
				"content": map[string]interface{}{
					"type":        "string",
					"description": "要写入的内容（仅write操作需要）",
				},
				"target_path": map[string]interface{}{
					"type":        "string",
					"description": "目标路径（copy和move操作需要）",
				},
				"recursive": map[string]interface{}{
					"type":        "boolean",
					"description": "是否递归操作（list和delete操作）",
					"default":     false,
				},
			},
			"required": []string{"operation", "path"},
		},
		Returns: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"success": map[string]interface{}{
					"type":        "boolean",
					"description": "操作是否成功",
				},
				"result": map[string]interface{}{
					"type":        "object",
					"description": "操作结果",
				},
				"error": map[string]interface{}{
					"type":        "string",
					"description": "错误信息",
				},
			},
		},
		Enabled: true,
		Tags:    []string{"filesystem", "io", "storage"},
	}

	tool := &FSTool{
		BaseTool:      NewBaseTool(definition),
		workspacePath: workspacePath,
	}

	return tool, nil
}

// Execute 执行文件系统操作
func (t *FSTool) Execute(ctx context.Context, execution ToolExecution) (*ToolResult, error) {
	// 获取参数
	operation, _ := execution.Parameters["operation"].(string)
	path, _ := execution.Parameters["path"].(string)
	content, _ := execution.Parameters["content"].(string)
	targetPath, _ := execution.Parameters["target_path"].(string)
	recursive, _ := execution.Parameters["recursive"].(bool)

	// 验证路径安全性
	safePath, err := t.validatePath(path)
	if err != nil {
		return &ToolResult{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	var result interface{}
	var operationErr error

	// 执行操作
	switch operation {
	case "read":
		result, operationErr = t.readFile(safePath)
	case "write":
		result, operationErr = t.writeFile(safePath, content)
	case "list":
		result, operationErr = t.listFiles(safePath, recursive)
	case "delete":
		result, operationErr = t.deleteFile(safePath, recursive)
	case "mkdir":
		result, operationErr = t.createDirectory(safePath)
	case "exists":
		result, operationErr = t.fileExists(safePath)
	case "copy":
		safeTargetPath, err := t.validatePath(targetPath)
		if err != nil {
			operationErr = err
		} else {
			result, operationErr = t.copyFile(safePath, safeTargetPath)
		}
	case "move":
		safeTargetPath, err := t.validatePath(targetPath)
		if err != nil {
			operationErr = err
		} else {
			result, operationErr = t.moveFile(safePath, safeTargetPath)
		}
	default:
		operationErr = fmt.Errorf("unsupported operation: %s", operation)
	}

	// 构建结果
	toolResult := &ToolResult{
		Success: operationErr == nil,
		Output:  result,
	}

	if operationErr != nil {
		toolResult.Error = operationErr.Error()
	}

	return toolResult, nil
}

// validatePath 验证路径安全性
func (t *FSTool) validatePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// 清理路径
	cleanPath := filepath.Clean(path)

	// 检查是否为绝对路径，如果是则转换为相对于工作区的路径
	if filepath.IsAbs(cleanPath) {
		// 确保路径在工作区内
		relPath, err := filepath.Rel(t.workspacePath, cleanPath)
		if err != nil {
			return "", fmt.Errorf("path is outside workspace: %s", cleanPath)
		}

		// 检查是否尝试向上访问
		if strings.HasPrefix(relPath, "..") {
			return "", fmt.Errorf("path traversal not allowed: %s", cleanPath)
		}

		cleanPath = filepath.Join(t.workspacePath, relPath)
	} else {
		// 相对路径，直接连接到工作区
		cleanPath = filepath.Join(t.workspacePath, cleanPath)
	}

	return cleanPath, nil
}

// readFile 读取文件
func (t *FSTool) readFile(path string) (interface{}, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return map[string]interface{}{
		"path":    path,
		"content": string(content),
		"size":    len(content),
	}, nil
}

// writeFile 写入文件
func (t *FSTool) writeFile(path, content string) (interface{}, error) {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	// 获取文件信息
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return map[string]interface{}{
		"path":     path,
		"size":     info.Size(),
		"modified": info.ModTime().Unix(),
	}, nil
}

// listFiles 列出文件
func (t *FSTool) listFiles(path string, recursive bool) (interface{}, error) {
	var files []map[string]interface{}

	if recursive {
		// 递归列出
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 跳过目录本身（如果是目录）
			if filePath == path {
				return nil
			}

			files = append(files, map[string]interface{}{
				"path":     filePath,
				"name":     info.Name(),
				"size":     info.Size(),
				"is_dir":   info.IsDir(),
				"modified": info.ModTime().Unix(),
				"mode":     info.Mode().String(),
			})

			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("failed to walk directory: %w", err)
		}
	} else {
		// 非递归列出
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory: %w", err)
		}

		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			files = append(files, map[string]interface{}{
				"path":     filepath.Join(path, entry.Name()),
				"name":     entry.Name(),
				"size":     info.Size(),
				"is_dir":   entry.IsDir(),
				"modified": info.ModTime().Unix(),
				"mode":     info.Mode().String(),
			})
		}
	}

	return map[string]interface{}{
		"path":  path,
		"files": files,
		"count": len(files),
	}, nil
}

// deleteFile 删除文件或目录
func (t *FSTool) deleteFile(path string, recursive bool) (interface{}, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	if info.IsDir() {
		if recursive {
			if err := os.RemoveAll(path); err != nil {
				return nil, fmt.Errorf("failed to delete directory recursively: %w", err)
			}
		} else {
			// 检查目录是否为空
			entries, err := os.ReadDir(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read directory: %w", err)
			}

			if len(entries) > 0 {
				return nil, fmt.Errorf("directory is not empty, use recursive=true to delete")
			}

			if err := os.Remove(path); err != nil {
				return nil, fmt.Errorf("failed to delete directory: %w", err)
			}
		}
	} else {
		if err := os.Remove(path); err != nil {
			return nil, fmt.Errorf("failed to delete file: %w", err)
		}
	}

	return map[string]interface{}{
		"path":    path,
		"deleted": true,
		"was_dir": info.IsDir(),
		"size":    info.Size(),
	}, nil
}

// createDirectory 创建目录
func (t *FSTool) createDirectory(path string) (interface{}, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory info: %w", err)
	}

	return map[string]interface{}{
		"path":     path,
		"created":  true,
		"mode":     info.Mode().String(),
		"modified": info.ModTime().Unix(),
	}, nil
}

// fileExists 检查文件是否存在
func (t *FSTool) fileExists(path string) (interface{}, error) {
	_, err := os.Stat(path)
	exists := err == nil

	return map[string]interface{}{
		"path":   path,
		"exists": exists,
	}, nil
}

// copyFile 复制文件
func (t *FSTool) copyFile(src, dst string) (interface{}, error) {
	// 检查源文件是否存在
	srcInfo, err := os.Stat(src)
	if err != nil {
		return nil, fmt.Errorf("source file not found: %w", err)
	}

	// 确保目标目录存在
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	// 复制文件
	srcFile, err := os.Open(src)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// 保持文件权限
	if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
		utils.Warn("Failed to preserve file permissions: %v", err)
	}

	dstInfo, err := os.Stat(dst)
	if err != nil {
		return nil, fmt.Errorf("failed to get destination file info: %w", err)
	}

	return map[string]interface{}{
		"source":      src,
		"destination": dst,
		"copied":      true,
		"size":        dstInfo.Size(),
	}, nil
}

// moveFile 移动文件
func (t *FSTool) moveFile(src, dst string) (interface{}, error) {
	// 检查源文件是否存在
	srcInfo, err := os.Stat(src)
	if err != nil {
		return nil, fmt.Errorf("source file not found: %w", err)
	}

	// 确保目标目录存在
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	// 移动文件
	if err := os.Rename(src, dst); err != nil {
		return nil, fmt.Errorf("failed to move file: %w", err)
	}

	return map[string]interface{}{
		"source":      src,
		"destination": dst,
		"moved":       true,
		"size":        srcInfo.Size(),
	}, nil
}

// Initialize 初始化工具
func (t *FSTool) Initialize() error {
	// 调用父类初始化
	if err := t.BaseTool.Initialize(); err != nil {
		return err
	}

	utils.Info("文件系统工具初始化完成，工作区: %s", t.workspacePath)
	return nil
}
