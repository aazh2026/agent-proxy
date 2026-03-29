package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"
	"sync"
	"time"
)

type CacheEntry struct {
	Response  []byte
	CreatedAt time.Time
	HitCount  int
}

type Cache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	order   []string
	maxSize int
	ttl     time.Duration
	stats   *Stats
}

type Stats struct {
	mu        sync.RWMutex
	Hits      int64
	Misses    int64
	Evictions int64
	Size      int
}

type Config struct {
	MaxSize         int  `json:"max_size"`
	TTLSeconds      int  `json:"ttl_seconds"`
	Enabled         bool `json:"enabled"`
	MaxResponseSize int  `json:"max_response_size"`
}

func DefaultConfig() *Config {
	return &Config{
		MaxSize:         1000,
		TTLSeconds:      3600,
		Enabled:         true,
		MaxResponseSize: 1024 * 1024,
	}
}

func New(config *Config) *Cache {
	if config == nil {
		config = DefaultConfig()
	}

	return &Cache{
		entries: make(map[string]*CacheEntry),
		order:   make([]string, 0, config.MaxSize),
		maxSize: config.MaxSize,
		ttl:     time.Duration(config.TTLSeconds) * time.Second,
		stats:   &Stats{},
	}
}

func GenerateKey(model string, messages []Message) string {
	var parts []string
	parts = append(parts, model)

	for _, msg := range messages {
		normalized := normalizeContent(msg.Content)
		parts = append(parts, msg.Role+":"+normalized)
	}

	combined := strings.Join(parts, "|")
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

func GenerateKeyFromInput(model string, input interface{}) string {
	inputBytes, _ := json.Marshal(input)
	combined := model + "|" + string(inputBytes)
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

func normalizeContent(content string) string {
	content = strings.ToLower(content)
	content = strings.TrimSpace(content)
	content = strings.Join(strings.Fields(content), " ")
	return content
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	entry, exists := c.entries[key]
	c.mu.RUnlock()

	if !exists {
		c.stats.recordMiss()
		return nil, false
	}

	if time.Since(entry.CreatedAt) > c.ttl {
		c.mu.Lock()
		delete(c.entries, key)
		c.removeFromOrder(key)
		c.mu.Unlock()
		c.stats.recordMiss()
		return nil, false
	}

	c.mu.Lock()
	entry.HitCount++
	c.mu.Unlock()

	c.stats.recordHit()
	return entry.Response, true
}

func (c *Cache) Set(key string, response []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.entries) >= c.maxSize {
		c.evictOldest()
	}

	c.entries[key] = &CacheEntry{
		Response:  response,
		CreatedAt: time.Now(),
		HitCount:  0,
	}
	c.order = append(c.order, key)
	c.stats.setSize(len(c.entries))
}

func (c *Cache) evictOldest() {
	if len(c.order) == 0 {
		return
	}

	oldestKey := c.order[0]
	delete(c.entries, oldestKey)
	c.order = c.order[1:]
	c.stats.recordEviction()
}

func (c *Cache) removeFromOrder(key string) {
	for i, k := range c.order {
		if k == key {
			c.order = append(c.order[:i], c.order[i+1:]...)
			break
		}
	}
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*CacheEntry)
	c.order = make([]string, 0, c.maxSize)
	c.stats.setSize(0)
}

func (c *Cache) InvalidateByPrefix(prefix string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0
	for key := range c.entries {
		if strings.HasPrefix(key, prefix) {
			delete(c.entries, key)
			c.removeFromOrder(key)
			count++
		}
	}
	c.stats.setSize(len(c.entries))
	return count
}

func (c *Cache) GetStats() map[string]interface{} {
	c.stats.mu.RLock()
	defer c.stats.mu.RUnlock()

	c.mu.RLock()
	size := len(c.entries)
	c.mu.RUnlock()

	total := c.stats.Hits + c.stats.Misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(c.stats.Hits) / float64(total) * 100
	}

	return map[string]interface{}{
		"hits":        c.stats.Hits,
		"misses":      c.stats.Misses,
		"evictions":   c.stats.Evictions,
		"size":        size,
		"max_size":    c.maxSize,
		"hit_rate":    hitRate,
		"ttl_seconds": int(c.ttl.Seconds()),
	}
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

func (s *Stats) recordHit() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Hits++
}

func (s *Stats) recordMiss() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Misses++
}

func (s *Stats) recordEviction() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Evictions++
}

func (s *Stats) setSize(size int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Size = size
}

func SortMessages(messages []Message) []Message {
	sorted := make([]Message, len(messages))
	copy(sorted, messages)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Role != sorted[j].Role {
			return sorted[i].Role < sorted[j].Role
		}
		return sorted[i].Content < sorted[j].Content
	})
	return sorted
}
