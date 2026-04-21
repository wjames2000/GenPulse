package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// SemanticMemory 语义记忆（L3）- 文件系统存储，构建用户画像
type SemanticMemory struct {
	baseDir    string
	memoryFile string
	userFile   string
	mu         sync.RWMutex
}

// UserProfile 用户画像
type UserProfile struct {
	UserID          string             `json:"user_id"`
	Username        string             `json:"username,omitempty"`
	Preferences     map[string]any     `json:"preferences,omitempty"`
	Skills          []string           `json:"skills,omitempty"`
	Interests       []string           `json:"interests,omitempty"`
	Goals           []string           `json:"goals,omitempty"`
	WorkingStyle    map[string]any     `json:"working_style,omitempty"`
	Communication   map[string]any     `json:"communication,omitempty"`
	KnowledgeAreas  map[string]int     `json:"knowledge_areas,omitempty"`  // 领域知识评分 0-100
	TaskPreferences map[string]float64 `json:"task_preferences,omitempty"` // 任务类型偏好 0-1
	SuccessPatterns []SuccessPattern   `json:"success_patterns,omitempty"`
	FailurePatterns []FailurePattern   `json:"failure_patterns,omitempty"`
	LearningHistory []LearningEvent    `json:"learning_history,omitempty"`
	LastUpdated     time.Time          `json:"last_updated"`
	CreatedAt       time.Time          `json:"created_at"`
}

// SuccessPattern 成功模式
type SuccessPattern struct {
	PatternID     string         `json:"pattern_id"`
	TaskType      string         `json:"task_type"`
	Description   string         `json:"description"`
	KeyFactors    []string       `json:"key_factors"`
	Conditions    map[string]any `json:"conditions"`
	Effectiveness float64        `json:"effectiveness"` // 0-1 有效性评分
	UsageCount    int            `json:"usage_count"`
	LastUsed      time.Time      `json:"last_used"`
	CreatedAt     time.Time      `json:"created_at"`
}

// FailurePattern 失败模式
type FailurePattern struct {
	PatternID    string         `json:"pattern_id"`
	TaskType     string         `json:"task_type"`
	Description  string         `json:"description"`
	ErrorType    string         `json:"error_type"`
	RootCause    string         `json:"root_cause"`
	Conditions   map[string]any `json:"conditions"`
	Frequency    int            `json:"frequency"`
	LastOccurred time.Time      `json:"last_occurred"`
	CreatedAt    time.Time      `json:"created_at"`
}

