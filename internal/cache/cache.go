// internal/cache/cache.go
package cache

import (
	"SomeProject/internal/models"
	"sync"
)

type BannerCache struct {
	mu      sync.Mutex
	banners map[string]models.Banner
}

func NewBannerCache() *BannerCache {
	return &BannerCache{
		banners: make(map[string]models.Banner),
	}
}

func (c *BannerCache) Get(key string) (models.Banner, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	banner, found := c.banners[key]
	return banner, found
}

func (c *BannerCache) Set(key string, banner models.Banner) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.banners[key] = banner
}
