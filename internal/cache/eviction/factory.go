package eviction

func GetCacheByEviction(evict string, cap, minTTL, maxTTL int) Cache {
	if evict == "lru" {
		return NewLRUCache(cap, minTTL, maxTTL)
	}
	return nil
}
