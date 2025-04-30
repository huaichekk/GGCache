package cache

import (
	"GGCache/configs"
	"GGCache/internal/cache/eviction"
	"crypto/md5"
	"encoding/hex"
	"sync"
)

type SafeCache struct {
	shardMutex []*sync.RWMutex
	cache      []eviction.Cache
}

func NewSafeCache(cap int) *SafeCache {
	shardMutex := make([]*sync.RWMutex, configs.GetConfig().CacheConfig.ShardNum)
	cacheSlice := make([]eviction.Cache, configs.GetConfig().CacheConfig.ShardNum)
	for i := 0; i < configs.GetConfig().CacheConfig.ShardNum; i++ {
		shardMutex[i] = &sync.RWMutex{}
		cacheSlice[i] = eviction.GetCacheByEviction(configs.GetConfig().CacheConfig.Eviction,
			cap)
	}
	return &SafeCache{
		shardMutex: shardMutex,
		cache:      cacheSlice,
	}
}

// getShardIndex 获取键对应的分片索引
func getShardIndex(key string) int {
	// 计算MD5哈希
	hasher := md5.New()
	hasher.Write([]byte(key))
	hashBytes := hasher.Sum(nil)

	// 将哈希转换为16进制字符串
	hashStr := hex.EncodeToString(hashBytes)

	// 取前8个字符作为uint32
	var hash uint32
	for i := 0; i < 8; i++ {
		hash = hash<<4 | uint32(hashStr[i]%16)
	}

	// 取模得到分片索引
	return int(hash % uint32(configs.GetConfig().CacheConfig.ShardNum))
}

func (s *SafeCache) Get(key string) (eviction.ByteView, bool) {
	index := getShardIndex(key)
	mu := s.shardMutex[index]
	mu.RLock()
	defer mu.RUnlock()
	//s.shardMutex.RLock()
	//defer s.shardMutex.RUnlock()
	if v, ok := s.cache[index].Get(key); ok {
		return v.(eviction.ByteView), true
	} else {
		return eviction.ByteView{}, false
	}
}

func (s *SafeCache) Put(key string, value eviction.ByteView) {
	index := getShardIndex(key)
	mu := s.shardMutex[index]
	mu.Lock()
	defer mu.Unlock()
	//s.shardMutex.Lock()
	//defer s.shardMutex.Unlock()
	s.cache[index].Put(key, value)
}
