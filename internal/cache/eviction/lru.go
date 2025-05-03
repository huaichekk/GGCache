package eviction

import (
	"container/list"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type LRUCache struct {
	cap    int
	size   int
	cache  map[string]*list.Element
	list   *list.List
	minTTL int
	maxTTL int
}

func NewLRUCache(cap, minTTL, maxTTL int) *LRUCache {
	res := &LRUCache{
		cap:    cap,
		size:   0,
		cache:  make(map[string]*list.Element),
		list:   list.New(),
		minTTL: minTTL,
		maxTTL: maxTTL,
	}
	return res
}
func (lru *LRUCache) Get(key string) (Value, bool) {
	if v, ok := lru.cache[key]; ok {
		b := v.Value.(*Block)
		if b.ttl < time.Now().Unix() { //惰性删除
			lru.list.Remove(v)
			delete(lru.cache, key)
			lru.size -= len(key) + b.value.Len()
			return nil, false
		}
		b.ttl += int64(lru.RandomInt())
		lru.list.MoveToBack(v)
		return b.value, true
	}
	return nil, false
}
func (lru *LRUCache) Put(key string, value Value) {
	if v, ok := lru.cache[key]; ok { //老数据
		b := v.Value.(*Block)
		lru.list.MoveToBack(v)
		lru.size += value.Len() - b.value.Len()
		b.value = value
		b.ttl += int64(lru.RandomInt()) //续约
	} else { //新数据
		b := &Block{
			key:   key,
			value: value,
			ttl:   time.Now().Unix() + int64(lru.RandomInt()),
		}
		back := lru.list.PushBack(b)
		lru.cache[key] = back
		lru.size += len(key) + value.Len()
	}
	lru.evictionCache() //淘汰多余数据
}

func (lru *LRUCache) evictionCache() {
	for lru.size > lru.cap {
		front := lru.list.Front()
		b := front.Value.(*Block)
		delete(lru.cache, b.key)
		lru.size -= len(b.key) + b.value.Len()
		lru.list.Remove(front)
	}
}

func (lru *LRUCache) show() {
	cur := lru.list.Front()
	for cur != nil {
		fmt.Println(cur.Value.(*Block).key)
		cur = cur.Next()
	}
}

func (lru *LRUCache) RandomInt() int {
	return rand.Intn(lru.maxTTL-lru.minTTL) + lru.minTTL
}

func (lru *LRUCache) Remove(key string, value *list.Element) {
	delete(lru.cache, key)
	lru.list.Remove(value)
	lru.size -= len(key) + value.Value.(*Block).value.Len()
}

func (lru *LRUCache) ScheduleDelete(mu *sync.RWMutex) {
	go func() { //定时删除
		ticker := time.NewTicker(time.Second * 60 * 30)
		defer ticker.Stop()
		for _ = range ticker.C {
			mu.Lock()
			for k, v := range lru.cache {
				if v.Value.(*Block).ttl < time.Now().Unix() {
					lru.Remove(k, v)
				}
			}
			mu.Unlock()
		}
	}()
}
