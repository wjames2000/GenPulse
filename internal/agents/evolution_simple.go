package agents

import (
	"fmt"
	"strings"
	"time"
)

// SimpleEvolutionManager 简化的自进化管理器
// 实现WBS文档中3.3节的核心功能
type SimpleEvolutionManager struct {
	// 技能存储（简化版本）
	skills map[string]*SimpleSkill

	// 事件记录
	events []*SimpleEvolutionEvent

	// 配置
	config SimpleEvolutionConfig
}

// SimpleSkill 简化技能结构
type SimpleSkill struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	UsageCount  int       `json:"usage_count"`
	SuccessRate float64   `json:"success_rate"`
	CreatedAt   time.Time `json:"created_at"`
	LastUsed    time.Time `json:"last_used"`
	Enabled     bool      `json:"enabled"`
}

// SimpleEvolutionConfig 简化配置
type SimpleEvolutionConfig struct {
	EnablePromptEnhancement bool `json:"enable_prompt_enhancement"`
	EnableFeedbackLoop      bool `json:"enable_feedback_loop"`
	EnableEventTracking     bool `json:"enable_event_tracking"`
}

// NewSimpleEvolutionManager 创建简化的自进化管理器
func NewSimpleEvolutionManager(config SimpleEvolutionConfig) *SimpleEvolutionManager {
	return &SimpleEvolutionManager{
		skills: make(map[string]*SimpleSkill),
		events: make([]*SimpleEvolutionEvent, 0),
		config: config,
	}
}

// 3.3.1 Agent执行前记忆注入功能
// EnhancePrompt 增强Agent提示词，注入记忆和技能
func (sem *SimpleEvolutionManager) EnhancePrompt(agent Agent, task string) string {
	if !sem.config.EnablePromptEnhancement {
		return task
	}

	// 获取相关技能
	relevantSkills := sem.getRelevantSkills(agent, task)

	// 构建增强提示词
	enhanced := sem.buildEnhancedPrompt(agent, task, relevantSkills)

	// 记录事件
	if sem.config.EnableEventTracking {
		sem.recordEvent("prompt_enhanced", agent.GetConfig().ID, agent.GetConfig().Name,
			fmt.Sprintf("为任务'%s'增强提示词，注入%d个技能", task, len(relevantSkills)), true)
	}

	return enhanced
}

// getRelevantSkills 获取相关技能
func (sem *SimpleEvolutionManager) getRelevantSkills(agent Agent, task string) []*SimpleSkill {
	relevant := make([]*SimpleSkill, 0)

	// 简单匹配：根据任务关键词和Agent角色
	keywords := extractSimpleKeywords(task)
	agentRole := string(agent.GetConfig().Role)

	for _, skill := range sem.skills {
		if !skill.Enabled {
			continue
		}

		// 检查技能是否相关
		if sem.isSkillRelevant(skill, keywords, agentRole) {
			relevant = append(relevant, skill)
		}
	}

	// 按使用次数和成功率排序
	sem.sortSkillsByRelevance(relevant)

	// 限制数量
	if len(relevant) > 3 {
		return relevant[:3]
	}

	return relevant
}

// isSkillRelevant 检查技能是否相关
func (sem *SimpleEvolutionManager) isSkillRelevant(skill *SimpleSkill, keywords []string, agentRole string) bool {
	// 检查技能名称和描述中是否包含关键词
	skillText := strings.ToLower(skill.Name + " " + skill.Description)

	for _, keyword := range keywords {
		if strings.Contains(skillText, strings.ToLower(keyword)) {
			return true
		}
	}

	// 检查技能分类是否匹配Agent角色
	if skill.Category == agentRole {
		return true
	}

	return false
}

// sortSkillsByRelevance 按相关性排序技能
func (sem *SimpleEvolutionManager) sortSkillsByRelevance(skills []*SimpleSkill) {
	// 按使用次数降序，成功率降序排序
	for i := 0; i < len(skills)-1; i++ {
		for j := i + 1; j < len(skills); j++ {
			// 比较使用次数
			if skills[j].UsageCount > skills[i].UsageCount {
				skills[i], skills[j] = skills[j], skills[i]
			} else if skills[j].UsageCount == skills[i].UsageCount {
				// 使用次数相同，比较成功率
				if skills[j].SuccessRate > skills[i].SuccessRate {
					skills[i], skills[j] = skills[j], skills[i]
				}
			}
		}
	}
}

