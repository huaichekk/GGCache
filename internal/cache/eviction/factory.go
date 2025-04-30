package eviction

func GetCacheByEviction(evict string, cap int) Cache {
	if evict == "lru" {
		return NewLRUCache(cap)
	}
	return nil
}
