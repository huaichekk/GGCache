package eviction

import (
	"fmt"
	"testing"
)

type String struct {
	s string
}

func (s String) Len() int {
	return len(s.s)
}

func TestLRUCache_Get(t *testing.T) {
	m := map[string]String{
		"k1": String{s: "hh"},
		"k2": String{s: "gg"},
		"k3": String{s: "ll"},
	}
	var capacity int
	for k, v := range m {
		capacity += len(k) + v.Len()
	}
	lru := NewLRUCache(capacity)
	lru.Put("k1", m["k1"])
	lru.Put("k2", m["k2"])
	lru.Put("k3", m["k3"])
	lru.show()
	fmt.Println()
	if value, ok := lru.Get("k1"); !ok {
		t.Error("insert error", value)
	}
	lru.show()
	fmt.Println()
	lru.Put("k4", String{s: "aa"})
	lru.show()
}
