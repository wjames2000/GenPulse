package genkit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	genkitconfig "GenPulse/internal/genkit/config"
	"GenPulse/internal/utils"
)

type StartupPhase int

const (
	PhaseUndefined StartupPhase = iota
	Phase1Critical
	Phase2Important
	Phase3Background
)

type PhaseConfig struct {
	Name    string
	Timeout time.Duration
	Phase   StartupPhase
}

type PhaseResult struct {
	Phase   StartupPhase
	Name    string
	Success bool
	Err     error
	Elapsed time.Duration
}

type StartupMetrics struct {
	mu           sync.Mutex
	StartTime    time.Time
	PhaseResults []PhaseResult
	TotalElapsed time.Duration
}

func (sm *StartupMetrics) RecordPhase(name string, phase StartupPhase, success bool, err error, elapsed time.Duration) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.PhaseResults = append(sm.PhaseResults, PhaseResult{
		Phase:   phase,
		Name:    name,
		Success: success,
		Err:     err,
		Elapsed: elapsed,
	})
}

func (sm *StartupMetrics) Report() map[string]interface{} {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	phases := make([]map[string]interface{}, len(sm.PhaseResults))
	for i, pr := range sm.PhaseResults {
		phases[i] = map[string]interface{}{
			"phase":   int(pr.Phase),
			"name":    pr.Name,
			"success": pr.Success,
			"elapsed": pr.Elapsed.Milliseconds(),
		}
	}
	return map[string]interface{}{
		"total_elapsed_ms": sm.TotalElapsed.Milliseconds(),
		"phases":           phases,
		"start_time":       sm.StartTime.Format(time.RFC3339),
	}
}

type PreloadConfig struct {
	Enabled     bool
	ConfigPaths []string
}

type StartupOptimizer struct {
	metrics        *StartupMetrics
	preloadConfig  PreloadConfig
	phase1Timeout  time.Duration
	lazyInitMCP    bool
	lazyInitSkills bool
	preloadData    map[string]interface{}
	preloadMu      sync.RWMutex
}

func NewStartupOptimizer() *StartupOptimizer {
	cfg := genkitconfig.GetConfig()
	optimizer := &StartupOptimizer{
		metrics: &StartupMetrics{
			StartTime: time.Now(),
		},
		phase1Timeout:  500 * time.Millisecond,
		lazyInitMCP:    true,
		lazyInitSkills: true,
		preloadConfig: PreloadConfig{
			Enabled: true,
			ConfigPaths: []string{
				"config/app_config.json",
			},
		},
		preloadData: make(map[string]interface{}),
	}

	if cfg != nil {
		if cfg.StartupPhase1TimeoutMs > 0 {
			optimizer.phase1Timeout = time.Duration(cfg.StartupPhase1TimeoutMs) * time.Millisecond
		}
		optimizer.lazyInitMCP = cfg.LazyInitMCP
		optimizer.lazyInitSkills = cfg.LazyInitSkills
		if cfg.PreloadConfig {
			optimizer.preloadConfig.Enabled = true
		}
	}

	return optimizer
}

func (so *StartupOptimizer) GetMetrics() *StartupMetrics {
	return so.metrics
}

func (so *StartupOptimizer) GetPhase1Timeout() time.Duration {
	return so.phase1Timeout
}

func (so *StartupOptimizer) LazyInitMCP() bool {
	return so.lazyInitMCP
}

func (so *StartupOptimizer) LazyInitSkills() bool {
	return so.lazyInitSkills
}

func (so *StartupOptimizer) PreloadConfigFiles() {
	if !so.preloadConfig.Enabled {
		return
	}

	var wg sync.WaitGroup
	for _, path := range so.preloadConfig.ConfigPaths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			absPath := p
			if !filepath.IsAbs(p) {
				absPath = filepath.Join(".", p)
			}
			data, err := os.ReadFile(absPath)
			if err != nil {
				utils.Warn("预加载配置文件失败: %s: %v", p, err)
				return
			}
			var parsed interface{}
			if err := json.Unmarshal(data, &parsed); err != nil {
				utils.Warn("解析预加载配置文件失败: %s: %v", p, err)
				return
			}
			so.preloadMu.Lock()
			so.preloadData[p] = parsed
			so.preloadMu.Unlock()
			utils.Info("预加载配置文件完成: %s (%d bytes)", p, len(data))
		}(path)
	}
	wg.Wait()
}

func (so *StartupOptimizer) GetPreloadedConfig(path string) interface{} {
	so.preloadMu.RLock()
	defer so.preloadMu.RUnlock()
	return so.preloadData[path]
}

func (so *StartupOptimizer) ReportMetrics() {
	so.metrics.TotalElapsed = time.Since(so.metrics.StartTime)
	report := so.metrics.Report()

	utils.Info("===== 启动速度指标 =====")
	utils.Info("总启动时间: %d ms", report["total_elapsed_ms"])
	for _, phase := range report["phases"].([]map[string]interface{}) {
		status := "OK"
		if !phase["success"].(bool) {
			status = "FAIL"
		}
		utils.Info("  阶段[%d] %s: %d ms [%s]", phase["phase"], phase["name"], phase["elapsed"], status)
	}
	utils.Info("=========================")
}

func RunPhaseWithTimeout(name string, phase StartupPhase, timeout time.Duration, metrics *StartupMetrics, fn func() error) error {
	start := time.Now()
	errCh := make(chan error, 1)

	go func() {
		errCh <- fn()
	}()

	var err error
	select {
	case err = <-errCh:
	case <-time.After(timeout):
		err = fmt.Errorf("阶段 '%s' 超时 (%v)", name, timeout)
		utils.Warn("启动阶段超时: %s (%v)", name, timeout)
	}

	elapsed := time.Since(start)
	success := err == nil
	metrics.RecordPhase(name, phase, success, err, elapsed)
	return err
}

var (
	globalStartupOptimizer *StartupOptimizer
	startupOnce            sync.Once
)

func GetStartupOptimizer() *StartupOptimizer {
	startupOnce.Do(func() {
		globalStartupOptimizer = NewStartupOptimizer()
	})
	return globalStartupOptimizer
}