// LearningEvent 学习事件
type LearningEvent struct {
	EventID      string    `json:"event_id"`
	TaskType     string    `json:"task_type"`
	Description  string    `json:"description"`
	Insight      string    `json:"insight"`
	Impact       string    `json:"impact"`     // positive, negative, neutral
	Confidence   float64   `json:"confidence"` // 0-1 置信度
	Tags         []string  `json:"tags,omitempty"`
	RelatedTasks []string  `json:"related_tasks,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// SemanticQuery 语义查询
type SemanticQuery struct {
	Query         string  `json:"query,omitempty"`
	TaskType      string  `json:"task_type,omitempty"`
	Category      string  `json:"category,omitempty"`
	MinConfidence float64 `json:"min_confidence,omitempty"`
	Limit         int     `json:"limit,omitempty"`
}

// NewSemanticMemory 创建语义记忆
func NewSemanticMemory(baseDir string) (*SemanticMemory, error) {
	// 确保目录存在
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	sm := &SemanticMemory{
		baseDir:    baseDir,
		memoryFile: filepath.Join(baseDir, "MEMORY.md"),
		userFile:   filepath.Join(baseDir, "USER.md"),
	}

	// 初始化文件
	if err := sm.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize semantic memory: %w", err)
	}

	return sm, nil
}

// initialize 初始化语义记忆文件
func (sm *SemanticMemory) initialize() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 检查并创建MEMORY.md
	if _, err := os.Stat(sm.memoryFile); os.IsNotExist(err) {
		defaultMemory := `# 语义记忆库

## 概述
这是AI助手的语义记忆库，存储长期知识和经验。

## 知识领域
- 编程开发
- 系统设计
- 问题解决
- 代码审查
- 文档编写

## 成功经验
- 保持代码简洁和可维护性
- 优先考虑用户体验
- 遵循最佳实践和设计模式

## 学习要点
- 从错误中学习并改进
- 持续更新知识和技能
- 适应不同的工作风格

## 更新记录
- 初始创建: ` + time.Now().Format("2006-01-02 15:04:05") + `
`
		if err := os.WriteFile(sm.memoryFile, []byte(defaultMemory), 0644); err != nil {
			return fmt.Errorf("failed to create MEMORY.md: %w", err)
		}
	}

	// 检查并创建USER.md
	if _, err := os.Stat(sm.userFile); os.IsNotExist(err) {
		defaultUser := `# 用户画像

## 基本信息
- 用户ID: default
- 创建时间: ` + time.Now().Format("2006-01-02 15:04:05") + `

## 偏好设置
- 工作风格: 注重细节，喜欢结构化方法
- 沟通方式: 直接、清晰、有条理
- 学习偏好: 通过实践学习，喜欢示例代码

## 技能兴趣
- 编程语言: Go, Python, JavaScript
- 技术领域: 后端开发，系统架构，数据库设计
- 工具使用: Git, Docker, Kubernetes

## 目标
- 提高代码质量
- 学习新技术
- 构建可扩展的系统

## 知识领域评分
- 编程开发: 85
- 系统设计: 80
- 问题解决: 90
- 代码审查: 75
- 文档编写: 70

## 任务类型偏好
- 代码实现: 0.9
- 问题调试: 0.8
- 系统设计: 0.7
- 文档编写: 0.6
- 代码审查: 0.8
`
		if err := os.WriteFile(sm.userFile, []byte(defaultUser), 0644); err != nil {
			return fmt.Errorf("failed to create USER.md: %w", err)
		}
	}

	return nil
}

// GetUserProfile 获取用户画像
func (sm *SemanticMemory) GetUserProfile() (*UserProfile, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// 读取USER.md文件
	content, err := os.ReadFile(sm.userFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read USER.md: %w", err)
	}

	// 解析Markdown文件
	return sm.parseUserProfile(string(content))
}

// UpdateUserProfile 更新用户画像
func (sm *SemanticMemory) UpdateUserProfile(profile *UserProfile) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 更新时间和ID
	now := time.Now()
	if profile.CreatedAt.IsZero() {
		profile.CreatedAt = now
	}
	profile.LastUpdated = now

	// 生成Markdown内容
	content := sm.generateUserProfileMarkdown(profile)

	// 写入文件
	if err := os.WriteFile(sm.userFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write USER.md: %w", err)
	}

	return nil
}

// GetMemoryContent 获取记忆内容
func (sm *SemanticMemory) GetMemoryContent() (string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	content, err := os.ReadFile(sm.memoryFile)
	if err != nil {
		return "", fmt.Errorf("failed to read MEMORY.md: %w", err)
	}

	return string(content), nil
}

// UpdateMemoryContent 更新记忆内容
func (sm *SemanticMemory) UpdateMemoryContent(content string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 添加更新时间戳
	updatedContent := content + "\n\n## 更新记录\n- 最后更新: " + time.Now().Format("2006-01-02 15:04:05")

	if err := os.WriteFile(sm.memoryFile, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write MEMORY.md: %w", err)
	}

	return nil
}

// AppendMemory 追加记忆内容
func (sm *SemanticMemory) AppendMemory(section, content string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 读取现有内容
	existingContent, err := os.ReadFile(sm.memoryFile)
	if err != nil {
		return fmt.Errorf("failed to read MEMORY.md: %w", err)
	}

	// 解析并追加内容
	updatedContent := sm.appendToSection(string(existingContent), section, content)

	// 写入文件
	if err := os.WriteFile(sm.memoryFile, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write MEMORY.md: %w", err)
	}

	return nil
}

// Search 语义搜索
func (sm *SemanticMemory) Search(query *SemanticQuery) ([]string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// 读取记忆内容
	content, err := os.ReadFile(sm.memoryFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read MEMORY.md: %w", err)
	}

	// 简单关键词搜索
	lines := strings.Split(string(content), "\n")
	var results []string

	for _, line := range lines {
		if query.Query != "" && strings.Contains(strings.ToLower(line), strings.ToLower(query.Query)) {
			results = append(results, line)
		}
	}

	// 限制结果数量
	if query.Limit > 0 && len(results) > query.Limit {
		results = results[:query.Limit]
	}

	return results, nil
}

// AddSuccessPattern 添加成功模式
func (sm *SemanticMemory) AddSuccessPattern(pattern *SuccessPattern) error {
	profile, err := sm.GetUserProfile()
	if err != nil {
		return err
	}

	// 设置默认值
	if pattern.PatternID == "" {
		pattern.PatternID = fmt.Sprintf("success_%d", time.Now().UnixNano())
	}
	if pattern.CreatedAt.IsZero() {
		pattern.CreatedAt = time.Now()
	}
	pattern.LastUsed = time.Now()

	// 添加到用户画像
	profile.SuccessPatterns = append(profile.SuccessPatterns, *pattern)

	// 更新任务类型偏好
	if _, exists := profile.TaskPreferences[pattern.TaskType]; !exists {
		profile.TaskPreferences[pattern.TaskType] = 0.7
	} else {
		// 根据成功模式提高偏好
		profile.TaskPreferences[pattern.TaskType] = min(1.0, profile.TaskPreferences[pattern.TaskType]+0.05)
	}

	// 更新用户画像
	return sm.UpdateUserProfile(profile)
}

// AddFailurePattern 添加失败模式
func (sm *SemanticMemory) AddFailurePattern(pattern *FailurePattern) error {
	profile, err := sm.GetUserProfile()
	if err != nil {
		return err
	}

	// 设置默认值
	if pattern.PatternID == "" {
		pattern.PatternID = fmt.Sprintf("failure_%d", time.Now().UnixNano())
	}
	if pattern.CreatedAt.IsZero() {
		pattern.CreatedAt = time.Now()
	}
	pattern.LastOccurred = time.Now()

	// 添加到用户画像
	profile.FailurePatterns = append(profile.FailurePatterns, *pattern)

	// 更新用户画像
	return sm.UpdateUserProfile(profile)
}

// AddLearningEvent 添加学习事件
func (sm *SemanticMemory) AddLearningEvent(event *LearningEvent) error {
	profile, err := sm.GetUserProfile()
	if err != nil {
		return err
	}

	// 设置默认值
	if event.EventID == "" {
		event.EventID = fmt.Sprintf("learn_%d", time.Now().UnixNano())
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}

	// 添加到用户画像
	profile.LearningHistory = append(profile.LearningHistory, *event)

	// 更新知识领域评分
	for _, tag := range event.Tags {
		if currentScore, exists := profile.KnowledgeAreas[tag]; exists {
			if event.Impact == "positive" {
				newScore := currentScore + 5
				if newScore > 100 {
					newScore = 100
				}
				profile.KnowledgeAreas[tag] = newScore
			} else if event.Impact == "negative" {
				newScore := currentScore - 3
				if newScore < 0 {
					newScore = 0
				}
				profile.KnowledgeAreas[tag] = newScore
			}
		} else {
			profile.KnowledgeAreas[tag] = 50 // 初始评分
		}
	}

	// 更新用户画像
	return sm.UpdateUserProfile(profile)
}

// GetTaskAdvice 获取任务建议
func (sm *SemanticMemory) GetTaskAdvice(taskType string) (string, error) {
	profile, err := sm.GetUserProfile()
	if err != nil {
		return "", err
	}

	var advice strings.Builder
	advice.WriteString(fmt.Sprintf("## %s任务建议\n\n", taskType))

	// 添加成功模式建议
	successCount := 0
	for _, pattern := range profile.SuccessPatterns {
		if pattern.TaskType == taskType && pattern.Effectiveness > 0.7 {
			advice.WriteString(fmt.Sprintf("### 成功模式: %s\n", pattern.Description))
			advice.WriteString(fmt.Sprintf("- 关键因素: %s\n", strings.Join(pattern.KeyFactors, ", ")))
			advice.WriteString(fmt.Sprintf("- 有效性: %.0f%%\n", pattern.Effectiveness*100))
			advice.WriteString(fmt.Sprintf("- 使用次数: %d\n\n", pattern.UsageCount))
			successCount++

			if successCount >= 3 {
				break
			}
		}
	}

	// 添加失败模式警告
	failureCount := 0
	for _, pattern := range profile.FailurePatterns {
		if pattern.TaskType == taskType && pattern.Frequency > 1 {
			advice.WriteString(fmt.Sprintf("### 注意避免: %s\n", pattern.Description))
			advice.WriteString(fmt.Sprintf("- 错误类型: %s\n", pattern.ErrorType))
			advice.WriteString(fmt.Sprintf("- 根本原因: %s\n", pattern.RootCause))
			advice.WriteString(fmt.Sprintf("- 发生频率: %d次\n\n", pattern.Frequency))
			failureCount++

			if failureCount >= 2 {
				break
			}
		}
	}

	// 添加学习要点
	learningCount := 0
	for _, event := range profile.LearningHistory {
		if event.TaskType == taskType && event.Confidence > 0.7 {
			advice.WriteString(fmt.Sprintf("### 学习要点: %s\n", event.Insight))
			advice.WriteString(fmt.Sprintf("- 影响: %s\n", event.Impact))
			advice.WriteString(fmt.Sprintf("- 置信度: %.0f%%\n\n", event.Confidence*100))
			learningCount++

			if learningCount >= 2 {
				break
			}
		}
	}

	return advice.String(), nil
}

// parseUserProfile 解析用户画像Markdown
func (sm *SemanticMemory) parseUserProfile(content string) (*UserProfile, error) {
	profile := &UserProfile{
		UserID:          "default",
		Preferences:     make(map[string]any),
		Skills:          []string{},
		Interests:       []string{},
		Goals:           []string{},
		WorkingStyle:    make(map[string]any),
		Communication:   make(map[string]any),
		KnowledgeAreas:  make(map[string]int),
		TaskPreferences: make(map[string]float64),
		SuccessPatterns: []SuccessPattern{},
		FailurePatterns: []FailurePattern{},
		LearningHistory: []LearningEvent{},
		CreatedAt:       time.Now(),
		LastUpdated:     time.Now(),
	}

	lines := strings.Split(content, "\n")
	currentSection := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "## ") {
			currentSection = strings.TrimPrefix(line, "## ")
			continue
		}

		if strings.HasPrefix(line, "- ") {
			item := strings.TrimPrefix(line, "- ")

			switch currentSection {
			case "基本信息":
				if strings.Contains(item, "用户ID:") {
					profile.UserID = strings.TrimSpace(strings.Split(item, ":")[1])
				} else if strings.Contains(item, "创建时间:") {
					// 解析创建时间
				}
			case "偏好设置":
				if strings.Contains(item, ":") {
					parts := strings.SplitN(item, ":", 2)
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					profile.Preferences[key] = value
				}
			case "技能兴趣":
				if strings.Contains(item, ":") {
					parts := strings.SplitN(item, ":", 2)
					category := strings.TrimSpace(parts[0])
					items := strings.Split(strings.TrimSpace(parts[1]), ",")

					for _, skill := range items {
						skill = strings.TrimSpace(skill)
						switch category {
						case "编程语言":
							profile.Skills = append(profile.Skills, skill)
						case "技术领域":
							profile.Interests = append(profile.Interests, skill)
						}
					}
				} else {
					profile.Skills = append(profile.Skills, item)
				}
			case "目标":
				profile.Goals = append(profile.Goals, item)
			case "知识领域评分":
				if strings.Contains(item, ":") {
					parts := strings.SplitN(item, ":", 2)
					area := strings.TrimSpace(parts[0])
					score := 0
					fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &score)
					profile.KnowledgeAreas[area] = score
				}
			case "任务类型偏好":
				if strings.Contains(item, ":") {
					parts := strings.SplitN(item, ":", 2)
					taskType := strings.TrimSpace(parts[0])
					preference := 0.0
					fmt.Sscanf(strings.TrimSpace(parts[1]), "%f", &preference)
					profile.TaskPreferences[taskType] = preference
				}
			}
		}
	}

	return profile, nil
}

// generateUserProfileMarkdown 生成用户画像Markdown
func (sm *SemanticMemory) generateUserProfileMarkdown(profile *UserProfile) string {
	var content strings.Builder

	content.WriteString("# 用户画像\n\n")

	// 基本信息
	content.WriteString("## 基本信息\n")
	content.WriteString(fmt.Sprintf("- 用户ID: %s\n", profile.UserID))
	if profile.Username != "" {
		content.WriteString(fmt.Sprintf("- 用户名: %s\n", profile.Username))
	}
	content.WriteString(fmt.Sprintf("- 创建时间: %s\n", profile.CreatedAt.Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("- 最后更新: %s\n\n", profile.LastUpdated.Format("2006-01-02 15:04:05")))

	// 偏好设置
	content.WriteString("## 偏好设置\n")
	for key, value := range profile.Preferences {
		content.WriteString(fmt.Sprintf("- %s: %v\n", key, value))
	}
	if len(profile.Preferences) == 0 {
		content.WriteString("- 工作风格: 注重细节，喜欢结构化方法\n")
		content.WriteString("- 沟通方式: 直接、清晰、有条理\n")
		content.WriteString("- 学习偏好: 通过实践学习，喜欢示例代码\n")
	}
	content.WriteString("\n")

	// 技能兴趣
	content.WriteString("## 技能兴趣\n")
	if len(profile.Skills) > 0 {
		content.WriteString("- 编程语言: " + strings.Join(profile.Skills, ", ") + "\n")
	}
	if len(profile.Interests) > 0 {
		content.WriteString("- 技术领域: " + strings.Join(profile.Interests, ", ") + "\n")
	}
	if len(profile.Skills) == 0 && len(profile.Interests) == 0 {
		content.WriteString("- 编程语言: Go, Python, JavaScript\n")
		content.WriteString("- 技术领域: 后端开发，系统架构，数据库设计\n")
		content.WriteString("- 工具使用: Git, Docker, Kubernetes\n")
	}
	content.WriteString("\n")

	// 目标
	content.WriteString("## 目标\n")
	if len(profile.Goals) > 0 {
		for _, goal := range profile.Goals {
			content.WriteString(fmt.Sprintf("- %s\n", goal))
		}
	} else {
		content.WriteString("- 提高代码质量\n")
		content.WriteString("- 学习新技术\n")
		content.WriteString("- 构建可扩展的系统\n")
	}
	content.WriteString("\n")

	// 知识领域评分
	content.WriteString("## 知识领域评分\n")
	if len(profile.KnowledgeAreas) > 0 {
		for area, score := range profile.KnowledgeAreas {
			content.WriteString(fmt.Sprintf("- %s: %d\n", area, score))
		}
	} else {
		content.WriteString("- 编程开发: 85\n")
		content.WriteString("- 系统设计: 80\n")
		content.WriteString("- 问题解决: 90\n")
		content.WriteString("- 代码审查: 75\n")
		content.WriteString("- 文档编写: 70\n")
	}
	content.WriteString("\n")

	// 任务类型偏好
	content.WriteString("## 任务类型偏好\n")
	if len(profile.TaskPreferences) > 0 {
		for taskType, preference := range profile.TaskPreferences {
			content.WriteString(fmt.Sprintf("- %s: %.1f\n", taskType, preference))
		}
	} else {
		content.WriteString("- 代码实现: 0.9\n")
		content.WriteString("- 问题调试: 0.8\n")
		content.WriteString("- 系统设计: 0.7\n")
		content.WriteString("- 文档编写: 0.6\n")
		content.WriteString("- 代码审查: 0.8\n")
	}
	content.WriteString("\n")

	// 成功模式
	if len(profile.SuccessPatterns) > 0 {
		content.WriteString("## 成功模式\n")
		for _, pattern := range profile.SuccessPatterns {
			content.WriteString(fmt.Sprintf("### %s\n", pattern.Description))
			content.WriteString(fmt.Sprintf("- 任务类型: %s\n", pattern.TaskType))
			content.WriteString(fmt.Sprintf("- 关键因素: %s\n", strings.Join(pattern.KeyFactors, ", ")))
			content.WriteString(fmt.Sprintf("- 有效性: %.0f%%\n", pattern.Effectiveness*100))
			content.WriteString(fmt.Sprintf("- 使用次数: %d\n", pattern.UsageCount))
			content.WriteString(fmt.Sprintf("- 最后使用: %s\n\n", pattern.LastUsed.Format("2006-01-02")))
		}
	}

	// 失败模式
	if len(profile.FailurePatterns) > 0 {
		content.WriteString("## 失败模式\n")
		for _, pattern := range profile.FailurePatterns {
			content.WriteString(fmt.Sprintf("### %s\n", pattern.Description))
			content.WriteString(fmt.Sprintf("- 任务类型: %s\n", pattern.TaskType))
			content.WriteString(fmt.Sprintf("- 错误类型: %s\n", pattern.ErrorType))
			content.WriteString(fmt.Sprintf("- 根本原因: %s\n", pattern.RootCause))
			content.WriteString(fmt.Sprintf("- 发生频率: %d次\n", pattern.Frequency))
			content.WriteString(fmt.Sprintf("- 最后发生: %s\n\n", pattern.LastOccurred.Format("2006-01-02")))
		}
	}

	// 学习历史
	if len(profile.LearningHistory) > 0 {
		content.WriteString("## 学习历史\n")
		for _, event := range profile.LearningHistory {
			content.WriteString(fmt.Sprintf("### %s\n", event.Insight))
			content.WriteString(fmt.Sprintf("- 任务类型: %s\n", event.TaskType))
			content.WriteString(fmt.Sprintf("- 描述: %s\n", event.Description))
			content.WriteString(fmt.Sprintf("- 影响: %s\n", event.Impact))
			content.WriteString(fmt.Sprintf("- 置信度: %.0f%%\n", event.Confidence*100))
			if len(event.Tags) > 0 {
				content.WriteString(fmt.Sprintf("- 标签: %s\n", strings.Join(event.Tags, ", ")))
			}
			content.WriteString(fmt.Sprintf("- 时间: %s\n\n", event.CreatedAt.Format("2006-01-02")))
		}
	}

	return content.String()
}

// appendToSection 追加内容到指定章节
func (sm *SemanticMemory) appendToSection(content, section, newContent string) string {
	lines := strings.Split(content, "\n")
	var result []string
	inTargetSection := false
	sectionAdded := false

	for _, line := range lines {
		result = append(result, line)

		if strings.HasPrefix(line, "## "+section) {
			inTargetSection = true
		} else if strings.HasPrefix(line, "## ") && inTargetSection {
			// 遇到下一个章节，插入新内容
			result = append(result, "- "+newContent+" ("+time.Now().Format("2006-01-02")+")")
			result = append(result, "")
			inTargetSection = false
			sectionAdded = true
		}
	}

	// 如果目标章节是最后一个章节
	if inTargetSection && !sectionAdded {
		result = append(result, "- "+newContent+" ("+time.Now().Format("2006-01-02")+")")
		result = append(result, "")
	}

	// 如果章节不存在，添加新章节
	if !sectionAdded {
		result = append(result, "## "+section)
		result = append(result, "- "+newContent+" ("+time.Now().Format("2006-01-02")+")")
		result = append(result, "")
	}

	return strings.Join(result, "\n")
}