// buildEnhancedPrompt 构建增强提示词
func (sem *SimpleEvolutionManager) buildEnhancedPrompt(agent Agent, task string, skills []*SimpleSkill) string {
	var sb strings.Builder

	// 基础提示词
	sb.WriteString(agent.GetConfig().PromptTemplates["default"])
	sb.WriteString("\n\n")

	// 添加任务
	sb.WriteString("## 当前任务\n")
	sb.WriteString(task)
	sb.WriteString("\n\n")

	// 添加相关技能
	if len(skills) > 0 {
		sb.WriteString("## 可用技能\n")
		sb.WriteString("以下是你可以使用的相关技能：\n")

		for i, skill := range skills {
			sb.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, skill.Name))
			sb.WriteString(fmt.Sprintf("   描述: %s\n", skill.Description))
			sb.WriteString(fmt.Sprintf("   使用次数: %d, 成功率: %.1f%%\n", skill.UsageCount, skill.SuccessRate*100))
			sb.WriteString("\n")
		}

		sb.WriteString("请根据任务需要，参考以上技能经验。\n")
	}

	return sb.String()
}

// 3.3.2 Agent执行后反馈闭环功能
// RecordExecution 记录执行结果并触发反馈闭环
func (sem *SimpleEvolutionManager) RecordExecution(agent Agent, task string, result *AgentResult) error {
	if !sem.config.EnableFeedbackLoop {
		return nil
	}

	// 记录执行结果
	sem.recordExecutionResult(agent, task, result)

	// 检查是否需要创建技能
	if sem.shouldCreateSkill(task, result) {
		sem.createSkillFromExecution(agent, task, result)
	}

	// 记录事件
	if sem.config.EnableEventTracking {
		eventType := "execution_success"
		if !result.Success {
			eventType = "execution_failed"
		}

		sem.recordEvent(eventType, agent.GetConfig().ID, agent.GetConfig().Name,
			fmt.Sprintf("记录任务执行结果: %s", task), result.Success)
	}

	return nil
}

// recordExecutionResult 记录执行结果
func (sem *SimpleEvolutionManager) recordExecutionResult(agent Agent, task string, result *AgentResult) {
	// 简化实现：只更新相关技能的使用统计
	keywords := extractSimpleKeywords(task)

	for _, skill := range sem.skills {
		if sem.isSkillRelevant(skill, keywords, string(agent.GetConfig().Role)) {
			skill.UsageCount++
			skill.LastUsed = time.Now()

			// 更新成功率
			if result.Success {
				// 简化计算：每次成功增加成功率
				skill.SuccessRate = (skill.SuccessRate*float64(skill.UsageCount-1) + 1.0) / float64(skill.UsageCount)
			} else {
				// 失败：降低成功率
				skill.SuccessRate = (skill.SuccessRate * float64(skill.UsageCount-1)) / float64(skill.UsageCount)
			}
		}
	}
}

// shouldCreateSkill 判断是否需要创建技能
func (sem *SimpleEvolutionManager) shouldCreateSkill(task string, result *AgentResult) bool {
	// 条件：任务成功且复杂度足够
	if !result.Success {
		return false
	}

	// 简单判断：任务描述长度和关键词
	wordCount := len(strings.Fields(task))
	hasComplexKeywords := false

	complexKeywords := []string{"实现", "开发", "构建", "创建", "设计", "架构", "集成", "部署"}
	taskLower := strings.ToLower(task)

	for _, keyword := range complexKeywords {
		if strings.Contains(taskLower, strings.ToLower(keyword)) {
			hasComplexKeywords = true
			break
		}
	}

	return wordCount > 15 || hasComplexKeywords
}

// createSkillFromExecution 从执行结果创建技能
func (sem *SimpleEvolutionManager) createSkillFromExecution(agent Agent, task string, result *AgentResult) {
	// 生成技能ID和名称
	skillID := fmt.Sprintf("skill-%s-%d", agent.GetConfig().ID, time.Now().Unix())
	skillName := sem.generateSkillName(task)

	// 创建技能
	skill := &SimpleSkill{
		ID:          skillID,
		Name:        skillName,
		Description: fmt.Sprintf("基于Agent %s执行'%s'任务的经验", agent.GetConfig().Name, task),
		Category:    string(agent.GetConfig().Role),
		UsageCount:  1,
		SuccessRate: 1.0, // 首次创建基于成功任务
		CreatedAt:   time.Now(),
		LastUsed:    time.Now(),
		Enabled:     true,
	}

	// 保存技能
	sem.skills[skillID] = skill

	// 记录事件
	if sem.config.EnableEventTracking {
		sem.recordEvent("skill_created", agent.GetConfig().ID, agent.GetConfig().Name,
			fmt.Sprintf("创建新技能: %s", skillName), true)
	}
}

