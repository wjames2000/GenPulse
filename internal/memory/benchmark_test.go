package memory

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func BenchmarkSearch(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "bench-search-*")
	if err != nil {
		b.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "bench.db")
	em, err := NewEpisodicMemory(dbPath)
	if err != nil {
		b.Fatalf("failed to create episodic memory: %v", err)
	}
	defer em.Close()

	insertBenchmarkRecords(b, em, 1000)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			results, err := em.Search(&SearchQuery{
				Query: "test content",
				Limit: 20,
			})
			if err != nil {
				b.Fatalf("search failed: %v", err)
			}
			if len(results) == 0 {
				b.Fatal("expected at least one result")
			}
		}
	})
}

func BenchmarkSearchWithCache(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "bench-cache-*")
	if err != nil {
		b.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "bench.db")
	em, err := NewEpisodicMemory(dbPath)
	if err != nil {
		b.Fatalf("failed to create episodic memory: %v", err)
	}
	defer em.Close()

	insertBenchmarkRecords(b, em, 1000)

	se := NewSearchEngine(nil, em, nil)
	se.SetSearchTimeout(30 * time.Second)
	defer func() {
		se.InvalidateCache()
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := se.Search(context.Background(), &MemoryQuery{
				Query:     "test content",
				Limit:     20,
				IncludeL2: true,
			})
			if err != nil {
				b.Fatalf("search failed: %v", err)
			}
			if len(resp.Results) == 0 {
				b.Fatal("expected at least one result")
			}
		}
	})
}

func BenchmarkBatchInsert(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "bench-batch-*")
	if err != nil {
		b.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "bench.db")
	em, err := NewEpisodicMemory(dbPath)
	if err != nil {
		b.Fatalf("failed to create episodic memory: %v", err)
	}
	defer em.Close()

	records := make([]*MemoryRecord, 100)
	for i := 0; i < 100; i++ {
		records[i] = &MemoryRecord{
			ID:          fmt.Sprintf("bench-batch-%d", i),
			SessionID:   "bench-session",
			TaskID:      fmt.Sprintf("task-%d", i),
			TaskType:    "benchmark",
			Description: fmt.Sprintf("benchmark record %d", i),
			Content:     fmt.Sprintf("this is benchmark content for record number %d with some extra text for fts search", i),
			Metadata:    map[string]any{"index": i},
			Tags:        []string{"bench", "test"},
			Category:    "benchmark",
			Importance:  0.5,
			Success:     true,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := em.StoreBatch(records); err != nil {
			b.Fatalf("batch insert failed: %v", err)
		}
	}
}

func BenchmarkConcurrentSearch(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "bench-concurrent-*")
	if err != nil {
		b.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "bench.db")
	em, err := NewEpisodicMemory(dbPath)
	if err != nil {
		b.Fatalf("failed to create episodic memory: %v", err)
	}
	defer em.Close()

	insertBenchmarkRecords(b, em, 500)

	queries := []string{
		"test content",
		"code generation",
		"error fix",
		"testing",
		"benchmark",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var idx int
		for pb.Next() {
			q := queries[idx%len(queries)]
			idx++
			results, err := em.Search(&SearchQuery{
				Query: q,
				Limit: 10,
			})
			if err != nil {
				b.Fatalf("concurrent search failed: %v", err)
			}
			_ = results
		}
	})
}

func BenchmarkSearchWithPreparedStmt(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "bench-prep-*")
	if err != nil {
		b.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "bench.db")
	em, err := NewEpisodicMemory(dbPath)
	if err != nil {
		b.Fatalf("failed to create episodic memory: %v", err)
	}
	defer em.Close()

	insertBenchmarkRecords(b, em, 1000)

	if err := em.prepareStatements(); err != nil {
		b.Fatalf("failed to prepare statements: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			results, err := em.Search(&SearchQuery{
				Query: "test content",
				Limit: 20,
			})
			if err != nil {
				b.Fatalf("search failed: %v", err)
			}
			if len(results) == 0 {
				b.Fatal("expected at least one result")
			}
		}
	})
}

func insertBenchmarkRecords(b *testing.B, em *EpisodicMemory, n int) {
	b.Helper()
	records := make([]*MemoryRecord, n)
	for i := 0; i < n; i++ {
		records[i] = &MemoryRecord{
			ID:          fmt.Sprintf("bench-%d", i),
			SessionID:   fmt.Sprintf("session-%d", i%10),
			TaskID:      fmt.Sprintf("task-%d", i),
			TaskType:    []string{"code_generation", "code_review", "error_fix", "testing", "benchmark"}[i%5],
			Description: fmt.Sprintf("benchmark test record %d", i),
			Content:     fmt.Sprintf("this is a test content for benchmarking fts search performance with record number %d", i),
			Metadata:    map[string]any{"index": i, "type": "bench"},
			Tags:        []string{"bench", "test", fmt.Sprintf("tag-%d", i%20)},
			Category:    []string{"development", "quality", "maintenance", "testing"}[i%4],
			Importance:  float64(i%100) / 100.0,
			Success:     i%2 == 0,
			CreatedAt:   time.Now().Add(-time.Duration(i) * time.Hour),
		}
	}

	if err := em.StoreBatch(records); err != nil {
		b.Fatalf("failed to insert benchmark records: %v", err)
	}
}

var (
	benchSearchResult    []*SearchResult
	benchSearchResponse  *SearchResponse
	benchCacheOnce       sync.Once
	benchCacheInstance   *SearchCache
	benchEpisodicMemOnce sync.Once
	benchEpisodicMem     *EpisodicMemory
	benchEmCleanup       func()
)

func getBenchEpisodicMemory(b *testing.B) *EpisodicMemory {
	benchEpisodicMemOnce.Do(func() {
		tempDir, err := os.MkdirTemp("", "bench-global-*")
		if err != nil {
			b.Fatalf("failed to create temp dir: %v", err)
		}
		dbPath := filepath.Join(tempDir, "bench.db")
		em, err := NewEpisodicMemory(dbPath)
		if err != nil {
			b.Fatalf("failed to create episodic memory: %v", err)
		}
		insertBenchmarkRecords(b, em, 2000)
		benchEpisodicMem = em
		benchEmCleanup = func() {
			em.Close()
			os.RemoveAll(tempDir)
		}
	})
	return benchEpisodicMem
}

func getBenchCache() *SearchCache {
	benchCacheOnce.Do(func() {
		benchCacheInstance = NewSearchCache(1000, 30*time.Minute)
	})
	return benchCacheInstance
}

func TestMain(m *testing.M) {
	code := m.Run()
	if benchEmCleanup != nil {
		benchEmCleanup()
	}
	if benchCacheInstance != nil {
		benchCacheInstance.Stop()
	}
	os.Exit(code)
}
