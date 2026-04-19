package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"GenPulse/internal/utils"
)

// ProjectType 项目类型
type ProjectType string

const (
	ProjectTypeGo        ProjectType = "go"
	ProjectTypeNodeJS    ProjectType = "nodejs"
	ProjectTypePython    ProjectType = "python"
	ProjectTypeReact     ProjectType = "react"
	ProjectTypeVue       ProjectType = "vue"
	ProjectTypeStatic    ProjectType = "static"
	ProjectTypeFullStack ProjectType = "fullstack"
)

// ProjectTemplate 项目模板
type ProjectTemplate struct {
	Type         ProjectType         `json:"type"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Structure    map[string]string   `json:"structure"`    // 文件路径 -> 内容模板
	Dependencies map[string][]string `json:"dependencies"` // 包管理器 -> 依赖列表
	Commands     map[string][]string `json:"commands"`     // 环境 -> 命令列表
}

// ProjectTool 项目管理工具
type ProjectTool struct {
	*BaseTool
	workspacePath string
	templates     map[ProjectType]*ProjectTemplate
}

// NewProjectTool 创建项目管理工具
func NewProjectTool(workspacePath string) (*ProjectTool, error) {
	// 确保工作区目录存在
	if err := os.MkdirAll(workspacePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workspace directory: %w", err)
	}

	// 加载默认模板
	templates := loadDefaultTemplates()

	definition := ToolDefinition{
		ID:          "project_tool",
		Name:        "Project Tool",
		Description: "提供项目管理功能，包括项目初始化、依赖安装等",
		Category:    ToolCategoryProject,
		Version:     "1.0.0",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"operation": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"init", "install_deps", "create_file", "run_command", "get_info"},
					"description": "操作类型",
				},
				"project_path": map[string]interface{}{
					"type":        "string",
					"description": "项目路径（相对于工作区）",
				},
				"project_type": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"go", "nodejs", "python", "react", "vue", "static", "fullstack"},
					"description": "项目类型（init操作需要）",
				},
				"project_name": map[string]interface{}{
					"type":        "string",
					"description": "项目名称（init操作需要）",
				},
				"package_manager": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"npm", "yarn", "pnpm", "go", "pip", "pip3"},
					"description": "包管理器（install_deps操作需要）",
				},
				"file_path": map[string]interface{}{
					"type":        "string",
					"description": "文件路径（create_file操作需要）",
				},
				"file_content": map[string]interface{}{
					"type":        "string",
					"description": "文件内容（create_file操作需要）",
				},
				"command": map[string]interface{}{
					"type":        "string",
					"description": "命令（run_command操作需要）",
				},
				"args": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "string",
					},
					"description": "命令参数（run_command操作需要）",
				},
			},
			"required": []string{"operation", "project_path"},
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
		Tags:    []string{"project", "management", "dependency"},
	}

	tool := &ProjectTool{
		BaseTool:      NewBaseTool(definition),
		workspacePath: workspacePath,
		templates:     templates,
	}

	return tool, nil
}

// Execute 执行项目管理操作
func (t *ProjectTool) Execute(ctx context.Context, execution ToolExecution) (*ToolResult, error) {
	// 获取参数
	operation, _ := execution.Parameters["operation"].(string)
	projectPath, _ := execution.Parameters["project_path"].(string)
	projectType, _ := execution.Parameters["project_type"].(string)
	projectName, _ := execution.Parameters["project_name"].(string)
	packageManager, _ := execution.Parameters["package_manager"].(string)
	filePath, _ := execution.Parameters["file_path"].(string)
	fileContent, _ := execution.Parameters["file_content"].(string)
	command, _ := execution.Parameters["command"].(string)
	args, _ := execution.Parameters["args"].([]interface{})

	// 验证项目路径安全性
	safeProjectPath, err := t.validateProjectPath(projectPath)
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
		pt := ProjectType(projectType)
		result, operationErr = t.initProject(safeProjectPath, pt, projectName)
	case "install_deps":
		result, operationErr = t.installDependencies(safeProjectPath, packageManager)
	case "create_file":
		result, operationErr = t.createFile(safeProjectPath, filePath, fileContent)
	case "run_command":
		cmdArgs := make([]string, len(args))
		for i, arg := range args {
			cmdArgs[i] = arg.(string)
		}
		result, operationErr = t.runProjectCommand(safeProjectPath, command, cmdArgs)
	case "get_info":
		result, operationErr = t.getProjectInfo(safeProjectPath)
	default:
		operationErr = fmt.Errorf("unsupported project operation: %s", operation)
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

// validateProjectPath 验证项目路径安全性
func (t *ProjectTool) validateProjectPath(projectPath string) (string, error) {
	if projectPath == "" {
		return "", fmt.Errorf("project path cannot be empty")
	}

	// 清理路径
	cleanPath := filepath.Clean(projectPath)

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

// initProject 初始化项目
func (t *ProjectTool) initProject(projectPath string, projectType ProjectType, projectName string) (interface{}, error) {
	// 检查项目目录是否已存在
	if _, err := os.Stat(projectPath); err == nil {
		// 目录存在，检查是否为空
		entries, err := os.ReadDir(projectPath)
		if err != nil {
			return nil, fmt.Errorf("failed to check directory: %w", err)
		}
		if len(entries) > 0 {
			return nil, fmt.Errorf("project directory is not empty")
		}
	}

	// 获取模板
	template, exists := t.templates[projectType]
	if !exists {
		return nil, fmt.Errorf("unsupported project type: %s", projectType)
	}

	// 如果未提供项目名称，使用目录名
	if projectName == "" {
		projectName = filepath.Base(projectPath)
	}

	// 创建项目目录
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create project directory: %w", err)
	}

	// 应用模板
	createdFiles := []string{}
	for filePath, contentTemplate := range template.Structure {
		// 替换模板变量
		content := strings.ReplaceAll(contentTemplate, "{{PROJECT_NAME}}", projectName)
		content = strings.ReplaceAll(content, "{{PROJECT_PATH}}", projectPath)

		// 创建完整文件路径
		fullPath := filepath.Join(projectPath, filePath)

		// 确保目录存在
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory for %s: %w", filePath, err)
		}

		// 写入文件
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("failed to write file %s: %w", filePath, err)
		}

		createdFiles = append(createdFiles, filePath)
	}

	// 初始化版本控制（可选）
	gitPath := filepath.Join(projectPath, ".git")
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		// 可以在这里初始化Git仓库，但为了简化，暂时不自动初始化
	}

	return map[string]interface{}{
		"project_path":  projectPath,
		"project_type":  string(projectType),
		"project_name":  projectName,
		"created_files": createdFiles,
		"template_used": template.Name,
		"initialized":   true,
	}, nil
}

// installDependencies 安装依赖
func (t *ProjectTool) installDependencies(projectPath, packageManager string) (interface{}, error) {
	// 检查项目目录是否存在
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("project directory does not exist: %s", projectPath)
	}

	// 确定项目类型
	projectType, err := t.detectProjectType(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to detect project type: %w", err)
	}

	// 获取模板
	template, exists := t.templates[projectType]
	if !exists {
		return nil, fmt.Errorf("unsupported project type: %s", projectType)
	}

	// 获取依赖列表
	dependencies, exists := template.Dependencies[packageManager]
	if !exists {
		return nil, fmt.Errorf("package manager %s not supported for project type %s", packageManager, projectType)
	}

	// 执行安装命令
	var installCmd string
	var installArgs []string

	switch packageManager {
	case "npm":
		installCmd = "npm"
		installArgs = []string{"install"}
		if len(dependencies) > 0 {
			installArgs = append(installArgs, dependencies...)
		}
	case "yarn":
		installCmd = "yarn"
		installArgs = []string{"add"}
		if len(dependencies) > 0 {
			installArgs = append(installArgs, dependencies...)
		} else {
			installArgs = []string{"install"}
		}
	case "pnpm":
		installCmd = "pnpm"
		installArgs = []string{"add"}
		if len(dependencies) > 0 {
			installArgs = append(installArgs, dependencies...)
		} else {
			installArgs = []string{"install"}
		}
	case "go":
		installCmd = "go"
		installArgs = []string{"mod", "tidy"}
	case "pip", "pip3":
		installCmd = packageManager
		installArgs = []string{"install"}
		if len(dependencies) > 0 {
			installArgs = append(installArgs, dependencies...)
		}
	default:
		return nil, fmt.Errorf("unsupported package manager: %s", packageManager)
	}

	// 这里可以调用Shell工具来执行命令
	// 为了简化，暂时只返回计划执行的命令
	return map[string]interface{}{
		"project_path":    projectPath,
		"project_type":    string(projectType),
		"package_manager": packageManager,
		"dependencies":    dependencies,
		"install_command": installCmd,
		"install_args":    installArgs,
		"install_planned": true,
		"note":            "实际安装需要调用Shell工具执行",
	}, nil
}

// createFile 创建文件
func (t *ProjectTool) createFile(projectPath, filePath, content string) (interface{}, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path is required")
	}

	if content == "" {
		return nil, fmt.Errorf("file content is required")
	}

	// 检查项目目录是否存在
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("project directory does not exist: %s", projectPath)
	}

	// 创建完整文件路径
	fullPath := filepath.Join(projectPath, filePath)

	// 确保目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	// 获取文件信息
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return map[string]interface{}{
		"project_path": projectPath,
		"file_path":    filePath,
		"full_path":    fullPath,
		"size":         info.Size(),
		"created":      true,
	}, nil
}

// runProjectCommand 运行项目命令
func (t *ProjectTool) runProjectCommand(projectPath, command string, args []string) (interface{}, error) {
	if command == "" {
		return nil, fmt.Errorf("command is required")
	}

	// 检查项目目录是否存在
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("project directory does not exist: %s", projectPath)
	}

	// 这里可以调用Shell工具来执行命令
	// 为了简化，暂时只返回计划执行的命令
	return map[string]interface{}{
		"project_path":      projectPath,
		"command":           command,
		"args":              args,
		"execution_planned": true,
		"note":              "实际执行需要调用Shell工具",
	}, nil
}

// getProjectInfo 获取项目信息
func (t *ProjectTool) getProjectInfo(projectPath string) (interface{}, error) {
	// 检查项目目录是否存在
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("project directory does not exist: %s", projectPath)
	}

	// 检测项目类型
	projectType, err := t.detectProjectType(projectPath)
	if err != nil {
		projectType = ProjectTypeStatic
	}

	// 获取目录结构
	files, err := t.scanProjectStructure(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to scan project structure: %w", err)
	}

	// 检查配置文件
	configFiles := t.detectConfigFiles(projectPath)

	// 获取项目大小
	size, err := t.calculateProjectSize(projectPath)
	if err != nil {
		size = -1
	}

	return map[string]interface{}{
		"project_path": projectPath,
		"project_type": string(projectType),
		"exists":       true,
		"file_count":   len(files),
		"total_size":   size,
		"config_files": configFiles,
		"structure":    files,
	}, nil
}

// detectProjectType 检测项目类型
func (t *ProjectTool) detectProjectType(projectPath string) (ProjectType, error) {
	// 检查Go项目
	goModPath := filepath.Join(projectPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		return ProjectTypeGo, nil
	}

	// 检查Node.js项目
	packageJsonPath := filepath.Join(projectPath, "package.json")
	if _, err := os.Stat(packageJsonPath); err == nil {
		// 检查是否是React项目
		reactFiles := []string{"src/App.js", "src/App.jsx", "src/App.tsx", "src/index.js"}
		for _, file := range reactFiles {
			if _, err := os.Stat(filepath.Join(projectPath, file)); err == nil {
				return ProjectTypeReact, nil
			}
		}
		return ProjectTypeNodeJS, nil
	}

	// 检查Python项目
	pyFiles := []string{"requirements.txt", "setup.py", "pyproject.toml"}
	for _, file := range pyFiles {
		if _, err := os.Stat(filepath.Join(projectPath, file)); err == nil {
			return ProjectTypePython, nil
		}
	}

	// 检查Vue项目
	vueConfigPath := filepath.Join(projectPath, "vue.config.js")
	if _, err := os.Stat(vueConfigPath); err == nil {
		return ProjectTypeVue, nil
	}

	// 默认返回静态项目
	return ProjectTypeStatic, nil
}

// scanProjectStructure 扫描项目结构
func (t *ProjectTool) scanProjectStructure(projectPath string) ([]map[string]interface{}, error) {
	var files []map[string]interface{}

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(projectPath, path)
		if err != nil {
			return err
		}

		// 跳过项目根目录
		if relPath == "." {
			return nil
		}

		// 跳过.git目录
		if strings.Contains(relPath, ".git") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 跳过node_modules目录
		if strings.Contains(relPath, "node_modules") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		files = append(files, map[string]interface{}{
			"path":     relPath,
			"name":     info.Name(),
			"size":     info.Size(),
			"is_dir":   info.IsDir(),
			"modified": info.ModTime().Unix(),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// detectConfigFiles 检测配置文件
func (t *ProjectTool) detectConfigFiles(projectPath string) map[string]bool {
	configFiles := map[string]bool{}

	commonConfigs := []string{
		"go.mod", "package.json", "requirements.txt", "pyproject.toml",
		"dockerfile", "docker-compose.yml", ".gitignore", ".env",
		"README.md", "LICENSE", "Makefile", "CMakeLists.txt",
	}

	for _, config := range commonConfigs {
		configPath := filepath.Join(projectPath, config)
		if _, err := os.Stat(configPath); err == nil {
			configFiles[config] = true
		}
	}

	return configFiles
}

// calculateProjectSize 计算项目大小
func (t *ProjectTool) calculateProjectSize(projectPath string) (int64, error) {
	var totalSize int64

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过.git目录
		if strings.Contains(path, ".git") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 跳过node_modules目录
		if strings.Contains(path, "node_modules") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			totalSize += info.Size()
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return totalSize, nil
}

// loadDefaultTemplates 加载默认模板
func loadDefaultTemplates() map[ProjectType]*ProjectTemplate {
	templates := make(map[ProjectType]*ProjectTemplate)

	// Go项目模板
	templates[ProjectTypeGo] = &ProjectTemplate{
		Type:        ProjectTypeGo,
		Name:        "Go Project",
		Description: "标准的Go语言项目模板",
		Structure: map[string]string{
			"go.mod": `module {{PROJECT_NAME}}

go 1.24.1`,
			"main.go": `package main

import "fmt"

func main() {
	fmt.Println("Hello, {{PROJECT_NAME}}!")
}`,
			"README.md": "# {{PROJECT_NAME}}\n\n这是一个Go语言项目。\n\n## 构建和运行\n\n```bash\ngo build\n./{{PROJECT_NAME}}\n```",
			".gitignore": `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with go test -c
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/`,
		},
		Dependencies: map[string][]string{
			"go": {},
		},
		Commands: map[string][]string{
			"build": {"go build"},
			"run":   {"go run main.go"},
			"test":  {"go test ./..."},
		},
	}

	// Node.js项目模板
	templates[ProjectTypeNodeJS] = &ProjectTemplate{
		Type:        ProjectTypeNodeJS,
		Name:        "Node.js Project",
		Description: "标准的Node.js项目模板",
		Structure: map[string]string{
			"package.json": `{
  "name": "{{PROJECT_NAME}}",
  "version": "1.0.0",
  "description": "A Node.js project",
  "main": "index.js",
  "scripts": {
    "start": "node index.js",
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "keywords": [],
  "author": "",
  "license": "ISC"
}`,
			"index.js":  "console.log('Hello, {{PROJECT_NAME}}!');",
			"README.md": "# {{PROJECT_NAME}}\n\n这是一个Node.js项目。\n\n## 安装和运行\n\n```bash\nnpm install\nnpm start\n```",
			".gitignore": `# Logs
logs
*.log
npm-debug.log*
yarn-debug.log*
yarn-error.log*
lerna-debug.log*

# Runtime data
pids
*.pid
*.seed
*.pid.lock

# Directory for instrumented libs generated by jscoverage/JSCover
lib-cov

# Coverage directory used by tools like istanbul
coverage
*.lcov

# nyc test coverage
.nyc_output

# Grunt intermediate storage (https://gruntjs.com/creating-plugins#storing-task-files)
.grunt

# Bower dependency directory (https://bower.io/)
bower_components

# node-waf configuration
.lock-wscript

# Compiled binary addons (https://nodejs.org/api/addons.html)
build/Release

# Dependency directories
node_modules/
jspm_packages/

# TypeScript v1 declaration files
typings/

# TypeScript cache
*.tsbuildinfo

# Optional npm cache directory
.npm

# Optional eslint cache
.eslintcache

# Microbundle cache
.rpt2_cache/
.rts2_cache_cjs/
.rts2_cache_es/
.rts2_cache_umd/

# Optional REPL history
.node_repl_history

# Output of 'npm pack'
*.tgz

# Yarn Integrity file
.yarn-integrity

# dotenv environment variables file
.env
.env.test

# parcel-bundler cache (https://parceljs.org/)
.cache
.parcel-cache

# Next.js build output
.next

# Nuxt.js build / generate output
.nuxt
dist

# Gatsby files
.cache/
# Comment in the public folder in case your project uses Gatsby and not Next.js
# public

# Storybook build outputs
.out
.storybook-out

# vuepress build output
.vuepress/dist

# Serverless directories
.serverless/

# FuseBox cache
.fusebox/

# DynamoDB Local files
.dynamodb/

# TernJS port file
.tern-port

# Stores VSCode versions used for testing VSCode extensions
.vscode-test

# yarn v2
.yarn/cache
.yarn/unplugged
.yarn/build-state.yml
.yarn/install-state.gz
.pnp.*`,
		},
		Dependencies: map[string][]string{
			"npm":  {},
			"yarn": {},
			"pnpm": {},
		},
		Commands: map[string][]string{
			"start": {"npm start"},
			"dev":   {"npm run dev"},
			"build": {"npm run build"},
		},
	}

	// React项目模板（简化版）
	templates[ProjectTypeReact] = &ProjectTemplate{
		Type:        ProjectTypeReact,
		Name:        "React Project",
		Description: "React项目模板",
		Structure: map[string]string{
			"package.json": `{
  "name": "{{PROJECT_NAME}}",
  "version": "1.0.0",
  "private": true,
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  },
  "scripts": {
    "start": "react-scripts start",
    "build": "react-scripts build",
    "test": "react-scripts test",
    "eject": "react-scripts eject"
  },
  "devDependencies": {
    "react-scripts": "5.0.1"
  },
  "browserslist": {
    "production": [
      ">0.2%",
      "not dead",
      "not op_mini all"
    ],
    "development": [
      "last 1 chrome version",
      "last 1 firefox version",
      "last 1 safari version"
    ]
  }
}`,
			"src/App.js": `import React from 'react';

function App() {
  return (
    <div className="App">
      <h1>Hello, {{PROJECT_NAME}}!</h1>
      <p>Welcome to your React app.</p>
    </div>
  );
}

export default App;`,
			"src/index.js": `import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);`,
			"public/index.html": `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <link rel="icon" href="%PUBLIC_URL%/favicon.ico" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="theme-color" content="#000000" />
    <meta
      name="description"
      content="Web site created using create-react-app"
    />
    <title>{{PROJECT_NAME}}</title>
  </head>
  <body>
    <noscript>You need to enable JavaScript to run this app.</noscript>
    <div id="root"></div>
  </body>
