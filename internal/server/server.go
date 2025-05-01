package server

import (
	"GGCache/configs"
	"GGCache/internal/cache/eviction"
	"GGCache/internal/client"
	"GGCache/internal/consistent"
	"GGCache/internal/group"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type HTTPPool struct {
	selfAddr   string
	mu         sync.Mutex
	consistent *consistent.Consistent
	nodes      map[string]*client.Client
}

func NewHTTPPool(addr string) *HTTPPool {
	return &HTTPPool{
		selfAddr:   addr,
		mu:         sync.Mutex{},
		consistent: consistent.NewConsistent(configs.GetConfig().CacheConfig.Replicas, nil),
		nodes:      make(map[string]*client.Client),
	}
}

func (s *HTTPPool) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	//http://127.0.0.1/group/key
	params := strings.Split(req.URL.Path, "/")
	params = params[1:]
	if len(params) != 2 {
		fmt.Println(params, len(params))
		http.Error(resp, "bad request", http.StatusBadRequest)
		return
	}
	if v, err := s.Get(params[0], params[1]); err == nil {
		resp.Header().Set("Content-Type", "application/octet-stream")
		_, _ = resp.Write(v.ByteSlice())
		return
	} else {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *HTTPPool) RegisterNode(addr string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.consistent.AddNode(addr)
	s.nodes[addr] = client.NewClient(addr)
}

func (s *HTTPPool) Get(groupName, key string) (eviction.ByteView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if addr := s.consistent.ChooseNode(key); addr == s.selfAddr { //本地找
		g := group.GetGroup(groupName)
		if g == nil {
			return eviction.ByteView{}, fmt.Errorf("no such group")
		}
		if v, ok := g.Get(key); ok {
			return v, nil
		} else {
			return eviction.ByteView{}, fmt.Errorf("[Local]key not find by cache and local")
		}
	} else { //远程节点找
		c := s.nodes[addr]
		if v, ok := c.Get(groupName, key); ok {
			fmt.Println("get from", addr)
			return v, nil
		} else {
			return eviction.ByteView{}, fmt.Errorf("[Peer]key not find by cache and local")
		}
	}
}
func (s *HTTPPool) DeleteNode(addr string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.consistent.DeleteNode(addr)
	delete(s.nodes, addr)
}
