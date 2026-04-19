package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"GenPulse/internal/utils"
)

// GitTool Git操作工具
type GitTool struct {
	*BaseTool
	workspacePath string
}

// NewGitTool 创建Git操作工具
func NewGitTool(workspacePath string) (*GitTool, error) {
	// 确保工作区目录存在
	if err := os.MkdirAll(workspacePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workspace directory: %w", err)
	}

	definition := ToolDefinition{
		ID:          "git_tool",
		Name:        "Git Tool",
		Description: "提供Git版本控制操作功能，包括初始化、提交、推送、拉取等",
		Category:    ToolCategoryGit,
		Version:     "1.0.0",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"operation": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"init", "status", "add", "commit", "push", "pull", "clone", "branch", "checkout", "log", "remote"},
					"description": "Git操作类型",
				},
				"path": map[string]interface{}{
					"type":        "string",
					"description": "仓库路径（相对于工作区）",
				},
				"message": map[string]interface{}{
					"type":        "string",
					"description": "提交消息（commit操作需要）",
				},
				"files": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "string",
					},
					"description": "要添加的文件列表（add操作需要）",
				},
				"remote_url": map[string]interface{}{
					"type":        "string",
					"description": "远程仓库URL（clone和remote操作需要）",
				},
				"remote_name": map[string]interface{}{
					"type":        "string",
					"description": "远程仓库名称（默认origin）",
					"default":     "origin",
				},
				"branch_name": map[string]interface{}{
					"type":        "string",
					"description": "分支名称（branch和checkout操作需要）",
				},
				"username": map[string]interface{}{
					"type":        "string",
					"description": "Git用户名（push/pull操作需要）",
				},
				"token": map[string]interface{}{
					"type":        "string",
					"description": "访问令牌（push/pull操作需要）",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "日志条目限制（log操作）",
					"default":     10,
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
		Tags:    []string{"git", "version-control", "vcs"},
	}

	tool := &GitTool{
		BaseTool:      NewBaseTool(definition),
		workspacePath: workspacePath,
	}

	return tool, nil
}