</html>`,
		},
		Dependencies: map[string][]string{
			"npm":  {"react", "react-dom", "react-scripts"},
			"yarn": {"react", "react-dom", "react-scripts"},
		},
		Commands: map[string][]string{
			"start": {"npm start"},
			"build": {"npm run build"},
			"test":  {"npm test"},
		},
	}

	// 静态网站模板
	templates[ProjectTypeStatic] = &ProjectTemplate{
		Type:        ProjectTypeStatic,
		Name:        "Static Website",
		Description: "静态网站模板",
		Structure: map[string]string{
			"index.html": `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{PROJECT_NAME}}</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            line-height: 1.6;
        }
        h1 {
            color: #333;
        }
    </style>
</head>
<body>
    <h1>Welcome to {{PROJECT_NAME}}</h1>
    <p>This is a static website.</p>
</body>
</html>`,
			"README.md": "# {{PROJECT_NAME}}\n\n这是一个静态网站。\n\n## 使用\n\n直接在浏览器中打开index.html文件即可。",
			".gitignore": `# Logs
logs
*.log

# Runtime data
pids
*.pid
*.seed

# Directory for instrumented libs generated by jscoverage/JSCover
lib-cov

# Coverage directory used by tools like istanbul
coverage

# nyc test coverage
.nyc_output

