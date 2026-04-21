package skills

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// Registry 技能注册表
type Registry struct {
	skills    map[string]*Skill
	metadata  map[string]*SkillMetadata
	index     *Index
	storage   Storage
	mu        sync.RWMutex
	skillsDir string
}

// Index 技能索引
type Index struct {
	byCategory   map[string][]string
	byTag        map[string][]string
	byAgent      map[string][]string
	byComplexity map[string][]string
}

// Storage 存储接口
type Storage interface {
	Save(skill *Skill) error
	Load(id string) (*Skill, error)
	Delete(id string) error
	List() ([]*SkillMetadata, error)
	Exists(id string) bool
}

// FileStorage 文件系统存储
type FileStorage struct {
	baseDir string
}

// NewRegistry 创建新注册表
func NewRegistry(skillsDir string) (*Registry, error) {
	// 确保目录存在
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create skills directory: %w", err)
	}

	storage := &FileStorage{baseDir: skillsDir}
	registry := &Registry{
		skills:    make(map[string]*Skill),
		metadata:  make(map[string]*SkillMetadata),
		index:     newIndex(),
		storage:   storage,
		skillsDir: skillsDir,
	}

	// 加载现有技能
	if err := registry.loadAllSkills(); err != nil {
		return nil, fmt.Errorf("failed to load existing skills: %w", err)
	}

	return registry, nil
}

// newIndex 创建新索引
func newIndex() *Index {
	return &Index{
		byCategory:   make(map[string][]string),
		byTag:        make(map[string][]string),
		byAgent:      make(map[string][]string),
		byComplexity: make(map[string][]string),
	}
}

// Register 注册技能
func (r *Registry) Register(skill *Skill) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 检查是否已存在
	if _, exists := r.skills[skill.ID]; exists {
		return fmt.Errorf("skill with ID %s already exists", skill.ID)
	}

	// 保存到存储
	if err := r.storage.Save(skill); err != nil {
		return fmt.Errorf("failed to save skill: %w", err)
	}

	// 更新内存缓存
	r.skills[skill.ID] = skill
	r.metadata[skill.ID] = skill.ToMetadata()

	// 更新索引
	r.updateIndex(skill)

	return nil
}

// Get 获取技能
func (r *Registry) Get(id string) (*Skill, error) {
	r.mu.RLock()
	skill, exists := r.skills[id]
	r.mu.RUnlock()

	if exists {
		return skill, nil
	}

	// 从存储加载
	r.mu.Lock()
	defer r.mu.Unlock()

	// 再次检查（防止并发重复加载）
	if skill, exists := r.skills[id]; exists {
		return skill, nil
	}

	skill, err := r.storage.Load(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load skill: %w", err)
	}

	r.skills[id] = skill
	r.metadata[id] = skill.ToMetadata()
	r.updateIndex(skill)

	return skill, nil
}

// GetMetadata 获取技能元数据
func (r *Registry) GetMetadata(id string) (*SkillMetadata, error) {
	r.mu.RLock()
	metadata, exists := r.metadata[id]
	r.mu.RUnlock()

	if exists {
		return metadata, nil
	}

	// 如果元数据不存在，尝试获取完整技能
	skill, err := r.Get(id)
	if err != nil {
		return nil, err
	}

	return skill.ToMetadata(), nil
}

// Update 更新技能
func (r *Registry) Update(skill *Skill) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 检查是否存在
	if _, exists := r.skills[skill.ID]; !exists {
		return fmt.Errorf("skill with ID %s does not exist", skill.ID)
	}

	// 更新修改时间
	skill.UpdatedAt = time.Now()

	// 保存到存储
	if err := r.storage.Save(skill); err != nil {
		return fmt.Errorf("failed to save skill: %w", err)
	}

	// 更新内存缓存
	r.skills[skill.ID] = skill
	r.metadata[skill.ID] = skill.ToMetadata()

	// 重新构建索引
	r.rebuildIndex()

	return nil
}

// Delete 删除技能
func (r *Registry) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 从存储删除
	if err := r.storage.Delete(id); err != nil {
		return fmt.Errorf("failed to delete skill from storage: %w", err)
	}

	// 从内存删除
	delete(r.skills, id)
	delete(r.metadata, id)

	// 重新构建索引
	r.rebuildIndex()

	return nil
}

// List 列出所有技能元数据
func (r *Registry) List() ([]*SkillMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metadatas := make([]*SkillMetadata, 0, len(r.metadata))
	for _, metadata := range r.metadata {
		metadatas = append(metadatas, metadata)
	}

	// 按名称排序
	sort.Slice(metadatas, func(i, j int) bool {
		return metadatas[i].Name < metadatas[j].Name
	})

	return metadatas, nil
}

