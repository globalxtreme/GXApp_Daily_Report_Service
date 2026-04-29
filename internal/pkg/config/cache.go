package config

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	XtremeCache *cache.Cache
)

func InitCache(defaultExpiration, cleanupInterval time.Duration) {
	XtremeCache = cache.New(defaultExpiration, cleanupInterval)
}