# Grunt intermediate storage (https://gruntjs.com/creating-plugins#storing-task-files)
.grunt

# node-waf configuration
.lock-wscript

# Compiled binary addons (https://nodejs.org/api/addons.html)
build/Release

# Dependency directories
node_modules/
jspm_packages/

# Optional npm cache directory
.npm

# Optional eslint cache
.eslintcache

# Optional REPL history
.node_repl_history

# Output of 'npm pack'
*.tgz

# Yarn Integrity file
.yarn-integrity

# dotenv environment variables file
.env
.env.test

# parcel-bundler cache (https://parceljs.org/)
.cache
.parcel-cache

# Next.js build output
.next

# Nuxt.js build / generate output
.nuxt
dist

# Gatsby files
.cache/
public

# vuepress build output
.vuepress/dist

# Serverless directories
.serverless/

# FuseBox cache
.fusebox/

# DynamoDB Local files
.dynamodb/

# TernJS port file
.tern-port

# Stores VSCode versions used for testing VSCode extensions
.vscode-test

# yarn v2
.yarn/cache
.yarn/unplugged
.yarn/build-state.yml
.yarn/install-state.gz
.pnp.*`,
		},
		Dependencies: map[string][]string{},
		Commands:     map[string][]string{},
	}

	return templates
}

// GetTemplate 获取项目模板
func (t *ProjectTool) GetTemplate(projectType ProjectType) (*ProjectTemplate, error) {
	template, exists := t.templates[projectType]
	if !exists {
		return nil, fmt.Errorf("template not found for project type: %s", projectType)
	}
	return template, nil
}

// AddTemplate 添加项目模板
func (t *ProjectTool) AddTemplate(template *ProjectTemplate) {
	t.templates[template.Type] = template
}

// RemoveTemplate 移除项目模板
func (t *ProjectTool) RemoveTemplate(projectType ProjectType) {
	delete(t.templates, projectType)
}

// ListTemplates 列出所有模板
func (t *ProjectTool) ListTemplates() []map[string]interface{} {
	var templates []map[string]interface{}
	for _, template := range t.templates {
		templates = append(templates, map[string]interface{}{
			"type":        string(template.Type),
			"name":        template.Name,
			"description": template.Description,
		})
	}
	return templates
}

// Initialize 初始化工具
func (t *ProjectTool) Initialize() error {
	// 调用父类初始化
	if err := t.BaseTool.Initialize(); err != nil {
		return err
	}

	utils.Info("项目管理工具初始化完成，工作区: %s", t.workspacePath)
	utils.Info("可用模板: %v", t.ListTemplates())
	return nil
}