// Search 搜索技能
func (r *Registry) Search(query string, filters map[string]string) ([]*SkillMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]*SkillMetadata, 0)

	// 简单搜索实现：按名称和描述匹配
	for _, metadata := range r.metadata {
		matched := false

		// 文本搜索
		if query != "" {
			if contains(metadata.Name, query) || contains(metadata.Description, query) {
				matched = true
			} else {
				// 检查标签
				for _, tag := range metadata.Tags {
					if contains(tag, query) {
						matched = true
						break
					}
				}
			}
		} else {
			matched = true
		}

		// 应用过滤器
		if matched && filters != nil {
			for key, value := range filters {
				switch key {
				case "category":
					if metadata.Category != value {
						matched = false
					}
				case "complexity":
					if metadata.Complexity != value {
						matched = false
					}
				case "enabled":
					enabled := value == "true"
					if metadata.Enabled != enabled {
						matched = false
					}
				case "validated":
					validated := value == "true"
					if metadata.Validated != validated {
						matched = false
					}
				case "tag":
					hasTag := false
					for _, tag := range metadata.Tags {
						if tag == value {
							hasTag = true
							break
						}
					}
					if !hasTag {
						matched = false
					}
				}
			}
		}

		if matched {
			results = append(results, metadata)
		}
	}

	// 按使用次数排序（最常用的在前）
	sort.Slice(results, func(i, j int) bool {
		return results[i].UsageCount > results[j].UsageCount
	})

	return results, nil
}

// GetByCategory 按分类获取技能
func (r *Registry) GetByCategory(category string) ([]*SkillMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	skillIDs := r.index.byCategory[category]
	results := make([]*SkillMetadata, 0, len(skillIDs))

	for _, id := range skillIDs {
		if metadata, exists := r.metadata[id]; exists {
			results = append(results, metadata)
		}
	}

	return results, nil
}

// GetByTag 按标签获取技能
func (r *Registry) GetByTag(tag string) ([]*SkillMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	skillIDs := r.index.byTag[tag]
	results := make([]*SkillMetadata, 0, len(skillIDs))

	for _, id := range skillIDs {
		if metadata, exists := r.metadata[id]; exists {
			results = append(results, metadata)
		}
	}

	return results, nil
}

// IncrementUsage 增加技能使用计数
func (r *Registry) IncrementUsage(id string, success bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	skill, exists := r.skills[id]
	if !exists {
		return fmt.Errorf("skill with ID %s does not exist", id)
	}

	skill.IncrementUsage(success)

	// 保存更新
	if err := r.storage.Save(skill); err != nil {
		return fmt.Errorf("failed to save skill after usage increment: %w", err)
	}

	// 更新元数据
	r.metadata[id] = skill.ToMetadata()

	return nil
}

// Enable 启用技能
func (r *Registry) Enable(id string) error {
	return r.setEnabled(id, true)
}

// Disable 禁用技能
func (r *Registry) Disable(id string) error {
	return r.setEnabled(id, false)
}

// Validate 验证技能
func (r *Registry) Validate(id string) error {
	return r.setValidated(id, true)
}

// Invalidate 标记技能为未验证
func (r *Registry) Invalidate(id string) error {
	return r.setValidated(id, false)
}

// setEnabled 设置启用状态
func (r *Registry) setEnabled(id string, enabled bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	skill, exists := r.skills[id]
	if !exists {
		return fmt.Errorf("skill with ID %s does not exist", id)
	}

	skill.Enabled = enabled
	skill.UpdatedAt = time.Now()

	// 保存更新
	if err := r.storage.Save(skill); err != nil {
		return fmt.Errorf("failed to save skill: %w", err)
	}

	// 更新元数据
	r.metadata[id] = skill.ToMetadata()

	return nil
}

// setValidated 设置验证状态
func (r *Registry) setValidated(id string, validated bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	skill, exists := r.skills[id]
	if !exists {
		return fmt.Errorf("skill with ID %s does not exist", id)
	}

	skill.Validated = validated
	skill.UpdatedAt = time.Now()

	// 保存更新
	if err := r.storage.Save(skill); err != nil {
		return fmt.Errorf("failed to save skill: %w", err)
	}

	// 更新元数据
	r.metadata[id] = skill.ToMetadata()

	return nil
}

