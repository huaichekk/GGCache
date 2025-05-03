package eviction

import "sync"

type Value interface {
	Len() int
}

type Cache interface {
	Get(key string) (value Value, ok bool)
	Put(key string, value Value)
	ScheduleDelete(mu *sync.RWMutex)
}
type Block struct {
	key   string
	value Value
	ttl   int64
}
