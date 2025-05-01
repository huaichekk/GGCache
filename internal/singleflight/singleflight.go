package singleflight

import (
	"fmt"
	"sync"
)

type CallBack func() (interface{}, bool)
type Call struct {
	wg    sync.WaitGroup
	ok    bool
	value interface{}
}
type Singleflight struct {
	mu sync.Mutex
	m  map[string]*Call
}

func (s *Singleflight) Do(key string, fn CallBack) (interface{}, bool) {
	s.mu.Lock()
	if s.m == nil {
		s.m = make(map[string]*Call)
	}
	if v, ok := s.m[key]; ok { //已经有请求去执行回调函数了
		s.mu.Unlock()
		fmt.Println("wait first goroutine return")
		v.wg.Wait()
		return v.value, v.ok
	} else { //第一个开放执行回调
		c := new(Call)
		c.wg.Add(1)
		s.m[key] = c
		s.mu.Unlock()
		c.value, c.ok = fn()
		c.wg.Done()

		s.mu.Lock()
		delete(s.m, key)
		s.mu.Unlock()
		return c.value, c.ok
	}
}