// updateIndex 更新索引
func (r *Registry) updateIndex(skill *Skill) {
	// 按分类索引
	r.index.byCategory[skill.Category] = appendUnique(r.index.byCategory[skill.Category], skill.ID)

	// 按标签索引
	for _, tag := range skill.Tags {
		r.index.byTag[tag] = appendUnique(r.index.byTag[tag], skill.ID)
	}

	// 按Agent类型索引
	for _, agentType := range skill.AgentTypes {
		r.index.byAgent[agentType] = appendUnique(r.index.byAgent[agentType], skill.ID)
	}

	// 按复杂度索引
	r.index.byComplexity[skill.Complexity] = appendUnique(r.index.byComplexity[skill.Complexity], skill.ID)
}

// rebuildIndex 重新构建索引
func (r *Registry) rebuildIndex() {
	r.index = newIndex()

	for _, skill := range r.skills {
		r.updateIndex(skill)
	}
}

// loadAllSkills 加载所有技能
func (r *Registry) loadAllSkills() error {
	metadatas, err := r.storage.List()
	if err != nil {
		return err
	}

	for _, metadata := range metadatas {
		skill, err := r.storage.Load(metadata.ID)
		if err != nil {
			// 记录错误但继续加载其他技能
			fmt.Printf("Warning: failed to load skill %s: %v\n", metadata.ID, err)
			continue
		}

		r.skills[skill.ID] = skill
		r.metadata[skill.ID] = skill.ToMetadata()
		r.updateIndex(skill)
	}

	return nil
}

// appendUnique 添加唯一元素
func appendUnique(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}

// contains 检查字符串是否包含子串（不区分大小写）
func contains(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}

// FileStorage 实现

// Save 保存技能到文件
func (fs *FileStorage) Save(skill *Skill) error {
	// 创建技能目录
	skillDir := filepath.Join(fs.baseDir, skill.ID)
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return fmt.Errorf("failed to create skill directory: %w", err)
	}

	// 保存为YAML
	yamlPath := filepath.Join(skillDir, "skill.yaml")
	yamlData, err := yaml.Marshal(skill)
	if err != nil {
		return fmt.Errorf("failed to marshal skill to YAML: %w", err)
	}

	if err := os.WriteFile(yamlPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write skill YAML: %w", err)
	}

	// 保存为JSON（用于快速读取元数据）
	jsonPath := filepath.Join(skillDir, "metadata.json")
	metadata := skill.ToMetadata()
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata to JSON: %w", err)
	}

	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write metadata JSON: %w", err)
	}

	// 保存Markdown文档
	mdPath := filepath.Join(skillDir, "SKILL.md")
	mdContent := generateMarkdown(skill)
	if err := os.WriteFile(mdPath, []byte(mdContent), 0644); err != nil {
		return fmt.Errorf("failed to write skill markdown: %w", err)
	}

	return nil
}

// Load 从文件加载技能
func (fs *FileStorage) Load(id string) (*Skill, error) {
	yamlPath := filepath.Join(fs.baseDir, id, "skill.yaml")

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read skill YAML: %w", err)
	}

	var skill Skill
	if err := yaml.Unmarshal(data, &skill); err != nil {
		return nil, fmt.Errorf("failed to unmarshal skill YAML: %w", err)
	}

	return &skill, nil
}

// Delete 删除技能文件
func (fs *FileStorage) Delete(id string) error {
	skillDir := filepath.Join(fs.baseDir, id)
	return os.RemoveAll(skillDir)
}

