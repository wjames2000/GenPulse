package memory

import (
	"container/list"
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type cacheEntry struct {
	key       string
	value     []*SearchResult
	expiresAt time.Time
	size      int
}

type SearchCache struct {
	mu          sync.RWMutex
	maxSize     int
	ttl         time.Duration
	items       map[string]*list.Element
	lruList     *list.List
	hits        int64
	misses      int64
	stopCh      chan struct{}
	cleanupTick time.Duration
}

type cacheKey struct {
	query       string
	filtersHash string
	limit       int
	offset      int
	sortBy      string
	sortOrder   string
}

func NewSearchCache(maxSize int, ttl time.Duration) *SearchCache {
	if maxSize <= 0 {
		maxSize = 1000
	}
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}

	c := &SearchCache{
		maxSize:     maxSize,
		ttl:         ttl,
		items:       make(map[string]*list.Element),
		lruList:     list.New(),
		stopCh:      make(chan struct{}),
		cleanupTick: time.Minute,
	}

	go c.cleanupLoop()
	return c
}

func (c *SearchCache) Stop() {
	close(c.stopCh)
}

func (k *cacheKey) String() string {
	return fmt.Sprintf("%s|%s|%d|%d|%s|%s",
		strings.ToLower(strings.TrimSpace(k.query)),
		k.filtersHash,
		k.limit, k.offset, k.sortBy, k.sortOrder)
}

func NormalizeCacheKey(query string, filters map[string]any, limit, offset int, sortBy, sortOrder string) string {
	fh := filtersHash(filters)
	ck := &cacheKey{
		query:       query,
		filtersHash: fh,
		limit:       limit,
		offset:      offset,
		sortBy:      sortBy,
		sortOrder:   sortOrder,
	}
	return ck.String()
}

func filtersHash(filters map[string]any) string {
	if len(filters) == 0 {
		return ""
	}
	keys := make([]string, 0, len(filters))
	for k := range filters {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(fmt.Sprintf("%v", filters[k]))
		b.WriteString("&")
	}
	h := sha256.Sum256([]byte(b.String()))
	return fmt.Sprintf("%x", h[:8])
}

func (c *SearchCache) Get(key string) ([]*SearchResult, bool) {
	c.mu.RLock()
	elem, exists := c.items[key]
	if !exists {
		c.mu.RUnlock()
		c.mu.Lock()
		c.misses++
		c.mu.Unlock()
		return nil, false
	}

	entry := elem.Value.(*cacheEntry)
	if time.Now().After(entry.expiresAt) {
		c.mu.RUnlock()
		c.mu.Lock()
		c.removeElement(elem)
		c.misses++
		c.mu.Unlock()
		return nil, false
	}

	c.mu.RUnlock()

	c.mu.Lock()
	c.lruList.MoveToFront(elem)
	c.hits++
	c.mu.Unlock()

	return entry.value, true
}

func (c *SearchCache) Set(key string, value []*SearchResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.items[key]; exists {
		c.lruList.MoveToFront(elem)
		elem.Value.(*cacheEntry).value = value
		elem.Value.(*cacheEntry).expiresAt = time.Now().Add(c.ttl)
		return
	}

	entry := &cacheEntry{
		key:       key,
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
		size:      len(value),
	}
	elem := c.lruList.PushFront(entry)
	c.items[key] = elem

	if c.lruList.Len() > c.maxSize {
		c.removeOldest()
	}
}

func (c *SearchCache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element)
	c.lruList.Init()
}

func (c *SearchCache) InvalidatePattern(pattern string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, elem := range c.items {
		if strings.Contains(key, pattern) {
			c.removeElement(elem)
			delete(c.items, key)
		}
	}
}

func (c *SearchCache) Stats() (hits, misses int64, size int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hits, c.misses, c.lruList.Len()
}

func (c *SearchCache) HitRate() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	total := c.hits + c.misses
	if total == 0 {
		return 0
	}
	return float64(c.hits) / float64(total)
}

func (c *SearchCache) removeElement(elem *list.Element) {
	c.lruList.Remove(elem)
	delete(c.items, elem.Value.(*cacheEntry).key)
}

func (c *SearchCache) removeOldest() {
	elem := c.lruList.Back()
	if elem != nil {
		c.removeElement(elem)
	}
}

func (c *SearchCache) cleanupLoop() {
	ticker := time.NewTicker(c.cleanupTick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.removeExpired()
		case <-c.stopCh:
			return
		}
	}
}

func (c *SearchCache) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, elem := range c.items {
		if now.After(elem.Value.(*cacheEntry).expiresAt) {
			c.removeElement(elem)
			delete(c.items, key)
		}
	}
}
