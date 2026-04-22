package skills

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// SkillSource 技能来源
type SkillSource string

const (
	SourceClaudeOfficial    SkillSource = "claude_official"
	SourceAnthropicOfficial SkillSource = "anthropic_official"
	SourceOfficialRepo      SkillSource = "official_repo"
	SourceGitHub            SkillSource = "github"
	SourceCustomURL         SkillSource = "custom_url"
	SourceMarketplace       SkillSource = "marketplace"
)

// 官方技能仓库配置
const (
	OfficialRepoOwner  = "hpds-cc"
	OfficialRepoName   = "genpulse-skills"
	OfficialRepoBranch = "main"
	OfficialSkillDir   = "skills"
)

// OnlineSkill 在线技能表示（来自外部来源的标准化格式）
type OnlineSkill struct {
	ID            string         `json:"id" yaml:"id"`
	Name          string         `json:"name" yaml:"name"`
	Version       string         `json:"version" yaml:"version"`
	Description   string         `json:"description" yaml:"description"`
	Author        string         `json:"author" yaml:"author"`
	Category      string         `json:"category" yaml:"category"`
	Tags          []string       `json:"tags" yaml:"tags"`
	Complexity    string         `json:"complexity" yaml:"complexity"`
	Source        SkillSource    `json:"source" yaml:"source"`
	SourceURL     string         `json:"source_url" yaml:"source_url"`
	Steps         []OnlineStep   `json:"steps" yaml:"steps"`
	Examples      []string       `json:"examples" yaml:"examples"`
	Tips          []string       `json:"tips" yaml:"tips"`
	Warnings      []string       `json:"warnings" yaml:"warnings"`
	Prerequisites []string       `json:"prerequisites" yaml:"prerequisites"`
	RelatedTools  []string       `json:"related_tools" yaml:"related_tools"`
	AgentTypes    []string       `json:"agent_types" yaml:"agent_types"`
	TokenEstimate int            `json:"token_estimate" yaml:"token_estimate"`
	Metadata      map[string]any `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

type OnlineStep struct {
	Order      int         `json:"order" yaml:"order"`
	Action     string      `json:"action" yaml:"action"`
	Tool       string      `json:"tool,omitempty" yaml:"tool,omitempty"`
	Parameters []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Expected   string      `json:"expected,omitempty" yaml:"expected,omitempty"`
}

// OnlineImporter 在线技能导入器
type OnlineImporter struct {
	client   *http.Client
	registry *Registry
	sources  map[SkillSource]SourceFetcher
}

// SourceFetcher 来源获取器接口
type SourceFetcher interface {
	Name() string
	Fetch(ctx context.Context, query string) ([]*OnlineSkill, error)
	FetchByID(ctx context.Context, id string) (*OnlineSkill, error)
}

// NewOnlineImporter 创建在线导入器
func NewOnlineImporter(registry *Registry) *OnlineImporter {
	importer := &OnlineImporter{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		registry: registry,
		sources:  make(map[SkillSource]SourceFetcher),
	}

	importer.sources[SourceCustomURL] = &URLFetcher{client: importer.client}
	importer.sources[SourceMarketplace] = &MarketplaceFetcher{client: importer.client}
	importer.sources[SourceOfficialRepo] = NewOfficialSkillFetcher(importer.client)

	return importer
}

// RegisterSource 注册自定义来源
func (oi *OnlineImporter) RegisterSource(source SkillSource, fetcher SourceFetcher) {
	oi.sources[source] = fetcher
}

// ListSources 列出所有可用来源
func (oi *OnlineImporter) ListSources() []map[string]any {
	var result []map[string]any
	for source, fetcher := range oi.sources {
		result = append(result, map[string]any{
			"id":   string(source),
			"name": fetcher.Name(),
		})
	}
	return result
}

// SearchOnline 搜索在线技能
func (oi *OnlineImporter) SearchOnline(ctx context.Context, query string, source SkillSource) ([]*OnlineSkill, error) {
	fetcher, ok := oi.sources[source]
	if !ok {
		return nil, fmt.Errorf("unsupported source: %s", source)
	}
	return fetcher.Fetch(ctx, query)
}

// ImportFromOnline 从在线来源导入技能
func (oi *OnlineImporter) ImportFromOnline(ctx context.Context, source SkillSource, id string) (*Skill, error) {
	fetcher, ok := oi.sources[source]
	if !ok {
		return nil, fmt.Errorf("unsupported source: %s", source)
	}

	onlineSkill, err := fetcher.FetchByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch skill from %s: %w", source, err)
	}

	return oi.convertAndRegister(onlineSkill)
}

// ImportFromURL 从URL导入技能
func (oi *OnlineImporter) ImportFromURL(ctx context.Context, rawURL string) (*Skill, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	var onlineSkill *OnlineSkill

	switch {
	case strings.Contains(u.Host, "github.com"):
		onlineSkill, err = fetchFromGitHub(ctx, oi.client, rawURL)
	case strings.HasSuffix(rawURL, ".yaml") || strings.HasSuffix(rawURL, ".yml"):
		onlineSkill, err = fetchFromYAMLURL(ctx, oi.client, rawURL)
	case strings.HasSuffix(rawURL, ".json"):
		onlineSkill, err = fetchFromJSONURL(ctx, oi.client, rawURL)
	default:
		onlineSkill, err = fetchFromGenericURL(ctx, oi.client, rawURL)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch from URL: %w", err)
	}

	return oi.convertAndRegister(onlineSkill)
}

// ImportFromRawData 从原始数据导入技能
func (oi *OnlineImporter) ImportFromRawData(ctx context.Context, data []byte, format string, source SkillSource) (*Skill, error) {
	var onlineSkill OnlineSkill
	var err error

	switch format {
	case "yaml", "yml":
		err = yaml.Unmarshal(data, &onlineSkill)
	case "json":
		err = json.Unmarshal(data, &onlineSkill)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse skill data: %w", err)
	}

	onlineSkill.Source = source

	return oi.convertAndRegister(&onlineSkill)
}

// convertAndRegister 转换在线技能为内部技能并注册
func (oi *OnlineImporter) convertAndRegister(online *OnlineSkill) (*Skill, error) {
	skill := NewSkill(online.Name, online.Description, online.Author)

	if online.Version != "" {
		skill.Version = online.Version
	}
	if online.Category != "" {
		skill.Category = online.Category
	}
	if len(online.Tags) > 0 {
		skill.Tags = online.Tags
	}
	if online.Complexity != "" {
		skill.Complexity = online.Complexity
	}
	if len(online.Examples) > 0 {
		skill.Examples = online.Examples
	}
	if len(online.Tips) > 0 {
		skill.Tips = online.Tips
	}
	if len(online.Warnings) > 0 {
		skill.Warnings = online.Warnings
	}
	if len(online.Prerequisites) > 0 {
		skill.Prerequisites = online.Prerequisites
	}
	if len(online.RelatedTools) > 0 {
		skill.RelatedTools = online.RelatedTools
	}
	if len(online.AgentTypes) > 0 {
		skill.AgentTypes = online.AgentTypes
	}
	if online.TokenEstimate > 0 {
		skill.TokenEstimate = online.TokenEstimate
	}

	for _, os := range online.Steps {
		step := Step{
			ID:         generateStepID(skill.ID, os.Order),
			Order:      os.Order,
			Action:     os.Action,
			Tool:       os.Tool,
			Parameters: os.Parameters,
			Expected:   os.Expected,
		}
		skill.Steps = append(skill.Steps, step)
	}

	if err := oi.registry.Register(skill); err != nil {
		return nil, fmt.Errorf("failed to register imported skill: %w", err)
	}

	return skill, nil
}

// ==================== URL Fetcher ====================

type URLFetcher struct {
	client *http.Client
}

func (f *URLFetcher) Name() string {
	return "自定义 URL"
}

func (f *URLFetcher) Fetch(ctx context.Context, query string) ([]*OnlineSkill, error) {
	return nil, fmt.Errorf("URL fetcher does not support search")
}

func (f *URLFetcher) FetchByID(ctx context.Context, url string) (*OnlineSkill, error) {
	return fetchFromGenericURL(ctx, f.client, url)
}

// ==================== Marketplace Fetcher ====================

type MarketplaceFetcher struct {
	client *http.Client
}

func (f *MarketplaceFetcher) Name() string {
	return "技能市场"
}

func (f *MarketplaceFetcher) Fetch(ctx context.Context, query string) ([]*OnlineSkill, error) {
	return nil, fmt.Errorf("marketplace search not yet implemented")
}

func (f *MarketplaceFetcher) FetchByID(ctx context.Context, id string) (*OnlineSkill, error) {
	return nil, fmt.Errorf("marketplace fetch by ID not yet implemented")
}

// ==================== GitHub Fetcher ====================

func fetchFromGitHub(ctx context.Context, client *http.Client, rawURL string) (*OnlineSkill, error) {
	apiURL := convertGitHubToRawURL(rawURL)
	if apiURL == "" {
		return nil, fmt.Errorf("unable to parse GitHub URL: %s", rawURL)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3.raw")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return parseSkillFromData(data, rawURL)
}

func convertGitHubToRawURL(githubURL string) string {
	u, err := url.Parse(githubURL)
	if err != nil {
		return ""
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return ""
	}

	owner, repo := parts[0], parts[1]

	if strings.Contains(u.Path, "/blob/") {
		blobIdx := indexOf(parts, "blob")
		if blobIdx >= 0 && blobIdx+2 < len(parts) {
			branch := parts[blobIdx+1]
			filePath := strings.Join(parts[blobIdx+2:], "/")
			return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, branch, filePath)
		}
	}

	return ""
}

func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}

// ==================== URL-based fetchers ====================

func fetchFromYAMLURL(ctx context.Context, client *http.Client, rawURL string) (*OnlineSkill, error) {
	data, err := fetchURL(ctx, client, rawURL)
	if err != nil {
		return nil, err
	}
	return parseSkillFromData(data, rawURL)
}

func fetchFromJSONURL(ctx context.Context, client *http.Client, rawURL string) (*OnlineSkill, error) {
	data, err := fetchURL(ctx, client, rawURL)
	if err != nil {
		return nil, err
	}
	return parseSkillFromData(data, rawURL)
}

func fetchFromGenericURL(ctx context.Context, client *http.Client, rawURL string) (*OnlineSkill, error) {
	data, err := fetchURL(ctx, client, rawURL)
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(rawURL)
	path := u.Path

	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		return parseSkillFromData(data, rawURL)
	}
	return parseSkillFromData(data, rawURL)
}

func fetchURL(ctx context.Context, client *http.Client, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "GenPulse/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("URL returned status %d", resp.StatusCode)
	}

	limitedReader := io.LimitReader(resp.Body, 10*1024*1024)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return data, nil
}

// parseSkillFromData 从数据解析技能（尝试YAML和JSON）
func parseSkillFromData(data []byte, sourceURL string) (*OnlineSkill, error) {
	var skill OnlineSkill

	if err := yaml.Unmarshal(data, &skill); err == nil && skill.Name != "" {
		skill.Source = detectSourceFromURL(sourceURL)
		skill.SourceURL = sourceURL
		return &skill, nil
	}

	if err := json.Unmarshal(data, &skill); err == nil && skill.Name != "" {
		skill.Source = detectSourceFromURL(sourceURL)
		skill.SourceURL = sourceURL
		return &skill, nil
	}

	return nil, fmt.Errorf("unable to parse skill data: data does not contain valid skill definition")
}

func detectSourceFromURL(rawURL string) SkillSource {
	u, err := url.Parse(rawURL)
	if err != nil {
		return SourceCustomURL
	}

	if strings.Contains(u.Host, "github.com") {
		return SourceGitHub
	}
	if strings.Contains(u.Host, "claude.ai") || strings.Contains(u.Host, "anthropic.com") {
		return SourceAnthropicOfficial
	}

	return SourceCustomURL
}

// ==================== SkillManager integration ====================

// GetOnlineImporter 获取在线导入器
func (sm *SkillManager) GetOnlineImporter() *OnlineImporter {
	return sm.onlineImporter
}

// Add onlineImporter field to SkillManager
// This is declared separately - the actual field is added via the init function approach

// SetOnlineImporter 设置在线导入器（在NewSkillManager中调用）
func (sm *SkillManager) SetOnlineImporter(oi *OnlineImporter) {
	sm.onlineImporter = oi
}

// OnlineImportResult 在线导入结果
type OnlineImportResult struct {
	Skill  *Skill `json:"skill"`
	Source string `json:"source"`
}

// SearchOnlineSkills 搜索在线技能
func (sm *SkillManager) SearchOnlineSkills(ctx context.Context, query string, source SkillSource) ([]*OnlineSkill, error) {
	if sm.onlineImporter == nil {
		return nil, fmt.Errorf("online importer not initialized")
	}
	return sm.onlineImporter.SearchOnline(ctx, query, source)
}

// ImportOnlineSkill 从在线来源导入技能
func (sm *SkillManager) ImportOnlineSkill(ctx context.Context, source SkillSource, id string) (*Skill, error) {
	if sm.onlineImporter == nil {
		return nil, fmt.Errorf("online importer not initialized")
	}
	return sm.onlineImporter.ImportFromOnline(ctx, source, id)
}

// ImportSkillFromURL 从URL导入技能
func (sm *SkillManager) ImportSkillFromURL(ctx context.Context, rawURL string) (*Skill, error) {
	if sm.onlineImporter == nil {
		return nil, fmt.Errorf("online importer not initialized")
	}
	return sm.onlineImporter.ImportFromURL(ctx, rawURL)
}

// ListOnlineSources 列出在线技能来源
func (sm *SkillManager) ListOnlineSources() []map[string]any {
	if sm.onlineImporter == nil {
		return nil
	}
	return sm.onlineImporter.ListSources()
}

// ==================== Official Skill Fetcher ====================

// OfficialSkillIndex 官方技能仓库索引
type OfficialSkillIndex struct {
	Skills []OfficialSkillEntry `json:"skills"`
}

type OfficialSkillEntry struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Complexity  string   `json:"complexity"`
	Author      string   `json:"author"`
	Version     string   `json:"version"`
	FileName    string   `json:"file_name"`
}

type OfficialSkillFetcher struct {
	client   *http.Client
	owner    string
	repo     string
	branch   string
	skillDir string
}

func NewOfficialSkillFetcher(client *http.Client) *OfficialSkillFetcher {
	return &OfficialSkillFetcher{
		client:   client,
		owner:    OfficialRepoOwner,
		repo:     OfficialRepoName,
		branch:   OfficialRepoBranch,
		skillDir: OfficialSkillDir,
	}
}

func (f *OfficialSkillFetcher) Name() string {
	return "官方技能库"
}

func (f *OfficialSkillFetcher) Fetch(ctx context.Context, query string) ([]*OnlineSkill, error) {
	index, err := f.fetchIndex(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch official skill index: %w", err)
	}

	var result []*OnlineSkill
	for _, entry := range index.Skills {
		if query != "" {
			if !contains(entry.Name, query) && !contains(entry.Description, query) {
				matched := false
				for _, tag := range entry.Tags {
					if contains(tag, query) {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			}
		}
		result = append(result, &OnlineSkill{
			ID:          entry.ID,
			Name:        entry.Name,
			Description: entry.Description,
			Category:    entry.Category,
			Tags:        entry.Tags,
			Complexity:  entry.Complexity,
			Author:      entry.Author,
			Version:     entry.Version,
			Source:      SourceOfficialRepo,
		})
	}

	return result, nil
}

func (f *OfficialSkillFetcher) FetchByID(ctx context.Context, id string) (*OnlineSkill, error) {
	index, err := f.fetchIndex(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch official skill index: %w", err)
	}

	var entry *OfficialSkillEntry
	for _, e := range index.Skills {
		if e.ID == id || e.Name == id {
			entry = &e
			break
		}
	}
	if entry == nil {
		return nil, fmt.Errorf("skill '%s' not found in official repository", id)
	}

	rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s/%s",
		f.owner, f.repo, f.branch, f.skillDir, entry.FileName)
	return fetchFromYAMLURL(ctx, f.client, rawURL)
}

func (f *OfficialSkillFetcher) fetchIndex(ctx context.Context) (*OfficialSkillIndex, error) {
	indexURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s/index.json",
		f.owner, f.repo, f.branch, f.skillDir)

	data, err := fetchURL(ctx, f.client, indexURL)
	if err != nil {
		return nil, err
	}

	var index OfficialSkillIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("failed to parse skill index: %w", err)
	}

	return &index, nil
}