// Execute 执行Git操作
func (t *GitTool) Execute(ctx context.Context, execution ToolExecution) (*ToolResult, error) {
	// 获取参数
	operation, _ := execution.Parameters["operation"].(string)
	path, _ := execution.Parameters["path"].(string)
	message, _ := execution.Parameters["message"].(string)
	files, _ := execution.Parameters["files"].([]interface{})
	remoteURL, _ := execution.Parameters["remote_url"].(string)
	remoteName, _ := execution.Parameters["remote_name"].(string)
	branchName, _ := execution.Parameters["branch_name"].(string)
	username, _ := execution.Parameters["username"].(string)
	token, _ := execution.Parameters["token"].(string)
	limit, _ := execution.Parameters["limit"].(float64)

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
	case "init":
		result, operationErr = t.initRepository(safePath)
	case "status":
		result, operationErr = t.getStatus(safePath)
	case "add":
		fileList := make([]string, len(files))
		for i, f := range files {
			fileList[i] = f.(string)
		}
		result, operationErr = t.addFiles(safePath, fileList)
	case "commit":
		result, operationErr = t.commitChanges(safePath, message)
	case "push":
		result, operationErr = t.pushToRemote(safePath, remoteName, username, token)
	case "pull":
		result, operationErr = t.pullFromRemote(safePath, remoteName, username, token)
	case "clone":
		result, operationErr = t.cloneRepository(remoteURL, safePath)
	case "branch":
		result, operationErr = t.manageBranch(safePath, branchName)
	case "checkout":
		result, operationErr = t.checkoutBranch(safePath, branchName)
	case "log":
		result, operationErr = t.getLog(safePath, int(limit))
	case "remote":
		result, operationErr = t.manageRemote(safePath, remoteName, remoteURL)
	default:
		operationErr = fmt.Errorf("unsupported Git operation: %s", operation)
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
func (t *GitTool) validatePath(path string) (string, error) {
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

// initRepository 初始化Git仓库
func (t *GitTool) initRepository(path string) (interface{}, error) {
	// 确保目录存在
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// 初始化Git仓库
	repo, err := git.PlainInit(path, false)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// 获取仓库信息
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	// 获取配置
	cfg, err := repo.Config()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	return map[string]interface{}{
		"path":        path,
		"initialized": true,
		"bare":        false,
		"worktree":    worktree.Filesystem.Root(),
		"config": map[string]interface{}{
			"core": cfg.Core,
		},
	}, nil
}

// getStatus 获取仓库状态
func (t *GitTool) getStatus(path string) (interface{}, error) {
	// 打开仓库
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	// 获取工作树
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	// 获取状态
	status, err := worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	// 转换状态为可序列化格式
	statusMap := make(map[string]interface{})
	for file, fileStatus := range status {
		statusMap[file] = map[string]interface{}{
			"staging":  fileStatus.Staging,
			"worktree": fileStatus.Worktree,
			"extra":    fileStatus.Extra,
		}
	}

	// 获取当前分支
	head, err := repo.Head()
	var currentBranch string
	if err == nil {
		currentBranch = head.Name().Short()
	}

	return map[string]interface{}{
		"path":           path,
		"current_branch": currentBranch,
		"is_clean":       status.IsClean(),
		"files":          statusMap,
		"file_count":     len(status),
	}, nil
}

// addFiles 添加文件到暂存区
func (t *GitTool) addFiles(path string, files []string) (interface{}, error) {
	// 打开仓库
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	// 获取工作树
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	// 添加文件
	var addedFiles []string
	if len(files) == 0 {
		// 添加所有文件
		if err := worktree.AddWithOptions(&git.AddOptions{All: true}); err != nil {
			return nil, fmt.Errorf("failed to add all files: %w", err)
		}
		addedFiles = []string{"*"}
	} else {
		// 添加指定文件
		for _, file := range files {
			if _, err := worktree.Add(file); err != nil {
				return nil, fmt.Errorf("failed to add file %s: %w", file, err)
			}
			addedFiles = append(addedFiles, file)
		}
	}

	return map[string]interface{}{
		"path":        path,
		"added_files": addedFiles,
		"count":       len(addedFiles),
	}, nil
}

// commitChanges 提交更改
func (t *GitTool) commitChanges(path, message string) (interface{}, error) {
	if message == "" {
		return nil, fmt.Errorf("commit message is required")
	}

	// 打开仓库
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	// 获取工作树
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	// 提交更改
	commit, err := worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "GenPulse Agent",
			Email: "agent@genpulse.ai",
			When:  time.Now(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	// 获取提交信息
	commitObj, err := repo.CommitObject(commit)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit object: %w", err)
	}

	return map[string]interface{}{
		"path":         path,
		"commit_hash":  commit.String(),
		"message":      commitObj.Message,
		"author":       commitObj.Author.String(),
		"committer":    commitObj.Committer.String(),
		"timestamp":    commitObj.Author.When.Unix(),
		"parent_count": len(commitObj.ParentHashes),
	}, nil
}

// pushToRemote 推送到远程仓库
func (t *GitTool) pushToRemote(path, remoteName, username, token string) (interface{}, error) {
	// 打开仓库
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	// 创建认证信息
	var auth *http.BasicAuth
	if username != "" && token != "" {
		auth = &http.BasicAuth{
			Username: username,
			Password: token,
		}
	}

	// 推送
	err = repo.Push(&git.PushOptions{
		RemoteName: remoteName,
		Auth:       auth,
		Progress:   os.Stdout,
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			// 已经是最新版本
			return map[string]interface{}{
				"path":    path,
				"remote":  remoteName,
				"pushed":  false,
				"message": "already up to date",
			}, nil
		}
		return nil, fmt.Errorf("failed to push: %w", err)
	}

	return map[string]interface{}{
		"path":   path,
		"remote": remoteName,
		"pushed": true,
	}, nil
}

// pullFromRemote 从远程仓库拉取
func (t *GitTool) pullFromRemote(path, remoteName, username, token string) (interface{}, error) {
	// 打开仓库
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	// 获取工作树
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	// 创建认证信息
	var auth *http.BasicAuth
	if username != "" && token != "" {
		auth = &http.BasicAuth{
			Username: username,
			Password: token,
		}
	}

	// 拉取
	err = worktree.Pull(&git.PullOptions{
		RemoteName: remoteName,
		Auth:       auth,
		Progress:   os.Stdout,
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			// 已经是最新版本
			return map[string]interface{}{
				"path":    path,
				"remote":  remoteName,
				"pulled":  false,
				"message": "already up to date",
			}, nil
		}
		return nil, fmt.Errorf("failed to pull: %w", err)
	}

	return map[string]interface{}{
		"path":   path,
		"remote": remoteName,
		"pulled": true,
	}, nil
}

// cloneRepository 克隆仓库
func (t *GitTool) cloneRepository(url, path string) (interface{}, error) {
	if url == "" {
		return nil, fmt.Errorf("remote URL is required for clone")
	}

	// 确保目录不存在或为空
	if _, err := os.Stat(path); err == nil {
		// 目录存在，检查是否为空
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, fmt.Errorf("failed to check directory: %w", err)
		}
		if len(entries) > 0 {
			return nil, fmt.Errorf("target directory is not empty")
		}
	}

	// 克隆仓库
	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	// 获取仓库信息
	head, err := repo.Head()
	var currentBranch string
	if err == nil {
		currentBranch = head.Name().Short()
	}

	// 获取远程信息
	remotes, err := repo.Remotes()
	var remoteURLs []string
	if err == nil && len(remotes) > 0 {
		for _, remote := range remotes {
			remoteURLs = append(remoteURLs, remote.Config().URLs...)
		}
	}

	return map[string]interface{}{
		"path":           path,
		"url":            url,
		"cloned":         true,
		"current_branch": currentBranch,
		"remotes":        remoteURLs,
	}, nil
}

// manageBranch 管理分支
func (t *GitTool) manageBranch(path, branchName string) (interface{}, error) {
	if branchName == "" {
		// 列出分支
		return t.listBranches(path)
	}

	// 创建分支
	return t.createBranch(path, branchName)
}

// listBranches 列出分支
func (t *GitTool) listBranches(path string) (interface{}, error) {
	// 打开仓库
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	// 获取分支迭代器
	branches, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	// 获取当前分支
	head, err := repo.Head()
	var currentBranch string
	if err == nil {
		currentBranch = head.Name().Short()
	}

	// 收集分支信息
	var branchList []map[string]interface{}
	err = branches.ForEach(func(ref *plumbing.Reference) error {
		branchList = append(branchList, map[string]interface{}{
			"name":    ref.Name().Short(),
			"hash":    ref.Hash().String(),
			"is_head": ref.Name().Short() == currentBranch,
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate branches: %w", err)
	}

	return map[string]interface{}{
		"path":           path,
		"current_branch": currentBranch,
		"branches":       branchList,
		"count":          len(branchList),
	}, nil
}

// createBranch 创建分支
func (t *GitTool) createBranch(path, branchName string) (interface{}, error) {
	// 打开仓库
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	// 获取当前HEAD
	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	// 创建分支引用
	branchRef := plumbing.NewBranchReferenceName(branchName)
	ref := plumbing.NewHashReference(branchRef, head.Hash())

	// 存储引用
	if err := repo.Storer.SetReference(ref); err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}

	return map[string]interface{}{
		"path":        path,
		"branch_name": branchName,
		"created":     true,
		"from_hash":   head.Hash().String(),
	}, nil
}

// checkoutBranch 切换分支
func (t *GitTool) checkoutBranch(path, branchName string) (interface{}, error) {
	if branchName == "" {
		return nil, fmt.Errorf("branch name is required for checkout")
	}

	// 打开仓库
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	// 获取工作树
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	// 切换分支
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Create: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to checkout branch: %w", err)
	}

	return map[string]interface{}{
		"path":        path,
		"branch_name": branchName,
		"checked_out": true,
	}, nil
}

// getLog 获取提交日志
func (t *GitTool) getLog(path string, limit int) (interface{}, error) {
	// 打开仓库
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	// 获取当前分支
	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	// 获取提交日志
	commitIter, err := repo.Log(&git.LogOptions{From: head.Hash()})
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	// 收集提交信息
	var commits []map[string]interface{}
	count := 0

	err = commitIter.ForEach(func(commit *object.Commit) error {
		if limit > 0 && count >= limit {
			return nil
		}

		commits = append(commits, map[string]interface{}{
			"hash":         commit.Hash.String(),
			"message":      commit.Message,
			"author":       commit.Author.String(),
			"committer":    commit.Committer.String(),
			"timestamp":    commit.Author.When.Unix(),
			"parent_count": len(commit.ParentHashes),
		})

		count++
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate commits: %w", err)
	}

	return map[string]interface{}{
		"path":    path,
		"commits": commits,
		"count":   len(commits),
		"limit":   limit,
	}, nil
}

// manageRemote 管理远程仓库
func (t *GitTool) manageRemote(path, remoteName, remoteURL string) (interface{}, error) {
	// 打开仓库
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	if remoteURL == "" {
		// 列出远程仓库
		return t.listRemotes(repo)
	}

	// 添加远程仓库
	return t.addRemote(repo, remoteName, remoteURL)
}

// listRemotes 列出远程仓库
func (t *GitTool) listRemotes(repo *git.Repository) (interface{}, error) {
	remotes, err := repo.Remotes()
	if err != nil {
		return nil, fmt.Errorf("failed to get remotes: %w", err)
	}

	var remoteList []map[string]interface{}
	for _, remote := range remotes {
		remoteList = append(remoteList, map[string]interface{}{
			"name": remote.Config().Name,
			"urls": remote.Config().URLs,
		})
	}

	return map[string]interface{}{
		"remotes": remoteList,
		"count":   len(remoteList),
	}, nil
}

// addRemote 添加远程仓库
func (t *GitTool) addRemote(repo *git.Repository, remoteName, remoteURL string) (interface{}, error) {
	// 创建远程配置
	remote := &config.RemoteConfig{
		Name: remoteName,
		URLs: []string{remoteURL},
	}

	// 添加远程仓库
	if _, err := repo.CreateRemote(remote); err != nil {
		return nil, fmt.Errorf("failed to create remote: %w", err)
	}

	return map[string]interface{}{
		"remote_name": remoteName,
		"remote_url":  remoteURL,
		"added":       true,
	}, nil
}

// Initialize 初始化工具
func (t *GitTool) Initialize() error {
	// 调用父类初始化
	if err := t.BaseTool.Initialize(); err != nil {
		return err
	}

	utils.Info("Git工具初始化完成，工作区: %s", t.workspacePath)
	return nil
}
