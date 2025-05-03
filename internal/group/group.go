package group

import (
	"GGCache/internal/cache"
	"GGCache/internal/cache/eviction"
	"GGCache/internal/singleflight"
	"fmt"
	"sync"
)

var (
	mu     = sync.RWMutex{}
	groups = make(map[string]*Group)
)

type GetterFunc func(key string) ([]byte, bool)

type Group struct {
	name       string
	mainCache  *cache.SafeCache
	getterFunc GetterFunc
	single     singleflight.Singleflight
}

func NewGroup(name string, cap, minTTL, maxTTL int, getterFunc GetterFunc) *Group {
	if getterFunc == nil {
		panic("Group 回调函数为空")
	}
	g := &Group{
		name:       name,
		mainCache:  cache.NewSafeCache(cap, minTTL, maxTTL),
		getterFunc: getterFunc,
		single:     singleflight.Singleflight{},
	}
	mu.Lock()
	defer mu.Unlock()
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	if v, ok := groups[name]; ok {
		return v
	} else {
		return nil
	}
}

func (g *Group) Get(key string) (eviction.ByteView, bool) {
	if key == "" {
		return eviction.ByteView{}, false
	}
	if v, ok := g.mainCache.Get(key); ok {
		fmt.Println("[GGCache] cache hit key:", key)
		return v, true
	} else {
		if local, ok := g.GetFromLocal(key); ok {
			return local, true
		} else {
			return eviction.ByteView{}, false
		}
	}
}

func (g *Group) GetFromLocal(key string) (eviction.ByteView, bool) {
	if value, ok := g.single.Do(key, func() (interface{}, bool) {
		if b, ok := g.getterFunc(key); ok {
			g.mainCache.Put(key, eviction.ByteView{B: b})
			return eviction.ByteView{B: b}, true
		} else {
			return eviction.ByteView{}, false
		}
	}); ok {
		return value.(eviction.ByteView), ok
	} else {
		return eviction.ByteView{}, false
	}
}

func (g *Group) Put(key string, value []byte) {
	g.mainCache.Put(key, eviction.ByteView{
		B: value,
	})
}
