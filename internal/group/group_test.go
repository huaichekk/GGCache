package group

import (
	"fmt"
	"log"
	"testing"
	"time"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("scores", 2<<10, 100, 50, GetterFunc(
		func(key string) ([]byte, bool) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), true
			}
			return nil, false
		}))

	for k, v := range db {
		if view, ok := gee.Get(k); !ok || view.String() != v {
			t.Fatal("failed to get value of Tom")
		} // load from callback function
		if _, ok := gee.Get(k); !ok || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss %d", k, loadCounts[k])
		} // cache hit
	}

	if view, ok := gee.Get("unknown"); ok {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
	fmt.Println("test singlefly")
}

func TestSingleFly(t *testing.T) {
	g := NewGroup("scores", 2<<10, 50, 100, GetterFunc(
		func(key string) ([]byte, bool) {
			log.Println("[SlowDB] search key", key)
			time.Sleep(100 * time.Millisecond)
			if v, ok := db[key]; ok {
				return []byte(v), true
			}
			return nil, false
		}))

	for i := 0; i < 10; i++ {
		go func() {
			if v, ok := g.Get("Tom"); !ok || v.String() != "630" {
				t.Errorf("错误的信息")
			}
		}()
	}
	time.Sleep(1 * time.Second)
}

func TestGroup_Get(t *testing.T) {
	g := NewGroup("scores", 2<<10, 50, 100, GetterFunc(
		func(key string) ([]byte, bool) {
			log.Println("[SlowDB] search key", key)
			time.Sleep(100 * time.Millisecond)
			if v, ok := db[key]; ok {
				return []byte(v), true
			}
			return nil, false
		}))
	if v, ok := g.Get("kwj"); ok {
		fmt.Println(v)
	}
	g.Get("kwj")
	g.Get("kkk")
	g.Get("kkk")
}

func TestGroup_GetTTL(t *testing.T) {
	g := NewGroup("scores", 2<<10, 100, 200, GetterFunc(
		func(key string) ([]byte, bool) {
			log.Println("[SlowDB] search key", key)
			time.Sleep(100 * time.Millisecond)
			if v, ok := db[key]; ok {
				return []byte(v), true
			}
			return nil, false
		}))
	if v, ok := g.Get("Tom"); ok {
		fmt.Println(v)
	}
	for i := 0; i < 3; i++ {
		time.Sleep(5 * time.Second)
		g.Get("Tom")
	}
}