// List 列出所有技能元数据
func (fs *FileStorage) List() ([]*SkillMetadata, error) {
	var metadatas []*SkillMetadata

	err := filepath.Walk(fs.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理metadata.json文件
		if !info.IsDir() && info.Name() == "metadata.json" {
			data, err := os.ReadFile(path)
			if err != nil {
				// 记录错误但继续处理其他文件
				fmt.Printf("Warning: failed to read metadata file %s: %v\n", path, err)
				return nil
			}

			var metadata SkillMetadata
			if err := json.Unmarshal(data, &metadata); err != nil {
				fmt.Printf("Warning: failed to unmarshal metadata file %s: %v\n", path, err)
				return nil
			}

			metadatas = append(metadatas, &metadata)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk skills directory: %w", err)
	}

	return metadatas, nil
}

// Exists 检查技能是否存在
func (fs *FileStorage) Exists(id string) bool {
	yamlPath := filepath.Join(fs.baseDir, id, "skill.yaml")
	_, err := os.Stat(yamlPath)
	return err == nil
}

// generateMarkdown 生成Markdown文档
func generateMarkdown(skill *Skill) string {
	var md strings.Builder

	md.WriteString(fmt.Sprintf("# %s\n\n", skill.Name))
	md.WriteString(fmt.Sprintf("**版本**: %s\n\n", skill.Version))
	md.WriteString(fmt.Sprintf("**描述**: %s\n\n", skill.Description))
	md.WriteString(fmt.Sprintf("**作者**: %s\n\n", skill.Author))
	md.WriteString(fmt.Sprintf("**创建时间**: %s\n\n", skill.CreatedAt.Format("2006-01-02 15:04:05")))
	md.WriteString(fmt.Sprintf("**最后更新**: %s\n\n", skill.UpdatedAt.Format("2006-01-02 15:04:05")))

	md.WriteString("## 基本信息\n\n")
	md.WriteString(fmt.Sprintf("- **分类**: %s\n", skill.Category))
	md.WriteString(fmt.Sprintf("- **复杂度**: %s\n", skill.Complexity))
	md.WriteString(fmt.Sprintf("- **标签**: %s\n", strings.Join(skill.Tags, ", ")))
	md.WriteString(fmt.Sprintf("- **使用次数**: %d\n", skill.UsageCount))
	md.WriteString(fmt.Sprintf("- **成功率**: %.1f%%\n", skill.SuccessRate*100))
	md.WriteString(fmt.Sprintf("- **状态**: %s\n", getStatusText(skill.Enabled, skill.Validated)))

	if len(skill.Prerequisites) > 0 {
		md.WriteString("\n## 前置技能\n\n")
		for _, prereq := range skill.Prerequisites {
			md.WriteString(fmt.Sprintf("- %s\n", prereq))
		}
	}

	if len(skill.Steps) > 0 {
		md.WriteString("\n## 执行步骤\n\n")
		for _, step := range skill.Steps {
			md.WriteString(fmt.Sprintf("### 步骤 %d: %s\n\n", step.Order, step.Action))

			if step.Tool != "" {
				md.WriteString(fmt.Sprintf("**工具**: %s\n\n", step.Tool))
			}

			if len(step.Parameters) > 0 {
				md.WriteString("**参数**:\n\n")
				for _, param := range step.Parameters {
					required := ""
					if param.Required {
						required = " (必填)"
					}
					md.WriteString(fmt.Sprintf("- `%s` (%s)%s: %s\n",
						param.Name, param.Type, required, param.Description))
				}
				md.WriteString("\n")
			}

			if len(step.Conditions) > 0 {
				md.WriteString("**执行条件**:\n\n")
				for _, cond := range step.Conditions {
					md.WriteString(fmt.Sprintf("- %s: %s\n", cond.Type, cond.Check))
				}
				md.WriteString("\n")
			}

			if step.Expected != "" {
				md.WriteString(fmt.Sprintf("**预期结果**: %s\n\n", step.Expected))
			}
		}
	}

	if len(skill.Examples) > 0 {
		md.WriteString("\n## 使用示例\n\n")
		for i, example := range skill.Examples {
			md.WriteString(fmt.Sprintf("### 示例 %d\n\n", i+1))
			md.WriteString(fmt.Sprintf("%s\n\n", example))
		}
	}

	if len(skill.Tips) > 0 {
		md.WriteString("\n## 使用技巧\n\n")
		for _, tip := range skill.Tips {
			md.WriteString(fmt.Sprintf("- %s\n", tip))
		}
	}

	if len(skill.Warnings) > 0 {
		md.WriteString("\n## 注意事项\n\n")
		for _, warning := range skill.Warnings {
			md.WriteString(fmt.Sprintf("- ⚠️ %s\n", warning))
		}
	}

	if len(skill.RelatedTools) > 0 {
		md.WriteString("\n## 关联工具\n\n")
		for _, tool := range skill.RelatedTools {
			md.WriteString(fmt.Sprintf("- %s\n", tool))
		}
	}

	if len(skill.AgentTypes) > 0 {
		md.WriteString("\n## 适用Agent类型\n\n")
		for _, agentType := range skill.AgentTypes {
			md.WriteString(fmt.Sprintf("- %s\n", agentType))
		}
	}

	md.WriteString("\n## 性能指标\n\n")
	md.WriteString(fmt.Sprintf("- **平均执行时间**: %v\n", skill.AvgExecutionTime))
	md.WriteString(fmt.Sprintf("- **Token估算**: %d\n", skill.TokenEstimate))

	return md.String()
}

// getStatusText 获取状态文本
func getStatusText(enabled, validated bool) string {
	if !enabled {
		return "已禁用"
	}
	if validated {
		return "已启用（已验证）"
	}
	return "已启用（未验证）"
}
