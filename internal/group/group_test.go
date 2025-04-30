package group

import (
	"log"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("scores", 2<<10, GetterFunc(
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
}