// generateSkillName 生成技能名称
func (sem *SimpleEvolutionManager) generateSkillName(task string) string {
	// 提取任务中的关键词
	keywords := extractSimpleKeywords(task)

	if len(keywords) > 0 {
		// 使用前2个关键词
		nameLength := 2
		if len(keywords) < 2 {
			nameLength = len(keywords)
		}
		return fmt.Sprintf("skill_%s", strings.Join(keywords[:nameLength], "_"))
	}

	return fmt.Sprintf("skill_%d", time.Now().Unix())
}

// 3.3.3 自进化事件追踪功能
// recordEvent 记录事件
func (sem *SimpleEvolutionManager) recordEvent(eventType, agentID, agentName, description string, success bool) {
	if !sem.config.EnableEventTracking {
		return
	}

	event := &SimpleEvolutionEvent{
		EventType:   eventType,
		AgentID:     agentID,
		AgentName:   agentName,
		Description: description,
		Success:     success,
		Timestamp:   time.Now(),
	}

	sem.events = append(sem.events, event)

	// 限制事件数量（保留最近1000个）
	if len(sem.events) > 1000 {
		sem.events = sem.events[1:]
	}
}

// GetEvents 获取事件列表
func (sem *SimpleEvolutionManager) GetEvents(limit int) []*SimpleEvolutionEvent {
	if limit <= 0 || limit > len(sem.events) {
		limit = len(sem.events)
	}

	// 返回最近的事件
	start := len(sem.events) - limit
	if start < 0 {
		start = 0
	}

	return sem.events[start:]
}

// GetSkills 获取技能列表
func (sem *SimpleEvolutionManager) GetSkills() []*SimpleSkill {
	skills := make([]*SimpleSkill, 0, len(sem.skills))

	for _, skill := range sem.skills {
		skills = append(skills, skill)
	}

	// 按使用次数排序
	sem.sortSkillsByRelevance(skills)

	return skills
}

// GetStats 获取统计信息
func (sem *SimpleEvolutionManager) GetStats() map[string]interface{} {
	stats := map[string]interface{}{
		"total_skills":      len(sem.skills),
		"total_events":      len(sem.events),
		"enabled_skills":    0,
		"total_skill_usage": 0,
		"avg_success_rate":  0.0,
	}

	var totalUsage int
	var totalSuccessRate float64
	var enabledCount int

	for _, skill := range sem.skills {
		if skill.Enabled {
			enabledCount++
		}
		totalUsage += skill.UsageCount
		totalSuccessRate += skill.SuccessRate
	}

	stats["enabled_skills"] = enabledCount
	stats["total_skill_usage"] = totalUsage

	if len(sem.skills) > 0 {
		stats["avg_success_rate"] = totalSuccessRate / float64(len(sem.skills))
	}

	// 事件类型统计
	eventStats := make(map[string]int)
	for _, event := range sem.events {
		eventStats[event.EventType] = eventStats[event.EventType] + 1
	}
	stats["event_stats"] = eventStats

	return stats
}

// AddTestSkill 添加测试技能（用于测试）
func (sem *SimpleEvolutionManager) AddTestSkill(skill *SimpleSkill) {
	sem.skills[skill.ID] = skill
}

// Helper functions
// extractSimpleKeywords 提取简单关键词
func extractSimpleKeywords(text string) []string {
	words := strings.Fields(text)
	keywords := make([]string, 0)

	// 简单过滤：长度大于1且不是常见停用词
	stopWords := map[string]bool{
		"的": true, "了": true, "在": true, "是": true, "我": true,
		"有": true, "和": true, "就": true, "不": true, "人": true,
		"都": true, "一": true, "一个": true, "也": true, "很": true,
		"到": true, "说": true, "要": true, "去": true, "你": true,
		"会": true, "着": true, "没有": true, "看": true, "过": true,
		"来": true, "还": true, "呢": true, "那": true, "现在": true,
	}

	for _, word := range words {
		word = strings.Trim(word, ".,!?;:\"'()[]{}")
		if len(word) > 1 && !stopWords[word] {
			keywords = append(keywords, strings.ToLower(word))
		}
	}

	return keywords
}

// SimpleEvolutionEvent 简化进化事件
type SimpleEvolutionEvent struct {
	EventType   string    `json:"event_type"`
	AgentID     string    `json:"agent_id,omitempty"`
	AgentName   string    `json:"agent_name,omitempty"`
	Description string    `json:"description"`
	Success     bool      `json:"success"`
	Timestamp   time.Time `json:"timestamp"`
}

// DefaultSimpleConfig 默认简化配置
func DefaultSimpleConfig() SimpleEvolutionConfig {
	return SimpleEvolutionConfig{
		EnablePromptEnhancement: true,
		EnableFeedbackLoop:      true,
		EnableEventTracking:     true,
	}
}
