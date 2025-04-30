package eviction

import (
	"container/list"
	"fmt"
)

type LRUCache struct {
	cap   int
	size  int
	cache map[string]*list.Element
	list  *list.List
}
type Block struct {
	key   string
	value Value
}

func NewLRUCache(cap int) *LRUCache {
	return &LRUCache{
		cap:   cap,
		size:  0,
		cache: make(map[string]*list.Element),
		list:  list.New(),
	}
}
func (lru *LRUCache) Get(key string) (Value, bool) {
	if v, ok := lru.cache[key]; ok {
		b := v.Value.(*Block)
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
	} else { //新数据
		b := &Block{
			key:   key,
			value: value,
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
