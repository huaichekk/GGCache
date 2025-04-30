package cache

import (
	"GGCache/internal/cache/eviction"
	"fmt"
	"sync"
	"testing"
)

func TestSafeCache_GetPut(t *testing.T) {
	cache := NewSafeCache(100000000000)
	key := "testKey"
	value := eviction.ByteView{B: []byte("v")}

	// 测试Put和Get
	cache.Put(key, value)
	if v, ok := cache.Get(key); !ok {
		t.Error("Failed to get value that was just put")
	} else if v.String() != value.String() {
		t.Errorf("Got wrong value, expected %s, got %s", value.String(), v.String())
	}

	// 测试不存在的key
	if _, ok := cache.Get("nonExistentKey"); ok {
		t.Error("Got value for non-existent key")
	}
}
func TestSafeCache_ConcurrentAccess(t *testing.T) {
	cache := NewSafeCache(100000000)

	// 并发测试参数
	numGoroutines := 100
	keysPerGoroutine := 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// 启动多个goroutine并发读写
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < keysPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d-%d", goroutineID, j)
				value := eviction.ByteView{B: []byte(fmt.Sprintf("value-%d-%d", goroutineID, j))}

				// 写入
				cache.Put(key, value)

				// 读取验证
				if v, ok := cache.Get(key); !ok {
					t.Errorf("Failed to get key %s", key)
				} else if v.String() != value.String() {
					t.Errorf("Value mismatch for key %s", key)
				}
			}
		}(i)
	}

	wg.Wait()

	// 验证所有键值对
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < keysPerGoroutine; j++ {
			key := fmt.Sprintf("key-%d-%d", i, j)
			expectedValue := fmt.Sprintf("value-%d-%d", i, j)

			if v, ok := cache.Get(key); !ok {
				t.Errorf("Key %s missing after concurrent test", key)
			} else if v.String() != expectedValue {
				t.Errorf("Value mismatch for key %s after concurrent test", key)
			}
		}
	}
}
