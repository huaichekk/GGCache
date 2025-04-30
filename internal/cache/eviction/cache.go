package eviction

type Value interface {
	Len() int
}

type Cache interface {
	Get(key string) (value Value, ok bool)
	Put(key string, value Value)
}
